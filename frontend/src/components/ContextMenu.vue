<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'

export interface MenuItem {
  label: string
  icon?: string  // svg path d
  action: () => void
  divider?: boolean
}

const props = defineProps<{
  x: number
  y: number
  items: MenuItem[]
}>()

const emit = defineEmits<{ close: [] }>()

function handleClick(item: MenuItem) {
  item.action()
  emit('close')
}

function handleOverlay() {
  emit('close')
}

onMounted(() => {
  document.addEventListener('click', handleOverlay)
})
onBeforeUnmount(() => {
  document.removeEventListener('click', handleOverlay)
})
</script>

<template>
  <div class="fixed inset-0 z-50">
    <div
      class="absolute bg-[var(--bg-elevated)] border border-[var(--border)] rounded-lg shadow-2xl py-1 min-w-[160px]"
      :style="{ left: x + 'px', top: y + 'px' }"
      @click.stop
    >
      <button
        v-for="(item, i) in items" :key="i"
        @click="handleClick(item)"
        class="w-full flex items-center gap-2 px-3 py-1.5 text-xs hover:bg-[var(--bg-hover)] transition-colors cursor-pointer text-left"
        :class="item.divider ? 'border-t border-[var(--border)] mt-1 pt-2' : ''"
      >
        <svg v-if="item.icon" class="w-3.5 h-3.5 text-[var(--text-muted)] flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="item.icon" />
        </svg>
        <span class="text-[var(--text-primary)]">{{ item.label }}</span>
      </button>
    </div>
  </div>
</template>
