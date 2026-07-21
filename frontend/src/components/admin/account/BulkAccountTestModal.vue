<template>
  <BaseDialog
    :show="show"
    :title="t('admin.accounts.bulkTest.title')"
    width="extra-wide"
    :show-close-button="!running"
    @close="handleClose"
  >
    <div class="space-y-4">
      <div class="rounded-lg border border-gray-200 bg-gray-50 p-3 text-sm text-gray-700 dark:border-dark-600 dark:bg-dark-700/50 dark:text-gray-200">
        <div class="flex flex-wrap items-center gap-x-4 gap-y-1">
          <span>{{ t('admin.accounts.bulkTest.selected', { count: accountIds.length }) }}</span>
          <span v-if="running || finished">
            {{ t('admin.accounts.bulkTest.progress', { done: doneCount, total: accountIds.length }) }}
          </span>
          <span v-if="finished" class="text-green-600 dark:text-green-400">
            {{ t('admin.accounts.bulkTest.successCount', { count: successCount }) }}
          </span>
          <span v-if="finished && failedCount > 0" class="text-red-600 dark:text-red-400">
            {{ t('admin.accounts.bulkTest.failedCount', { count: failedCount }) }}
          </span>
        </div>
        <div class="mt-2 h-2 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-dark-600">
          <div
            class="h-2 rounded-full bg-primary-600 transition-all"
            :style="{ width: `${progressPercent}%` }"
          />
        </div>
      </div>

      <div v-if="!running && !finished" class="grid gap-3 sm:grid-cols-2">
        <div class="space-y-1.5">
          <label class="text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('admin.accounts.bulkTest.concurrency') }}
          </label>
          <input
            v-model.number="concurrency"
            type="number"
            min="1"
            max="20"
            class="input w-full"
          />
          <p class="text-xs text-gray-500 dark:text-gray-400">
            {{ t('admin.accounts.bulkTest.concurrencyHint') }}
          </p>
        </div>
        <div class="space-y-1.5">
          <label class="text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('admin.accounts.bulkTest.modelOptional') }}
          </label>
          <input
            v-model="modelId"
            type="text"
            class="input w-full"
            :placeholder="t('admin.accounts.bulkTest.modelPlaceholder')"
          />
          <p class="text-xs text-gray-500 dark:text-gray-400">
            {{ t('admin.accounts.bulkTest.modelHint') }}
          </p>
        </div>
      </div>

      <div v-if="results.length > 0" class="space-y-2">
        <div class="flex items-center justify-between">
          <h4 class="text-sm font-semibold text-gray-900 dark:text-gray-100">
            {{ t('admin.accounts.bulkTest.results') }}
          </h4>
          <label class="flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
            <input v-model="onlyFailed" type="checkbox" class="rounded border-gray-300 text-primary-600" />
            {{ t('admin.accounts.bulkTest.onlyFailed') }}
          </label>
        </div>
        <div class="max-h-80 overflow-auto rounded-lg border border-gray-200 dark:border-dark-600">
          <table class="min-w-full divide-y divide-gray-200 text-sm dark:divide-dark-600">
            <thead class="bg-gray-50 dark:bg-dark-700">
              <tr>
                <th class="px-3 py-2 text-left font-medium text-gray-600 dark:text-gray-300">ID</th>
                <th class="px-3 py-2 text-left font-medium text-gray-600 dark:text-gray-300">
                  {{ t('admin.accounts.bulkTest.colName') }}
                </th>
                <th class="px-3 py-2 text-left font-medium text-gray-600 dark:text-gray-300">
                  {{ t('admin.accounts.bulkTest.colStatus') }}
                </th>
                <th class="px-3 py-2 text-left font-medium text-gray-600 dark:text-gray-300">
                  {{ t('admin.accounts.bulkTest.colLatency') }}
                </th>
                <th class="px-3 py-2 text-left font-medium text-gray-600 dark:text-gray-300">
                  {{ t('admin.accounts.bulkTest.colMessage') }}
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-600">
              <tr v-for="item in visibleResults" :key="item.account_id">
                <td class="px-3 py-2 text-gray-700 dark:text-gray-200">{{ item.account_id }}</td>
                <td class="px-3 py-2 text-gray-700 dark:text-gray-200">
                  {{ item.name || accountNameMap[item.account_id] || '-' }}
                </td>
                <td class="px-3 py-2">
                  <span
                    :class="item.status === 'success'
                      ? 'text-green-600 dark:text-green-400'
                      : 'text-red-600 dark:text-red-400'"
                  >
                    {{ item.status === 'success'
                      ? t('admin.accounts.bulkTest.statusSuccess')
                      : t('admin.accounts.bulkTest.statusFailed') }}
                  </span>
                </td>
                <td class="px-3 py-2 text-gray-600 dark:text-gray-300">
                  {{ item.latency_ms ? `${item.latency_ms}ms` : '-' }}
                </td>
                <td class="max-w-md truncate px-3 py-2 text-gray-600 dark:text-gray-300" :title="item.error_message || item.response_text || ''">
                  {{ item.error_message || item.response_text || '-' }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>
    </div>

    <template #footer>
      <div class="flex w-full items-center justify-end gap-2">
        <button
          v-if="!running"
          type="button"
          class="btn btn-secondary btn-sm"
          @click="handleClose"
        >
          {{ finished ? t('common.close') : t('common.cancel') }}
        </button>
        <button
          v-if="!running && !finished"
          type="button"
          class="btn btn-primary btn-sm"
          :disabled="accountIds.length === 0"
          @click="start"
        >
          {{ t('admin.accounts.bulkTest.start') }}
        </button>
        <button
          v-if="running"
          type="button"
          class="btn btn-danger btn-sm"
          @click="cancel"
        >
          {{ t('admin.accounts.bulkTest.cancel') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { adminAPI } from '@/api/admin'
import type { BulkTestAccountResult } from '@/api/admin/accounts'
import { extractApiErrorMessage } from '@/utils/apiError'

const props = defineProps<{
  show: boolean
  accountIds: number[]
  accountNameMap?: Record<number, string>
}>()

const emit = defineEmits<{
  close: []
  completed: [payload: { success: number; failed: number; results: BulkTestAccountResult[] }]
}>()

const { t } = useI18n()

const BATCH_SIZE = 50
const concurrency = ref(5)
const modelId = ref('')
const running = ref(false)
const finished = ref(false)
const cancelled = ref(false)
const results = ref<BulkTestAccountResult[]>([])
const doneCount = ref(0)
const errorMessage = ref('')
const onlyFailed = ref(false)

const accountNameMap = computed(() => props.accountNameMap ?? {})
const successCount = computed(() => results.value.filter((r) => r.status === 'success').length)
const failedCount = computed(() => results.value.filter((r) => r.status !== 'success').length)
const progressPercent = computed(() => {
  if (props.accountIds.length === 0) return 0
  return Math.min(100, Math.round((doneCount.value / props.accountIds.length) * 100))
})
const visibleResults = computed(() => {
  if (!onlyFailed.value) return results.value
  return results.value.filter((r) => r.status !== 'success')
})

watch(
  () => props.show,
  (show) => {
    if (!show) return
    running.value = false
    finished.value = false
    cancelled.value = false
    results.value = []
    doneCount.value = 0
    errorMessage.value = ''
    onlyFailed.value = false
    concurrency.value = 5
    modelId.value = ''
  }
)

function chunkIds(ids: number[], size: number): number[][] {
  const chunks: number[][] = []
  for (let i = 0; i < ids.length; i += size) {
    chunks.push(ids.slice(i, i + size))
  }
  return chunks
}

async function start() {
  if (running.value || props.accountIds.length === 0) return
  running.value = true
  finished.value = false
  cancelled.value = false
  results.value = []
  doneCount.value = 0
  errorMessage.value = ''

  const safeConcurrency = Math.min(20, Math.max(1, Number(concurrency.value) || 5))
  const chunks = chunkIds(props.accountIds, BATCH_SIZE)

  try {
    for (const chunk of chunks) {
      if (cancelled.value) break
      const response = await adminAPI.accounts.bulkTest(chunk, {
        model_id: modelId.value.trim() || undefined,
        concurrency: safeConcurrency,
        timeoutMs: 300000
      })
      results.value.push(...(response.results ?? []))
      doneCount.value = Math.min(props.accountIds.length, doneCount.value + chunk.length)
    }
    finished.value = true
    emit('completed', {
      success: successCount.value,
      failed: failedCount.value,
      results: results.value
    })
  } catch (error) {
    errorMessage.value = extractApiErrorMessage(error, t('admin.accounts.bulkTest.failed'))
  } finally {
    running.value = false
    if (!finished.value && cancelled.value) {
      finished.value = true
    }
  }
}

function cancel() {
  cancelled.value = true
}

function handleClose() {
  if (running.value) return
  emit('close')
}
</script>
