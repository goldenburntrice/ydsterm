import { ref, onMounted, onBeforeUnmount } from 'vue'
import { Events } from '@wailsio/runtime'

const syncing = ref(false)
const currentModule = ref('')
const syncError = ref('')

export function useSyncProgress() {
  function onSyncProgress(payload: any) {
    const p = payload.data
    if (!p) return

    if (p.Error) {
      syncError.value = p.Error
      syncing.value = false
      currentModule.value = ''
      return
    }

    if (p.Done) {
      currentModule.value = ''
      syncing.value = false
    } else {
      currentModule.value = p.Module
      syncing.value = true
      syncError.value = ''
    }
  }

  onMounted(() => {
    Events.On('sync-progress', onSyncProgress)
  })

  onBeforeUnmount(() => {
    Events.Off('sync-progress')
  })

  return { syncing, currentModule, syncError }
}
