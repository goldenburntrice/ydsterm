import { reactive, ref, computed } from 'vue'
import type { Terminal } from '@xterm/xterm'
import type { FitAddon } from '@xterm/addon-fit'
import { Disconnect } from '../../bindings/ydsterm/internal/service/terminalserviceimpl'
import { Disconnect as SFTPDisconnect } from '../../bindings/ydsterm/internal/service/sftpserviceimpl'

export interface AppTab {
  id: string
  type: 'home' | 'terminal' | 'sftp'
  title: string
  closable: boolean
  // terminal
  sessionId?: string
  hostId?: string
  hostName?: string
  connected?: boolean
  terminal?: Terminal | null
  fitAddon?: FitAddon | null
  resizeObs?: ResizeObserver | null
  // sftp
  sftpSessionId?: string
}

const tabs = reactive<AppTab[]>([
  { id: 'home', type: 'home', title: '主页', closable: false },
])

const activeTabId = ref('home')

export function useAppTabs() {
  const activeTab = computed(() => tabs.find(t => t.id === activeTabId.value))

  function openTerminal(hostId: string, hostName: string, sessionId: string) {
    // close existing connection to same host? no — allow duplicates per spec
    const id = `term-${hostId}-${Date.now()}`
    const tab: AppTab = {
      id,
      type: 'terminal',
      title: hostName,
      closable: true,
      sessionId,
      hostId,
      hostName,
      connected: true,
    }
    tabs.push(tab)
    activeTabId.value = id
    return id
  }

  function openSFTP(hostId: string, hostName: string, sftpSessionId: string) {
    const id = `sftp-${hostId}-${Date.now()}`
    const tab: AppTab = {
      id, type: 'sftp', title: `SFTP: ${hostName}`, closable: true,
      hostId, hostName, sftpSessionId,
    }
    tabs.push(tab)
    activeTabId.value = id
    return id
  }

  function closeTab(tabId: string) {
    const idx = tabs.findIndex(t => t.id === tabId)
    if (idx === -1) return
    const tab = tabs[idx]
    if (!tab.closable) return

    if (tab.type === 'terminal' && tab.sessionId) {
      Disconnect(tab.sessionId)
    } else if (tab.type === 'sftp' && tab.sftpSessionId) {
      SFTPDisconnect(tab.sftpSessionId)
    }

    tab.resizeObs?.disconnect()
    tab.terminal?.dispose()
    tabs.splice(idx, 1)

    if (activeTabId.value === tabId) {
      activeTabId.value = tabs[Math.min(idx, tabs.length - 1)]?.id || 'home'
    }
  }

  function switchTab(tabId: string) {
    activeTabId.value = tabId
  }

  function goHome() {
    activeTabId.value = 'home'
  }

  return {
    tabs,
    activeTabId,
    activeTab,
    openTerminal,
    openSFTP,
    closeTab,
    switchTab,
    goHome,
  }
}

// singleton
let instance: ReturnType<typeof useAppTabs> | null = null
export function getAppTabs() {
  if (!instance) instance = useAppTabs()
  return instance
}
