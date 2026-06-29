<script setup lang="ts">
import { ref } from 'vue'
import TabBar from './components/TabBar.vue'
import HomeView from './views/HomeView.vue'
import TerminalView from './views/TerminalView.vue'
import SFTPView from './views/SFTPView.vue'
import SettingsPanel from './components/SettingsPanel.vue'
import { getAppTabs } from './composables/useAppTabs'

const { tabs, activeTabId } = getAppTabs()
const showSettings = ref(false)
</script>

<template>
  <div class="flex flex-col h-screen w-screen bg-[var(--bg-base)]">
    <TabBar @open-settings="showSettings = true" />
    <main class="flex-1 overflow-hidden">
      <HomeView v-show="activeTabId === 'home'" />
      <template v-for="tab in tabs" :key="tab.id">
        <TerminalView v-if="tab.type === 'terminal'"
          :tab-id="tab.id" :session-id="tab.sessionId || ''"
          :host-name="tab.hostName || ''" :host-id="tab.hostId || ''"
          :active="activeTabId === tab.id" />
        <SFTPView v-if="tab.type === 'sftp'"
          :tab-id="tab.id" :sftp-session-id="tab.sftpSessionId || ''"
          :host-name="tab.hostName || ''"
          :active="activeTabId === tab.id" />
      </template>
    </main>
    <SettingsPanel v-if="showSettings" @close="showSettings = false" />
  </div>
</template>
