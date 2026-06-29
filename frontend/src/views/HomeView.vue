<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import type { Host, HostGroup, HostCreateInput, HostUpdateInput, KeyCreateInput, HostGroupCreateInput } from '../../bindings/ydsterm/internal/types/models'
import * as HostService from '../../bindings/ydsterm/internal/service/hostserviceimpl'
import { Connect as termConnect } from '../../bindings/ydsterm/internal/service/terminalserviceimpl'
import { Connect as sftpConnect } from '../../bindings/ydsterm/internal/service/sftpserviceimpl'
import { useModal } from '../composables/useModal'
import { getAppTabs } from '../composables/useAppTabs'
import { useTheme } from '../composables/useTheme'
import HostFormDialog from '../components/HostFormDialog.vue'
import KeyManagerDialog from '../components/KeyManagerDialog.vue'
import ContextMenu from '../components/ContextMenu.vue'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import type { MenuItem } from '../components/ContextMenu.vue'

const appVersion = import.meta.env.VITE_APP_VERSION || 'dev'

const hosts = ref<Host[]>([])
const groups = ref<HostGroup[]>([])
const keys = ref<any[]>([])
const searchQuery = ref('')
const hostForm = useModal()
const keyManager = useModal()
const editingHost = ref<Host | null>(null)
const ctxMenu = ref<{ x: number; y: number; items: MenuItem[] } | null>(null)
const selectedHostId = ref('')
const showAddMenu = ref(false)
const confirmDialog = ref<{
  visible: boolean
  mode: 'confirm' | 'prompt'
  title: string
  message: string
  defaultValue?: string
  danger?: boolean
  onConfirm?: (value: string) => void
}>({
  visible: false,
  mode: 'confirm',
  title: '',
  message: '',
})
const { theme, toggleTheme } = useTheme()

onMounted(async () => {
  const [h, g, k] = await Promise.all([
    HostService.List(''),
    HostService.ListGroups(),
    HostService.ListKeys(),
  ])
  hosts.value = h ?? []
  groups.value = g ?? []
  keys.value = k ?? []
})

function openNewHost(parentGroupId: string = '') {
  editingHost.value = null
  hostForm.open()
}

function openEditHost(h: Host) {
  editingHost.value = h
  hostForm.open()
}

async function handleHostSave(data: HostCreateInput | HostUpdateInput) {
  if (editingHost.value) {
    await HostService.Update(data as HostUpdateInput)
  } else {
    await HostService.Create(data as HostCreateInput)
  }
  hostForm.close()
  const [h, g] = await Promise.all([HostService.List(''), HostService.ListGroups()])
  hosts.value = h ?? []
  groups.value = g ?? []
}

async function handleDeleteHost(id: string) {
  await HostService.Delete(id)
  hosts.value = (await HostService.List('')) ?? []
}

async function handleKeyImport(data: { name: string; privateKey: string; passphrase?: string }) {
  await HostService.CreateKey(data as KeyCreateInput)
  keys.value = (await HostService.ListKeys()) ?? []
}

async function handleKeyDelete(id: string) {
  await HostService.DeleteKey(id)
  keys.value = (await HostService.ListKeys()) ?? []
}

// Build tree: groups with their hosts
interface TreeNode {
  id: string
  name: string
  type: 'group' | 'host'
  parentId: string
  children: TreeNode[]
  host?: Host
  expanded: boolean
}

const treeRef = ref<Record<string, boolean>>({})

const tree = computed(() => {
  const q = searchQuery.value.toLowerCase()
  const filteredHosts = q
    ? hosts.value.filter(h => h.name.toLowerCase().includes(q) || h.hostname.toLowerCase().includes(q))
    : hosts.value

  const root: TreeNode[] = []

  // hosts without a group (or all if searching)
  if (!q) {
    const ungrouped = filteredHosts.filter(h => !h.groupId)
    ungrouped.forEach(h => {
      root.push({ id: h.id, name: h.name, type: 'host', parentId: '', children: [], host: h, expanded: false })
    })
  } else {
    filteredHosts.forEach(h => {
      root.push({ id: h.id, name: h.name, type: 'host', parentId: h.groupId, children: [], host: h, expanded: false })
    })
  }

  // groups with their hosts
  if (!q) {
    groups.value.forEach(g => {
      const groupHosts = filteredHosts.filter(h => h.groupId === g.id)
      const node: TreeNode = {
        id: g.id, name: g.name, type: 'group', parentId: g.parentId, children: [], expanded: treeRef.value[g.id] ?? true,
      }
      groupHosts.forEach(h => {
        node.children.push({ id: h.id, name: h.name, type: 'host', parentId: g.id, children: [], host: h, expanded: false })
      })
      root.push(node)
    })
  }

  return root
})

function toggleGroup(id: string) {
  treeRef.value[id] = !(treeRef.value[id] ?? true)
}

async function connectHost(host: Host) {
  try {
    const sessionId = await termConnect(host.id) as string
    const { openTerminal } = getAppTabs()
    openTerminal(host.id, host.name, sessionId)
  } catch (e: any) {
    alert('连接失败: ' + (e?.message ?? e))
  }
}

async function connectSFTP(host: Host) {
  try {
    const sftpSessionId = await sftpConnect(host.id) as string
    const { openSFTP } = getAppTabs()
    openSFTP(host.id, host.name, sftpSessionId)
  } catch (e: any) {
    alert('SFTP 连接失败: ' + (e?.message ?? e))
  }
}

async function handleGroupsChanged() {
  groups.value = (await HostService.ListGroups()) ?? []
}

function showContextMenu(e: MouseEvent, host: Host) {
  e.preventDefault()
  e.stopPropagation()
  selectedHostId.value = host.id
  ctxMenu.value = {
    x: e.clientX,
    y: e.clientY,
    items: [
      { label: 'SSH 连接', icon: 'M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z', action: () => connectHost(host) },
      { label: 'SFTP 连接', icon: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z', action: () => connectSFTP(host) },
    ],
  }
}

function showGroupContextMenu(e: MouseEvent, node: TreeNode) {
  e.preventDefault()
  e.stopPropagation()
  ctxMenu.value = {
    x: e.clientX,
    y: e.clientY,
    items: [
      { label: '重命名', icon: 'M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z', action: () => renameGroup(node) },
      { label: '删除', icon: 'M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16', action: () => deleteGroup(node) },
    ],
  }
}

function createGroup() {
  confirmDialog.value = {
    visible: true,
    mode: 'prompt',
    title: '新增分组',
    message: '请输入分组名称：',
    defaultValue: '',
    onConfirm: async (name) => {
      if (!name.trim()) return
      await HostService.CreateGroup({ name: name.trim(), parentId: '' })
      groups.value = (await HostService.ListGroups()) ?? []
    }
  }
}

function renameGroup(node: TreeNode) {
  confirmDialog.value = {
    visible: true,
    mode: 'prompt',
    title: '重命名分组',
    message: `请输入「${node.name}」的新名称：`,
    defaultValue: node.name,
    onConfirm: async (newName) => {
      if (!newName.trim() || newName === node.name) return
      await HostService.UpdateGroup({ id: node.id, name: newName.trim() })
      groups.value = (await HostService.ListGroups()) ?? []
    }
  }
}

function deleteGroup(node: TreeNode) {
  confirmDialog.value = {
    visible: true,
    mode: 'confirm',
    title: '删除分组',
    message: `确定删除分组「${node.name}」吗？分组内的主机将变为未分组状态。`,
    danger: true,
    onConfirm: async () => {
      await HostService.DeleteGroup(node.id)
      groups.value = (await HostService.ListGroups()) ?? []
      hosts.value = (await HostService.List('')) ?? []
    }
  }
}

function handleConfirmDialog() {
  confirmDialog.value.onConfirm?.('')
  confirmDialog.value.visible = false
}

function handleConfirmDialogPrompt(value: string) {
  confirmDialog.value.onConfirm?.(value)
  confirmDialog.value.visible = false
}
</script>

<template>
  <div class="flex flex-col items-center h-full overflow-y-auto py-12 px-4 gap-8" @click="showAddMenu = false">
    <!-- Logo -->
    <div class="select-none mt-6">
      <pre class="font-mono text-[var(--accent)] text-center leading-tight mb-3"
style="font-size: 11px; line-height: 1.15; letter-spacing: 0">
██╗   ██╗██████╗ ███████╗████████╗███████╗██████╗ ███╗   ███╗
╚██╗ ██╔╝██╔══██╗██╔════╝╚══██╔══╝██╔════╝██╔══██╗████╗ ████║
 ╚████╔╝ ██║  ██║███████╗   ██║   █████╗  ██████╔╝██╔████╔██║
  ╚██╔╝  ██║  ██║╚════██║   ██║   ██╔══╝  ██╔══██╗██║╚██╔╝██║
   ██║   ██████╔╝███████║   ██║   ███████╗██║  ██║██║ ╚═╝ ██║
   ╚═╝   ╚═════╝ ╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝</pre>
      <div class="text-center mt-2 text-xs text-[var(--text-muted)] tracking-wide">
        Cross-platform terminal · {{ appVersion }}
      </div>
    </div>

    <!-- Search + Actions -->
    <div class="w-full max-w-xl flex items-center gap-2">
      <div class="relative flex-1">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
        <input v-model="searchQuery" type="text"
          class="w-full h-9 pr-4 bg-[var(--bg-surface)] border border-[var(--border)] rounded-lg text-sm focus:outline-none focus:border-[var(--accent)] transition-colors placeholder:text-[var(--text-muted)]"
          style="padding-left: 2.5rem"
          placeholder="搜索主机 (名称 / 地址)..." />
      </div>
      <div class="relative flex-shrink-0" @click.stop>
        <button @click="showAddMenu = !showAddMenu"
          class="w-9 h-9 flex items-center justify-center rounded-lg bg-[var(--accent)] text-white hover:opacity-90 transition-opacity cursor-pointer"
          title="新建">
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
        </button>
        <div v-if="showAddMenu" @click.stop
          class="absolute top-full right-0 mt-1 w-36 bg-[var(--bg-elevated)] border border-[var(--border)] rounded-lg shadow-lg z-10 py-1">
          <button @click="showAddMenu = false; openNewHost()"
            class="w-full flex items-center gap-2 px-3 py-2 text-sm text-[var(--text-primary)] hover:bg-[var(--bg-hover)] transition-colors cursor-pointer text-left">
            <svg class="w-4 h-4 text-[var(--text-muted)] flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/></svg>
            新增主机
          </button>
          <button @click="showAddMenu = false; createGroup()"
            class="w-full flex items-center gap-2 px-3 py-2 text-sm text-[var(--text-primary)] hover:bg-[var(--bg-hover)] transition-colors cursor-pointer text-left">
            <svg class="w-4 h-4 text-[var(--text-muted)] flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
            新增分组
          </button>
        </div>
      </div>
      <button @click="toggleTheme"
        class="w-9 h-9 flex items-center justify-center rounded-lg bg-[var(--bg-surface)] border border-[var(--border)] text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:border-[var(--text-muted)] transition-colors cursor-pointer flex-shrink-0"
        title="切换主题">
        <svg v-if="theme === 'dark'" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"/></svg>
        <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"/></svg>
      </button>
      <button @click="keyManager.open()"
        class="w-9 h-9 flex items-center justify-center rounded-lg bg-[var(--bg-surface)] border border-[var(--border)] text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:border-[var(--text-muted)] transition-colors cursor-pointer flex-shrink-0"
        title="密钥管理">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/></svg>
      </button>
    </div>

    <!-- Tree View -->
    <div class="w-full max-w-xl">
      <div v-if="hosts.length === 0" class="text-center text-[var(--text-muted)] py-8 text-sm">
        还没有主机，点击 + 添加
      </div>

      <div v-for="node in tree" :key="node.id">
        <!-- Group folder -->
        <template v-if="node.type === 'group'">
          <button
            @click="toggleGroup(node.id)"
            @contextmenu.prevent="showGroupContextMenu($event, node)"
            class="w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm hover:bg-[var(--bg-hover)] transition-colors cursor-pointer text-left"
          >
            <svg class="w-4 h-4 text-[var(--text-muted)] transition-transform flex-shrink-0" :class="{ 'rotate-90': node.expanded }" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
            <svg class="w-4 h-4 text-[var(--warning)] flex-shrink-0" fill="currentColor" viewBox="0 0 24 24"><path d="M10 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"/></svg>
            <span class="font-medium">{{ node.name }}</span>
            <span class="text-xs text-[var(--text-muted)] ml-auto">{{ node.children.length }}</span>
          </button>
          <div v-if="node.expanded" class="ml-4 border-l border-[var(--border)] pl-3 my-0.5">
            <div v-for="child in node.children" :key="child.id"
              @dblclick="child.host && connectHost(child.host)"
              @click="selectedHostId = child.host?.id || ''"
              @contextmenu.prevent="child.host && showContextMenu($event, child.host)"
              class="group flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm hover:bg-[var(--bg-hover)] cursor-pointer transition-colors"
              :class="selectedHostId === child.host?.id ? 'bg-[var(--bg-hover)]' : ''">
              <div class="w-2 h-2 rounded-full flex-shrink-0" :style="{ backgroundColor: child.host?.color || 'var(--accent)' }" />
              <span class="flex-1 truncate">{{ child.name }}</span>
              <span class="text-xs text-[var(--text-muted)] font-mono hidden group-hover:inline truncate max-w-[150px]">{{ child.host?.hostname }}</span>
              <button @click.stop="openEditHost(child.host!)" class="opacity-0 group-hover:opacity-100 p-1 rounded hover:bg-[var(--border)] text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-all flex-shrink-0">
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
              </button>
              <button @click.stop="handleDeleteHost(child.host!.id)" class="opacity-0 group-hover:opacity-100 p-1 rounded hover:bg-[var(--danger)]/15 text-[var(--text-muted)] hover:text-[var(--danger)] transition-all flex-shrink-0">
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
              </button>
            </div>
          </div>
        </template>

        <!-- Ungrouped host -->
        <template v-if="node.type === 'host' && !searchQuery && !node.host?.groupId">
          <div
            @dblclick="node.host && connectHost(node.host)"
            @click="selectedHostId = node.host?.id || ''"
            @contextmenu.prevent="node.host && showContextMenu($event, node.host)"
            class="group flex items-center gap-2 px-3 py-2 rounded-lg text-sm hover:bg-[var(--bg-hover)] cursor-pointer transition-colors"
            :class="selectedHostId === node.host?.id ? 'bg-[var(--bg-hover)]' : ''"
          >
            <div class="w-2 h-2 rounded-full flex-shrink-0" :style="{ backgroundColor: node.host?.color || 'var(--accent)' }" />
            <svg class="w-4 h-4 text-[var(--text-muted)] flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/></svg>
            <span class="flex-1 truncate">{{ node.name }}</span>
            <span class="text-xs text-[var(--text-muted)] font-mono hidden group-hover:inline truncate max-w-[150px]">{{ node.host?.hostname }}</span>
            <button @click.stop="openEditHost(node.host!)" class="opacity-0 group-hover:opacity-100 p-1 rounded hover:bg-[var(--border)] text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-all flex-shrink-0">
              <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
            </button>
            <button @click.stop="handleDeleteHost(node.host!.id)" class="opacity-0 group-hover:opacity-100 p-1 rounded hover:bg-[var(--danger)]/15 text-[var(--text-muted)] hover:text-[var(--danger)] transition-all flex-shrink-0">
              <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
            </button>
          </div>
        </template>

        <!-- Search results: host shown regardless of group -->
        <template v-if="node.type === 'host' && searchQuery">
          <div
            @dblclick="node.host && connectHost(node.host)"
            @click="selectedHostId = node.host?.id || ''"
            @contextmenu.prevent="node.host && showContextMenu($event, node.host)"
            class="group flex items-center gap-2 px-3 py-2 rounded-lg text-sm hover:bg-[var(--bg-hover)] cursor-pointer transition-colors"
            :class="selectedHostId === node.host?.id ? 'bg-[var(--bg-hover)]' : ''"
          >
            <div class="w-2 h-2 rounded-full flex-shrink-0" :style="{ backgroundColor: node.host?.color || 'var(--accent)' }" />
            <svg class="w-4 h-4 text-[var(--text-muted)] flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"/></svg>
            <span class="flex-1 truncate">{{ node.name }}</span>
            <span class="text-xs text-[var(--text-muted)] font-mono truncate max-w-[150px]">{{ node.host?.hostname }}</span>
          </div>
        </template>
      </div>
    </div>

    <HostFormDialog :visible="hostForm.visible.value" :host="editingHost" :groups="groups" :keys="keys"
      @close="hostForm.close()" @save="handleHostSave" @groupsChanged="handleGroupsChanged" />
    <KeyManagerDialog :visible="keyManager.visible.value" :keys="keys"
      @close="keyManager.close()" @import="handleKeyImport" @deleteKey="handleKeyDelete" />
    <ContextMenu v-if="ctxMenu" :x="ctxMenu.x" :y="ctxMenu.y" :items="ctxMenu.items" @close="ctxMenu = null" />
    <ConfirmDialog
      :visible="confirmDialog.visible"
      :mode="confirmDialog.mode"
      :title="confirmDialog.title"
      :message="confirmDialog.message"
      :default-value="confirmDialog.defaultValue"
      :danger="confirmDialog.danger"
      @close="confirmDialog.visible = false"
      @confirm="confirmDialog.mode === 'prompt' ? handleConfirmDialogPrompt($event) : handleConfirmDialog()"
    />
  </div>
</template>
