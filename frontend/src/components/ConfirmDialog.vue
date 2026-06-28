<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  visible: boolean
  title: string
  message: string
  mode: 'confirm' | 'prompt'
  defaultValue?: string
  danger?: boolean
}>()

const emit = defineEmits<{
  close: []
  confirm: [value: string]
}>()

const inputValue = ref('')

watch(() => props.visible, (val) => {
  if (val && props.mode === 'prompt') {
    inputValue.value = props.defaultValue ?? ''
  }
})

function handleConfirm() {
  emit('confirm', props.mode === 'prompt' ? inputValue.value : '')
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="emit('close')" />
      <div class="relative w-full max-w-sm mx-4 bg-[var(--bg-elevated)] border border-[var(--border)] rounded-xl shadow-2xl">
        <div class="flex items-center justify-between px-5 py-3.5 border-b border-[var(--border)]">
          <h2 class="text-sm font-semibold">{{ title }}</h2>
          <button @click="emit('close')" class="p-1 rounded hover:bg-[var(--bg-hover)] text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors cursor-pointer">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>

        <div class="p-5 space-y-4">
          <p class="text-sm text-[var(--text-secondary)]">{{ message }}</p>

          <input v-if="mode === 'prompt'" v-model="inputValue" type="text"
            class="w-full px-3 py-2 bg-[var(--bg-base)] border border-[var(--border)] rounded-md text-sm focus:outline-none focus:border-[var(--accent)] transition-colors"
            @keyup.enter="handleConfirm"
            @keyup.escape="emit('close')" />

          <div class="flex justify-end gap-2 pt-2">
            <button @click="emit('close')"
              class="px-4 py-2 rounded-md text-xs bg-[var(--bg-base)] border border-[var(--border)] text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors cursor-pointer">取消</button>
            <button @click="handleConfirm"
              class="px-4 py-2 rounded-md text-xs text-white hover:opacity-90 transition-opacity font-medium cursor-pointer"
              :class="danger ? 'bg-[var(--danger)]' : 'bg-[var(--accent)]'">确定</button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
