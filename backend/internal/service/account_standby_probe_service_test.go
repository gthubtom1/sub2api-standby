package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func schedulableTestAccount(id int64, priority int, lastUsed *time.Time, platform string, extra map[string]any) *Account {
	return &Account{
		ID:          id,
		Priority:    priority,
		LastUsedAt:  lastUsed,
		Platform:    platform,
		Status:      StatusActive,
		Schedulable: true,
		Extra:       extra,
	}
}

func TestSelectStandbyProbeCandidates_NeedsTwoPerActive(t *testing.T) {
	now := time.Now()
	recent := now.Add(-5 * time.Minute)
	old := now.Add(-2 * time.Hour)

	// 2 active + 5 idle, none verified -> need 4 standbys probed
	pool := []*Account{
		schedulableTestAccount(1, 1, &recent, "anthropic", nil),
		schedulableTestAccount(2, 1, &recent, "anthropic", nil),
		schedulableTestAccount(3, 2, &old, "anthropic", nil),
		schedulableTestAccount(4, 3, &old, "anthropic", nil),
		schedulableTestAccount(5, 4, &old, "anthropic", nil),
		schedulableTestAccount(6, 5, &old, "anthropic", nil),
		schedulableTestAccount(7, 6, &old, "anthropic", nil),
	}

	got := selectStandbyProbeCandidates(pool, now, 30*time.Minute, 15*time.Minute, 2, 0, 5*time.Minute, 5*time.Minute)
	require.Len(t, got, 4)
	// Should not include active accounts
	for _, acc := range got {
		require.NotContains(t, []int64{1, 2}, acc.ID)
	}
}

func TestSelectStandbyProbeCandidates_SkipsWhenNoActive(t *testing.T) {
	now := time.Now()
	old := now.Add(-2 * time.Hour)
	pool := []*Account{
		schedulableTestAccount(3, 1, &old, "openai", nil),
		schedulableTestAccount(4, 2, &old, "openai", nil),
	}
	got := selectStandbyProbeCandidates(pool, now, 30*time.Minute, 15*time.Minute, 2, 0, 5*time.Minute, 5*time.Minute)
	require.Empty(t, got)
}

func TestSelectStandbyProbeCandidates_RateLimitedStillCreatesDemand(t *testing.T) {
	now := time.Now()
	old := now.Add(-2 * time.Hour)
	resetAt := now.Add(10 * time.Minute)

	// Hot account is rate-limited (no recent LastUsedAt) + idle standbys.
	pool := []*Account{
		{
			ID:               1,
			Priority:         1,
			LastUsedAt:       &old,
			Platform:         "anthropic",
			Status:           StatusActive,
			Schedulable:      true,
			RateLimitResetAt: &resetAt,
		},
		schedulableTestAccount(10, 1, &old, "anthropic", nil),
		schedulableTestAccount(11, 2, &old, "anthropic", nil),
		schedulableTestAccount(12, 3, &old, "anthropic", nil),
	}

	got := selectStandbyProbeCandidates(pool, now, 30*time.Minute, 15*time.Minute, 2, 0, 5*time.Minute, 5*time.Minute)
	require.Len(t, got, 2)
	for _, acc := range got {
		require.NotEqual(t, int64(1), acc.ID)
		require.True(t, acc.IsSchedulable())
	}
}

func TestSelectStandbyProbeCandidates_HysteresisDemand(t *testing.T) {
	now := time.Now()
	old := now.Add(-2 * time.Hour)
	pool := []*Account{
		schedulableTestAccount(10, 1, &old, "openai", nil),
		schedulableTestAccount(11, 2, &old, "openai", nil),
		schedulableTestAccount(12, 3, &old, "openai", nil),
	}
	// Live demand 0, but caller passes remembered demand=2 -> need 4, only 3 idle
	got := selectStandbyProbeCandidates(pool, now, 30*time.Minute, 15*time.Minute, 2, 2, 5*time.Minute, 5*time.Minute)
	require.Len(t, got, 3)
}

func TestSelectStandbyProbeCandidates_RespectsVerifiedTTL(t *testing.T) {
	now := time.Now()
	recent := now.Add(-2 * time.Minute)
	old := now.Add(-2 * time.Hour)
	verifiedTS := now.Add(-5 * time.Minute).Unix()

	pool := []*Account{
		schedulableTestAccount(1, 1, &recent, "gemini", nil),
		schedulableTestAccount(10, 1, &old, "gemini", map[string]any{standbyProbeOKAtExtraKey: verifiedTS}),
		schedulableTestAccount(11, 2, &old, "gemini", map[string]any{standbyProbeOKAtExtraKey: verifiedTS}),
		schedulableTestAccount(12, 3, &old, "gemini", nil),
	}
	// need = 1*2 = 2, already have 2 verified -> no probes
	got := selectStandbyProbeCandidates(pool, now, 30*time.Minute, 15*time.Minute, 2, 0, 5*time.Minute, 5*time.Minute)
	require.Empty(t, got)

	// Only 1 verified -> need one more
	pool[2].Extra = nil
	got = selectStandbyProbeCandidates(pool, now, 30*time.Minute, 15*time.Minute, 2, 0, 5*time.Minute, 5*time.Minute)
	require.Len(t, got, 1)
	require.Contains(t, []int64{11, 12}, got[0].ID)
}

func TestSelectStandbyProbeCandidates_SkipsUnschedulableIdle(t *testing.T) {
	now := time.Now()
	recent := now.Add(-1 * time.Minute)
	old := now.Add(-2 * time.Hour)
	pool := []*Account{
		schedulableTestAccount(1, 1, &recent, "anthropic", nil),
		{
			ID:          10,
			Priority:    1,
			LastUsedAt:  &old,
			Platform:    "anthropic",
			Status:      StatusActive,
			Schedulable: false, // manually disabled
		},
		schedulableTestAccount(11, 2, &old, "anthropic", nil),
	}
	got := selectStandbyProbeCandidates(pool, now, 30*time.Minute, 15*time.Minute, 2, 0, 5*time.Minute, 5*time.Minute)
	require.Len(t, got, 1)
	require.Equal(t, int64(11), got[0].ID)
}

func TestResolvePoolDemand_Hysteresis(t *testing.T) {
	svc := NewAccountStandbyProbeService(nil, nil, nil)
	now := time.Now()
	require.Equal(t, 3, svc.resolvePoolDemand("anthropic|g=0", 3, now, 30*time.Minute))
	// Cold pool remembers demand inside window
	require.Equal(t, 3, svc.resolvePoolDemand("anthropic|g=0", 0, now.Add(5*time.Minute), 30*time.Minute))
	// Expired memory
	require.Equal(t, 0, svc.resolvePoolDemand("anthropic|g=0", 0, now.Add(40*time.Minute), 30*time.Minute))
}

func TestBuildStandbyPools_ByGroupAndPlatform(t *testing.T) {
	pool := buildStandbyPools([]*Account{
		{ID: 1, Platform: "anthropic", GroupIDs: []int64{10}},
		{ID: 2, Platform: "anthropic", GroupIDs: []int64{10, 20}},
		{ID: 3, Platform: "openai", GroupIDs: nil},
	})
	require.Contains(t, pool, "anthropic|g=10")
	require.Contains(t, pool, "anthropic|g=20")
	require.Contains(t, pool, "openai|g=0")
	require.Len(t, pool["anthropic|g=10"], 2)
	require.Len(t, pool["anthropic|g=20"], 1)
}

func TestCompareStandbyProbeFreshness(t *testing.T) {
	now := time.Now()
	fresh := schedulableTestAccount(1, 1, nil, "openai", map[string]any{standbyProbeOKAtExtraKey: now.Add(-2 * time.Minute).Unix()})
	cold := schedulableTestAccount(2, 1, nil, "openai", nil)
	require.Equal(t, -1, compareStandbyProbeFreshness(fresh, cold, now))
	require.Equal(t, 1, compareStandbyProbeFreshness(cold, fresh, now))
	require.Equal(t, 0, compareStandbyProbeFreshness(fresh, fresh, now))
}

func TestSelectStandbyProbeCandidates_FailCooldownSkips(t *testing.T) {
	now := time.Now()
	recent := now.Add(-1 * time.Minute)
	old := now.Add(-2 * time.Hour)
	pool := []*Account{
		schedulableTestAccount(1, 1, &recent, "anthropic", nil),
		schedulableTestAccount(10, 1, &old, "anthropic", map[string]any{standbyProbeFailAtExtraKey: now.Add(-1 * time.Minute).Unix()}),
		schedulableTestAccount(11, 2, &old, "anthropic", nil),
	}
	got := selectStandbyProbeCandidates(pool, now, 30*time.Minute, 15*time.Minute, 2, 0, 5*time.Minute, 5*time.Minute)
	require.Len(t, got, 1)
	require.Equal(t, int64(11), got[0].ID)
}
