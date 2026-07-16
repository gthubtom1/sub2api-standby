<template>
  <div class="group/standby relative">
    <span :class="['badge text-xs', badgeClass]" :title="title">
      {{ label }}
    </span>
    <div
      class="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 w-52 -translate-x-1/2 whitespace-normal rounded bg-gray-900 px-3 py-2 text-center text-xs leading-relaxed text-white opacity-0 transition-opacity group-hover/standby:opacity-100 dark:bg-gray-700"
    >
      {{ title }}
      <div
        class="absolute left-1/2 top-full -translate-x-1/2 border-4 border-transparent border-t-gray-900 dark:border-t-gray-700"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Account } from '@/types'
import { formatStandbyProbeTime, getStandbyProbeInfo } from '@/utils/standbyProbe'

const props = defineProps<{
  account: Account
}>()

const { t } = useI18n()

const info = computed(() =>
  getStandbyProbeInfo(props.account.extra as Record<string, unknown> | undefined)
)

const label = computed(() => {
  switch (info.value.state) {
    case 'healthy':
      return t('admin.accounts.status.standbyHealthy')
    case 'stale':
      return t('admin.accounts.status.standbyStale')
    case 'failed':
      return t('admin.accounts.status.standbyFailed')
    default:
      return t('admin.accounts.status.standbyUnknown')
  }
})

const badgeClass = computed(() => {
  switch (info.value.state) {
    case 'healthy':
      return 'badge-success'
    case 'stale':
      return 'badge-warning'
    case 'failed':
      return 'badge-danger'
    default:
      return 'badge-gray'
  }
})

const title = computed(() => {
  const okText = formatStandbyProbeTime(info.value.okAt)
  const failText = formatStandbyProbeTime(info.value.failAt)
  switch (info.value.state) {
    case 'healthy':
      return t('admin.accounts.status.standbyHealthyHint', { time: okText || '-' })
    case 'stale':
      return t('admin.accounts.status.standbyStaleHint', { time: okText || '-' })
    case 'failed':
      return t('admin.accounts.status.standbyFailedHint', { time: failText || '-' })
    default:
      return t('admin.accounts.status.standbyUnknownHint')
  }
})
</script>

