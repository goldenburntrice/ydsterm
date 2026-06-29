<script setup lang="ts">
import { getAppTabs } from '../composables/useAppTabs'
import { useTheme } from '../composables/useTheme'

const emit = defineEmits<{ 'open-settings': [] }>()

const { tabs, activeTabId, switchTab, closeTab } = getAppTabs()
const { theme, toggleTheme } = useTheme()
</script>

<template>
  <div class="flex items-center bg-[var(--bg-surface)] border-b border-[var(--border)] h-10 select-none pl-1">
    <div class="flex-1 flex items-center overflow-x-auto">
      <div
        v-for="tab in tabs" :key="tab.id"
        @click="switchTab(tab.id)"
        :title="tab.title"
        class="group relative flex items-center gap-2 px-3 py-2 text-xs cursor-pointer transition-colors min-w-0 max-w-[200px] border-r border-[var(--border-muted)]"
        :class="activeTabId === tab.id
          ? 'bg-[var(--tab-active-bg)] text-[var(--text-primary)]'
          : 'text-[var(--text-secondary)] hover:bg-[var(--tab-hover-bg)] hover:text-[var(--text-primary)]'"
      >
        <div v-if="activeTabId === tab.id" class="absolute inset-x-0 bottom-0 h-0.5 bg-[var(--accent)] rounded-full" />

        <template v-if="tab.type === 'home'">
          <svg class="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" /></svg>
        </template>
        <template v-else-if="tab.type === 'sftp'">
          <svg class="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24"><path d="M10 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"/></svg>
        </template>
        <template v-else>
          <svg class="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
        </template>

        <span v-if="tab.type === 'terminal'" class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="tab.connected ? 'bg-emerald-400' : 'bg-red-400'" />

        <span class="truncate">{{ tab.title }}</span>

        <button
          v-if="tab.closable"
          @click.stop="closeTab(tab.id)"
          class="opacity-0 group-hover:opacity-100 p-0.5 rounded hover:bg-[var(--border)] transition-all flex-shrink-0 ml-0.5"
        >
          <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
        </button>
      </div>
    </div>

    <button
      @click="toggleTheme"
      class="px-2 h-full flex items-center justify-center hover:bg-[var(--bg-hover)] text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors cursor-pointer flex-shrink-0 border-l border-[var(--border-muted)]"
      :title="theme === 'dark' ? '切换到亮色' : '切换到暗色'"
    >
      <svg v-if="theme === 'dark'" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"/></svg>
      <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"/></svg>
    </button>

    <button
      @click="emit('open-settings')"
      class="px-2 h-full flex items-center justify-center hover:bg-[var(--bg-hover)] text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors cursor-pointer flex-shrink-0 border-l border-[var(--border-muted)]"
      title="设置"
    >
      <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
    </button>
  </div>
</template>
