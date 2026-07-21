package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

const (
	standbyProbeOKAtExtraKey     = "standby_probe_ok_at"
	standbyProbeFailAtExtraKey   = "standby_probe_fail_at"
	standbyProbeLastStatusKey    = "standby_probe_last_status"
	standbyProbeDefaultInterval  = 2
	standbyProbeDefaultActiveWin = 30
	standbyProbeDefaultPerActive = 2
	standbyProbeDefaultTTL       = 15
	standbyProbeDefaultConcur    = 3
	standbyProbeDefaultMaxProbes = 40
	standbyProbeDefaultTimeout   = 240
	standbyProbeLeaderLockKey    = "standby_probe:leader"
	standbyProbeDefaultDeficitSec  = 30
	standbyProbeDefaultRefreshMin  = 5
	standbyProbeDefaultFailCooldown  = 5
	standbyProbeDefaultFirstDelay  = 5
)

// AccountStandbyProbeService proactively connection-tests standby accounts so
// that when actively used accounts exhaust quota or fail, the scheduler can
// pick pre-verified healthy accounts instead of burning client retries.
type AccountStandbyProbeService struct {
	accountRepo    AccountRepository
	accountTestSvc *AccountTestService
	cfg            *config.Config

	lockCache  LeaderLockCache
	db         *sql.DB
	instanceID string

	startOnce sync.Once
	stopOnce  sync.Once
	stopCh    chan struct{}
	doneCh    chan struct{}

	cycleMu     sync.Mutex
	cycleCancel context.CancelFunc

	// demandMu guards pool demand hysteresis so a pool that just lost all
	// schedulable actives (rate-limit / overload / error) still warms standbys
	// for one active window.
	demandMu   sync.Mutex
	poolDemand map[string]standbyPoolDemand
}

type standbyPoolDemand struct {
	count int
	at    time.Time
}

// rankedStandbyCandidate is used to fairly cap max_probes across pools.
type rankedStandbyCandidate struct {
	account *Account
	deficit int
	age     int64
}

// NewAccountStandbyProbeService creates the worker (not started).
func NewAccountStandbyProbeService(
	accountRepo AccountRepository,
	accountTestSvc *AccountTestService,
	cfg *config.Config,
) *AccountStandbyProbeService {
	return &AccountStandbyProbeService{
		accountRepo:    accountRepo,
		accountTestSvc: accountTestSvc,
		cfg:            cfg,
		instanceID:     uuid.NewString(),
		stopCh:         make(chan struct{}),
		doneCh:         make(chan struct{}),
		poolDemand:     make(map[string]standbyPoolDemand),
	}
}

// SetLeaderLock injects multi-instance single-flight coordination.
func (s *AccountStandbyProbeService) SetLeaderLock(lockCache LeaderLockCache, db *sql.DB) {
	if s == nil {
		return
	}
	s.lockCache = lockCache
	s.db = db
}

// Start launches the background loop when enabled.
func (s *AccountStandbyProbeService) Start() {
	if s == nil {
		return
	}
	s.startOnce.Do(func() {
		if !s.enabled() {
			close(s.doneCh)
			logger.LegacyPrintf("service.standby_probe", "[StandbyProbe] disabled")
			return
		}
		go s.loop()
		logger.LegacyPrintf("service.standby_probe",
			"[StandbyProbe] started interval=%dm standby_per_active=%d health_ttl=%dm",
			s.intervalMinutes(), s.standbyPerActive(), s.healthTTLMinutes())
	})
}

// Stop stops the worker and cancels any in-flight cycle.
func (s *AccountStandbyProbeService) Stop() {
	if s == nil {
		return
	}
	s.stopOnce.Do(func() {
		select {
		case <-s.stopCh:
		default:
			close(s.stopCh)
		}
		s.cycleMu.Lock()
		if s.cycleCancel != nil {
			s.cycleCancel()
		}
		s.cycleMu.Unlock()

		wait := time.Duration(s.cycleTimeoutSeconds()+5) * time.Second
		if wait < 10*time.Second {
			wait = 10 * time.Second
		}
		select {
		case <-s.doneCh:
		case <-time.After(wait):
			logger.LegacyPrintf("service.standby_probe", "[StandbyProbe] stop timed out")
		}
	})
}

func (s *AccountStandbyProbeService) loop() {
	defer close(s.doneCh)

	// First cycle shortly after boot so cold pools warm up.
	timer := time.NewTimer(time.Duration(s.firstDelaySeconds()) * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-timer.C:
			deficit := s.runCycle()
			next := time.Duration(s.intervalMinutes()) * time.Minute
			if deficit > 0 {
				fast := time.Duration(s.deficitIntervalSeconds()) * time.Second
				if fast > 0 && fast < next {
					next = fast
				}
			}
			timer.Reset(next)
		}
	}
}

func (s *AccountStandbyProbeService) runCycle() int {
	if s.accountRepo == nil || s.accountTestSvc == nil {
		return 0
	}

	// Parent cancels when Stop() is called so long probes abort promptly.
	parent, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()
	go func() {
		select {
		case <-s.stopCh:
			parentCancel()
		case <-parent.Done():
		}
	}()

	timeout := time.Duration(s.cycleTimeoutSeconds()) * time.Second
	ctx, cancel := context.WithTimeout(parent, timeout)
	s.cycleMu.Lock()
	s.cycleCancel = cancel
	s.cycleMu.Unlock()
	defer func() {
		cancel()
		s.cycleMu.Lock()
		s.cycleCancel = nil
		s.cycleMu.Unlock()
	}()

	release, ok := tryAcquireSingletonLeaderLock(
		ctx, s.lockCache, s.db, standbyProbeLeaderLockKey, s.instanceID, timeout+30*time.Second,
	)
	if !ok {
		slog.Debug("standby_probe_skip_not_leader")
		return 0
	}
	defer release()

	// ListActive (not ListSchedulable): rate-limited / overloaded actives still
	// produce demand so we keep warming replacements while they are hot.
	accounts, err := s.accountRepo.ListActive(ctx)
	if err != nil {
		logger.LegacyPrintf("service.standby_probe", "[StandbyProbe] ListActive error: %v", err)
		return 0
	}
	if len(accounts) == 0 {
		return 0
	}

	now := time.Now()
	ptrs := make([]*Account, 0, len(accounts))
	for i := range accounts {
		acc := &accounts[i]
		if acc == nil || !acc.IsActive() {
			continue
		}
		ptrs = append(ptrs, acc)
	}
	if len(ptrs) == 0 {
		return 0
	}

	pools := buildStandbyPools(ptrs)
	activeWin := time.Duration(s.activeWindowMinutes()) * time.Minute
	healthTTL := time.Duration(s.healthTTLMinutes()) * time.Minute
	perActive := s.standbyPerActive()

	ranked := make([]rankedStandbyCandidate, 0)
	seen := make(map[int64]struct{})
	hotPools := 0
	plannedDeficit := 0

	for key, pool := range pools {
		currentDemand := countStandbyInUse(pool, now, activeWin)
		demand := s.resolvePoolDemand(key, currentDemand, now, activeWin)
		if demand == 0 {
			continue
		}
		hotPools++
		needIDs := selectStandbyProbeCandidates(pool, now, activeWin, healthTTL, perActive, demand, time.Duration(s.refreshBeforeTTLMinutes())*time.Minute, time.Duration(s.failCooldownMinutes())*time.Minute)
		if len(needIDs) == 0 {
			// Count real unmet verified slots only (not full need), so healthy pools stay calm.
			plannedDeficit += estimateStandbyDeficit(pool, now, activeWin, healthTTL, perActive, demand)
			continue
		}
		deficit := len(needIDs)
		plannedDeficit += deficit
		slog.Debug("standby_probe_pool",
			"pool", key,
			"demand", demand,
			"current_in_use", currentDemand,
			"need", deficit,
			"size", len(pool),
		)
		for _, acc := range needIDs {
			if acc == nil {
				continue
			}
			if _, exists := seen[acc.ID]; exists {
				continue
			}
			seen[acc.ID] = struct{}{}
			ranked = append(ranked, rankedStandbyCandidate{
				account: acc,
				deficit: deficit,
				age:     standbyProbeAgeScore(acc, now),
			})
		}
	}

	if len(ranked) == 0 {
		return plannedDeficit
	}

	// Prefer pools with larger deficit, then never-probed / stale, then priority/id.
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].deficit != ranked[j].deficit {
			return ranked[i].deficit > ranked[j].deficit
		}
		if ranked[i].age != ranked[j].age {
			return ranked[i].age > ranked[j].age
		}
		ai, aj := ranked[i].account, ranked[j].account
		if ai.Priority != aj.Priority {
			return ai.Priority < aj.Priority
		}
		return ai.ID < aj.ID
	})

	maxProbes := s.maxProbesPerCycle()
	if len(ranked) > maxProbes {
		ranked = ranked[:maxProbes]
	}

	candidates := make([]*Account, 0, len(ranked))
	for _, r := range ranked {
		candidates = append(candidates, r.account)
	}

	logger.LegacyPrintf("service.standby_probe",
		"[StandbyProbe] cycle start candidates=%d hot_pools=%d total_pools=%d",
		len(candidates), hotPools, len(pools))

	var (
		mu       sync.Mutex
		successN int
		failedN  int
	)

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(s.concurrency())

	for _, acc := range candidates {
		account := acc
		g.Go(func() error {
			if gctx.Err() != nil {
				return nil
			}
			okProbe, errMsg := s.probeOne(gctx, account)
			mu.Lock()
			if okProbe {
				successN++
			} else {
				failedN++
				if errMsg != "" {
					slog.Info("standby_probe_failed",
						"account_id", account.ID,
						"platform", account.Platform,
						"error", errMsg)
				}
			}
			mu.Unlock()
			return nil
		})
	}
	_ = g.Wait()

	remaining := failedN
	if plannedDeficit > successN+failedN {
		remaining += plannedDeficit - (successN + failedN)
	}
	logger.LegacyPrintf("service.standby_probe",
		"[StandbyProbe] cycle done success=%d failed=%d remaining_deficit~=%d", successN, failedN, remaining)
	return remaining
}

// resolvePoolDemand returns effective demand: live in-use count, or remembered
// demand while still inside the active window after the pool went cold.
func (s *AccountStandbyProbeService) resolvePoolDemand(key string, current int, now time.Time, activeWindow time.Duration) int {
	s.demandMu.Lock()
	defer s.demandMu.Unlock()

	if current > 0 {
		s.poolDemand[key] = standbyPoolDemand{count: current, at: now}
		return current
	}

	mem, ok := s.poolDemand[key]
	if !ok || mem.count <= 0 {
		return 0
	}
	if now.Sub(mem.at) > activeWindow {
		delete(s.poolDemand, key)
		return 0
	}
	return mem.count
}

func (s *AccountStandbyProbeService) probeOne(ctx context.Context, account *Account) (bool, string) {
	if account == nil {
		return false, "nil account"
	}

	// Empty model_id -> platform defaults inside TestAccountConnection.
	result, err := s.accountTestSvc.RunTestBackground(ctx, account.ID, "")
	now := time.Now()

	if err != nil {
		_ = s.accountRepo.UpdateExtra(ctx, account.ID, map[string]any{
			standbyProbeFailAtExtraKey: now.Unix(),
			standbyProbeLastStatusKey:  "failed",
		})
		return false, err.Error()
	}
	if result == nil {
		_ = s.accountRepo.UpdateExtra(ctx, account.ID, map[string]any{
			standbyProbeFailAtExtraKey: now.Unix(),
			standbyProbeLastStatusKey:  "failed",
		})
		return false, "empty result"
	}
	if result.Status != "success" {
		_ = s.accountRepo.UpdateExtra(ctx, account.ID, map[string]any{
			standbyProbeFailAtExtraKey: now.Unix(),
			standbyProbeLastStatusKey:  "failed",
		})
		msg := result.ErrorMessage
		if msg == "" {
			msg = "test failed"
		}
		return false, msg
	}

	_ = s.accountRepo.UpdateExtra(ctx, account.ID, map[string]any{
		standbyProbeOKAtExtraKey:  now.Unix(),
		standbyProbeLastStatusKey: "success",
	})
	return true, ""
}

type standbyPoolKey struct {
	Platform string
	GroupID  int64 // 0 = ungrouped
}

func (k standbyPoolKey) String() string {
	return fmt.Sprintf("%s|g=%d", k.Platform, k.GroupID)
}

func buildStandbyPools(accounts []*Account) map[string][]*Account {
	pools := make(map[string][]*Account)
	for _, acc := range accounts {
		if acc == nil {
			continue
		}
		platform := strings.TrimSpace(acc.Platform)
		if platform == "" {
			platform = "unknown"
		}
		groupIDs := acc.GroupIDs
		if len(groupIDs) == 0 {
			key := standbyPoolKey{Platform: platform, GroupID: 0}.String()
			pools[key] = append(pools[key], acc)
			continue
		}
		seenG := make(map[int64]struct{}, len(groupIDs))
		for _, gid := range groupIDs {
			if gid <= 0 {
				continue
			}
			if _, ok := seenG[gid]; ok {
				continue
			}
			seenG[gid] = struct{}{}
			key := standbyPoolKey{Platform: platform, GroupID: gid}.String()
			pools[key] = append(pools[key], acc)
		}
		if len(seenG) == 0 {
			key := standbyPoolKey{Platform: platform, GroupID: 0}.String()
			pools[key] = append(pools[key], acc)
		}
	}
	return pools
}


func estimateStandbyDeficit(pool []*Account, now time.Time, activeWindow, healthTTL time.Duration, standbyPerActive, demand int) int {
	if demand <= 0 || standbyPerActive <= 0 {
		return 0
	}
	need := demand * standbyPerActive
	verified := 0
	for _, acc := range pool {
		if acc == nil || isStandbyInUse(acc, now, activeWindow) || !isStandbyProbeEligible(acc) {
			continue
		}
		if isStandbyRecentlyVerified(acc, now, healthTTL) {
			verified++
		}
	}
	deficit := need - verified
	if deficit < 0 {
		return 0
	}
	return deficit
}

func countStandbyInUse(pool []*Account, now time.Time, activeWindow time.Duration) int {
	n := 0
	for _, acc := range pool {
		if isStandbyInUse(acc, now, activeWindow) {
			n++
		}
	}
	return n
}

// selectStandbyProbeCandidates returns idle *schedulable* accounts that should
// be probed so the pool has (demand * standbyPerActive) recently-verified standbys.
//
// demand is the effective in-use count (live or hysteretic). Accounts that are
// rate-limited/overloaded count toward demand but are never probe candidates.
// Pass demand<=0 to derive demand from live isStandbyInUse only.
func selectStandbyProbeCandidates(
	pool []*Account,
	now time.Time,
	activeWindow time.Duration,
	healthTTL time.Duration,
	standbyPerActive int,
	demand int,
	refreshBefore time.Duration,
	failCooldown time.Duration,
) []*Account {
	if standbyPerActive <= 0 {
		standbyPerActive = standbyProbeDefaultPerActive
	}
	if demand <= 0 {
		demand = countStandbyInUse(pool, now, activeWindow)
	}
	if demand <= 0 {
		return nil
	}

	need := demand * standbyPerActive
	verified := 0
	candidates := make([]*Account, 0)
	refresh := make([]*Account, 0)
	for _, acc := range pool {
		if acc == nil {
			continue
		}
		// Hot accounts are demand, not standbys.
		if isStandbyInUse(acc, now, activeWindow) {
			continue
		}
		// Only probe accounts the scheduler can actually pick next.
		if !isStandbyProbeEligible(acc) {
			continue
		}
		if isStandbyRecentlyVerified(acc, now, healthTTL) {
			verified++
			// Re-probe before TTL dies so unattended pools never go stale.
			if needsStandbyProbeRefresh(acc, now, healthTTL, refreshBefore) {
				refresh = append(refresh, acc)
			}
			continue
		}
		// Avoid thrashing the same dead account every deficit cycle.
		if isStandbyProbeFailCooldown(acc, now, failCooldown) {
			continue
		}
		candidates = append(candidates, acc)
	}

	deficit := need - verified
	sortStandbyProbeCandidates(candidates, now)
	sortStandbyProbeCandidates(refresh, now)

	out := make([]*Account, 0, deficit+len(refresh))
	if deficit > 0 {
		if len(candidates) > deficit {
			out = append(out, candidates[:deficit]...)
		} else {
			out = append(out, candidates...)
		}
	}
	// Keep aging verified accounts warm so TTL never silently expires unattended.
	refreshBudget := standbyPerActive
	if refreshBudget < 1 {
		refreshBudget = 1
	}
	addedRefresh := 0
	for _, acc := range refresh {
		if addedRefresh >= refreshBudget {
			break
		}
		dup := false
		for _, x := range out {
			if x != nil && acc != nil && x.ID == acc.ID {
				dup = true
				break
			}
		}
		if !dup {
			out = append(out, acc)
			addedRefresh++
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func sortStandbyProbeCandidates(candidates []*Account, now time.Time) {
	sort.SliceStable(candidates, func(i, j int) bool {
		ai, aj := candidates[i], candidates[j]
		ti := standbyProbeAgeScore(ai, now)
		tj := standbyProbeAgeScore(aj, now)
		if ti != tj {
			return ti > tj // older / never probed first
		}
		if ai.Priority != aj.Priority {
			return ai.Priority < aj.Priority
		}
		return ai.ID < aj.ID
	})
}

func needsStandbyProbeRefresh(acc *Account, now time.Time, healthTTL, refreshBefore time.Duration) bool {
	if acc == nil || healthTTL <= 0 {
		return false
	}
	if refreshBefore <= 0 {
		refreshBefore = time.Duration(standbyProbeDefaultRefreshMin) * time.Minute
	}
	ts := readStandbyProbeUnix(acc, standbyProbeOKAtExtraKey)
	if ts <= 0 {
		return false
	}
	age := now.Sub(time.Unix(ts, 0))
	if age < 0 {
		return false
	}
	// Refresh when remaining TTL is under refreshBefore.
	return age >= healthTTL-refreshBefore
}

func isStandbyProbeFailCooldown(acc *Account, now time.Time, cooldown time.Duration) bool {
	if acc == nil || cooldown <= 0 {
		return false
	}
	failAt := readStandbyProbeUnix(acc, standbyProbeFailAtExtraKey)
	if failAt <= 0 {
		return false
	}
	okAt := readStandbyProbeUnix(acc, standbyProbeOKAtExtraKey)
	if okAt > failAt {
		return false
	}
	return now.Sub(time.Unix(failAt, 0)) < cooldown
}


// isStandbyProbeEligible reports whether an idle account is worth probing as a
// healthy standby. Account-level schedulable is required; Grok near-exhausted
// quota windows are excluded so "预备健康" does not prefer soon-to-429 accounts.
func isStandbyProbeEligible(acc *Account) bool {
	if acc == nil || !acc.IsSchedulable() {
		return false
	}
	if acc.IsGrok() {
		if paused, _ := shouldAutoPauseGrokAccountByQuota(acc); paused {
			return false
		}
	}
	return true
}

func isStandbyInUse(acc *Account, now time.Time, activeWindow time.Duration) bool {
	if acc == nil {
		return false
	}
	// Recently used by gateway traffic.
	if acc.LastUsedAt != nil && now.Sub(*acc.LastUsedAt) <= activeWindow {
		return true
	}
	// Just hit rate-limit / overload window - treat as hot so we warm replacements.
	// Requires ListActive (not ListSchedulable) so these rows are still visible.
	if acc.IsRateLimited() {
		return true
	}
	if acc.IsOverloaded() {
		return true
	}
	return false
}

func isStandbyRecentlyVerified(acc *Account, now time.Time, healthTTL time.Duration) bool {
	ts := readStandbyProbeUnix(acc, standbyProbeOKAtExtraKey)
	if ts <= 0 {
		return false
	}
	okAt := time.Unix(ts, 0)
	return now.Sub(okAt) <= healthTTL
}

func standbyProbeAgeScore(acc *Account, now time.Time) int64 {
	// Higher score = should probe sooner.
	okAt := readStandbyProbeUnix(acc, standbyProbeOKAtExtraKey)
	if okAt <= 0 {
		return now.Unix() // never verified
	}
	return now.Unix() - okAt
}

func readStandbyProbeUnix(acc *Account, key string) int64 {
	if acc == nil || acc.Extra == nil {
		return 0
	}
	raw, ok := acc.Extra[key]
	if !ok || raw == nil {
		return 0
	}
	switch v := raw.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		if n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			return n
		}
		if t, err := time.Parse(time.RFC3339, strings.TrimSpace(v)); err == nil {
			return t.Unix()
		}
	}
	return 0
}

func (s *AccountStandbyProbeService) enabled() bool {
	if s == nil || s.cfg == nil {
		return true
	}
	return s.cfg.StandbyProbe.Enabled
}

func (s *AccountStandbyProbeService) intervalMinutes() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.CheckIntervalMinutes > 0 {
		return s.cfg.StandbyProbe.CheckIntervalMinutes
	}
	return standbyProbeDefaultInterval
}

func (s *AccountStandbyProbeService) activeWindowMinutes() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.ActiveWindowMinutes > 0 {
		return s.cfg.StandbyProbe.ActiveWindowMinutes
	}
	return standbyProbeDefaultActiveWin
}

func (s *AccountStandbyProbeService) standbyPerActive() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.StandbyPerActive > 0 {
		return s.cfg.StandbyProbe.StandbyPerActive
	}
	return standbyProbeDefaultPerActive
}

func (s *AccountStandbyProbeService) healthTTLMinutes() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.HealthTTLMinutes > 0 {
		return s.cfg.StandbyProbe.HealthTTLMinutes
	}
	return standbyProbeDefaultTTL
}

func (s *AccountStandbyProbeService) concurrency() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.Concurrency > 0 {
		return s.cfg.StandbyProbe.Concurrency
	}
	return standbyProbeDefaultConcur
}

func (s *AccountStandbyProbeService) maxProbesPerCycle() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.MaxProbesPerCycle > 0 {
		return s.cfg.StandbyProbe.MaxProbesPerCycle
	}
	return standbyProbeDefaultMaxProbes
}

func (s *AccountStandbyProbeService) cycleTimeoutSeconds() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.CycleTimeoutSeconds > 0 {
		return s.cfg.StandbyProbe.CycleTimeoutSeconds
	}
	return standbyProbeDefaultTimeout
}

func (s *AccountStandbyProbeService) deficitIntervalSeconds() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.DeficitIntervalSeconds > 0 {
		return s.cfg.StandbyProbe.DeficitIntervalSeconds
	}
	return standbyProbeDefaultDeficitSec
}

func (s *AccountStandbyProbeService) firstDelaySeconds() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.FirstDelaySeconds > 0 {
		return s.cfg.StandbyProbe.FirstDelaySeconds
	}
	return standbyProbeDefaultFirstDelay
}

func (s *AccountStandbyProbeService) refreshBeforeTTLMinutes() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.RefreshBeforeTTLMinutes > 0 {
		return s.cfg.StandbyProbe.RefreshBeforeTTLMinutes
	}
	return standbyProbeDefaultRefreshMin
}

func (s *AccountStandbyProbeService) failCooldownMinutes() int {
	if s != nil && s.cfg != nil && s.cfg.StandbyProbe.FailCooldownMinutes > 0 {
		return s.cfg.StandbyProbe.FailCooldownMinutes
	}
	return standbyProbeDefaultFailCooldown
}
