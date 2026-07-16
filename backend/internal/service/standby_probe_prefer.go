package service

import "time"

// DefaultStandbyProbeHealthTTL is the scheduler-side freshness window used when
// deciding whether a pre-probed standby account is still trusted.
// Kept in sync with standbyProbeDefaultTTL minutes.
const DefaultStandbyProbeHealthTTL = 15 * time.Minute

// IsStandbyProbeFresh reports whether account has a successful standby probe
// within the health TTL. Used by schedulers to prefer pre-verified accounts
// during unattended failover so clients do not burn retries on cold accounts.
func IsStandbyProbeFresh(acc *Account, now time.Time) bool {
	if acc == nil {
		return false
	}
	if now.IsZero() {
		now = time.Now()
	}
	return isStandbyRecentlyVerified(acc, now, DefaultStandbyProbeHealthTTL)
}

// compareStandbyProbeFreshness returns -1 if a is fresher/better, 1 if b is, 0 if equal.
func compareStandbyProbeFreshness(a, b *Account, now time.Time) int {
	af := IsStandbyProbeFresh(a, now)
	bf := IsStandbyProbeFresh(b, now)
	if af == bf {
		return 0
	}
	if af {
		return -1
	}
	return 1
}
