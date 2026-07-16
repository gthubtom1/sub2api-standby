/** Standby probe health helpers (matches backend standby_probe_* extra keys). */

export const STANDBY_PROBE_OK_AT_KEY = "standby_probe_ok_at"
export const STANDBY_PROBE_FAIL_AT_KEY = "standby_probe_fail_at"
export const STANDBY_PROBE_LAST_STATUS_KEY = "standby_probe_last_status"

/** Default health TTL aligned with backend standbyProbeDefaultTTL (15 minutes). */
export const STANDBY_PROBE_HEALTH_TTL_MS = 15 * 60 * 1000

export type StandbyProbeState = "healthy" | "stale" | "failed" | "unknown"

export interface StandbyProbeInfo {
  state: StandbyProbeState
  okAt: number | null
  failAt: number | null
  lastStatus: string | null
}

function readUnixSeconds(raw: unknown): number | null {
  if (raw == null) return null
  if (typeof raw === "number" && Number.isFinite(raw) && raw > 0) {
    return Math.floor(raw)
  }
  if (typeof raw === "string") {
    const trimmed = raw.trim()
    if (!trimmed) return null
    const asInt = Number.parseInt(trimmed, 10)
    if (Number.isFinite(asInt) && asInt > 0 && String(asInt) === trimmed) {
      return asInt
    }
    const ms = Date.parse(trimmed)
    if (Number.isFinite(ms) && ms > 0) {
      return Math.floor(ms / 1000)
    }
  }
  return null
}

export function getStandbyProbeInfo(
  extra: Record<string, unknown> | null | undefined,
  nowMs: number = Date.now(),
  healthTtlMs: number = STANDBY_PROBE_HEALTH_TTL_MS
): StandbyProbeInfo {
  const okAt = readUnixSeconds(extra?.[STANDBY_PROBE_OK_AT_KEY])
  const failAt = readUnixSeconds(extra?.[STANDBY_PROBE_FAIL_AT_KEY])
  const lastStatusRaw = extra?.[STANDBY_PROBE_LAST_STATUS_KEY]
  const lastStatus = typeof lastStatusRaw === "string" ? lastStatusRaw : null

  const failIsNewer = failAt != null && (okAt == null || failAt >= okAt)
  if (failIsNewer && (lastStatus === "failed" || okAt == null || (okAt != null && failAt > okAt))) {
    return { state: "failed", okAt, failAt, lastStatus }
  }

  if (okAt != null) {
    const ageMs = nowMs - okAt * 1000
    if (ageMs >= 0 && ageMs <= healthTtlMs) {
      return { state: "healthy", okAt, failAt, lastStatus }
    }
    return { state: "stale", okAt, failAt, lastStatus }
  }

  return { state: "unknown", okAt, failAt, lastStatus }
}

export function formatStandbyProbeTime(unixSec: number | null): string {
  if (unixSec == null || unixSec <= 0) return ""
  const d = new Date(unixSec * 1000)
  if (Number.isNaN(d.getTime())) return ""
  const pad = (n: number) => String(n).padStart(2, "0")
  return `${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

