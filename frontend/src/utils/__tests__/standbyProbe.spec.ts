import { describe, expect, it } from 'vitest'
import { getStandbyProbeInfo, STANDBY_PROBE_HEALTH_TTL_MS } from '../standbyProbe'

describe('getStandbyProbeInfo', () => {
  const now = Date.parse('2026-07-16T12:00:00Z')

  it('marks recent ok as healthy', () => {
    const okAt = Math.floor((now - 5 * 60 * 1000) / 1000)
    const info = getStandbyProbeInfo({ standby_probe_ok_at: okAt }, now)
    expect(info.state).toBe('healthy')
  })

  it('marks old ok as stale', () => {
    const okAt = Math.floor((now - STANDBY_PROBE_HEALTH_TTL_MS - 60_000) / 1000)
    const info = getStandbyProbeInfo({ standby_probe_ok_at: okAt }, now)
    expect(info.state).toBe('stale')
  })

  it('marks newer fail as failed', () => {
    const okAt = Math.floor((now - 60 * 60 * 1000) / 1000)
    const failAt = Math.floor((now - 2 * 60 * 1000) / 1000)
    const info = getStandbyProbeInfo({
      standby_probe_ok_at: okAt,
      standby_probe_fail_at: failAt,
      standby_probe_last_status: 'failed',
    }, now)
    expect(info.state).toBe('failed')
  })

  it('marks empty extra as unknown', () => {
    expect(getStandbyProbeInfo({}, now).state).toBe('unknown')
  })
})

