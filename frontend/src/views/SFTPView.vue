<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { Events } from '@wailsio/runtime'
import * as SFTP from '../../bindings/ydsterm/internal/service/sftpserviceimpl'
import type { FileEntry } from '../../bindings/ydsterm/internal/service/models'
import type { MenuItem } from '../components/ContextMenu.vue'
import ContextMenu from '../components/ContextMenu.vue'

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB'
  if (bytes < 1073741824) return (bytes / 1048576).toFixed(1) + ' MB'
  return (bytes / 1073741824).toFixed(1) + ' GB'
}

const props = defineProps<{ tabId: string; sftpSessionId: string; hostName: string; active: boolean }>()

interface Panel {
  id: string; label: string; path: string; entries: FileEntry[]; loading: boolean; isLocal: boolean
}

const leftPanels = ref<Panel[]>([])
const rightPanels = ref<Panel[]>([])
const leftActive = ref('')
const rightActive = ref('')
const pathInputLeft = ref('')
const pathInputRight = ref('')
const selectedLeft = ref('')
const selectedRight = ref('')
const toastMsg = ref('')
const toastType = ref<'error' | 'success'>('error')
let toastTimer: ReturnType<typeof setTimeout> | null = null

function showToast(msg: string, type: 'error' | 'success' = 'error') {
  toastMsg.value = msg
  toastType.value = type
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => { toastMsg.value = '' }, 4000)
}

// ---------- context menu ----------

interface DragTransfer {
  path: string; name: string; isDir: boolean; fromLocal: boolean
  targetId: string; targetIsLocal: boolean
}
const ctxMenuVisible = ref(false)
const ctxMenuX = ref(0)
const ctxMenuY = ref(0)
const ctxMenuItems = ref<MenuItem[]>([])
const ctxMenuPanelId = ref('')
const ctxMenuIsLocal = ref(false)
const ctxMenuEntry = ref<FileEntry | null>(null)
const conflict = ref<DragTransfer | null>(null)
const conflictNewName = ref('')

// ---------- transfer queue ----------

interface TransferItem {
  id: string
  name: string
  size: string
  dir: 'up' | 'down'
  side: 'left' | 'right'
  pct: number
  speed: string
  done: boolean
  err: string
  startedAt: number
  bytesDone: number
  bytesTotal: number
}
const transfers = ref<TransferItem[]>([])

onMounted(() => {
  addPanel(true)
  addPanel(false)
  Events.On('sftp-transfer-progress', onTransferProgress)
})

onBeforeUnmount(() => {
  Events.Off('sftp-transfer-progress')
})

function findPanel(isLocal: boolean, id: string) { return (isLocal ? leftPanels.value : rightPanels.value).find(p => p.id === id) }
function activePanel(isLocal: boolean): string { return isLocal ? leftActive.value : rightActive.value }

function addPanel(isLocal: boolean) {
  const id = (isLocal ? 'L' : 'R') + Date.now()
  const list = isLocal ? leftPanels : rightPanels
  list.value.push({ id, label: isLocal ? '\u672c\u5730' : props.hostName, path: '', entries: [], loading: false, isLocal })
  if (isLocal) leftActive.value = id; else rightActive.value = id
  const p = list.value[list.value.length - 1]
  loadDir(p)
}

async function loadDir(p: Panel) {
  p.loading = true
  try {
    if (p.isLocal) {
      p.entries = await SFTP.ListLocal(p.path) ?? []
      if (!p.path) p.path = await SFTP.HomeLocal()
    } else {
      const remotePath = p.path || '/'
      p.entries = await SFTP.List(props.sftpSessionId, remotePath) ?? []
      if (!p.path) p.path = remotePath
    }
  } catch (_) { }
  p.loading = false
  syncPathInput(p)
}

function refreshActive(isLocal: boolean) {
  const p = findPanel(isLocal, activePanel(isLocal))
  if (p) loadDir(p)
}

function syncPathInput(p: Panel) {
  const val = p.path || '/'
  if (p.isLocal) pathInputLeft.value = val
  else pathInputRight.value = val
}

function selectEntry(isLocal: boolean, e: FileEntry) {
  if (isLocal) selectedLeft.value = e.path
  else selectedRight.value = e.path
}

function navigate(panelId: string, isLocal: boolean, e: FileEntry) {
  selectEntry(isLocal, e)
  if (!e.isDir) return
  const p = findPanel(isLocal, panelId)
  if (!p) return
  p.path = e.path
  loadDir(p)
}

async function goUp(panelId: string, isLocal: boolean) {
  const p = findPanel(isLocal, panelId)
  if (!p) return
  p.path = await SFTP.ParentDir(p.path, p.isLocal)
  if (isLocal) selectedLeft.value = ''
  else selectedRight.value = ''
  loadDir(p)
}

function closePanel(panelId: string, isLocal: boolean) {
  const list = isLocal ? leftPanels : rightPanels
  const active = isLocal ? leftActive : rightActive
  const idx = list.value.findIndex(x => x.id === panelId)
  if (idx === -1) return
  list.value.splice(idx, 1)
  if (active.value === panelId) active.value = list.value[Math.min(idx, list.value.length - 1)]?.id || ''
}

const activeLeft = computed(() => leftPanels.value.find(p => p.id === leftActive.value))
const activeRight = computed(() => rightPanels.value.find(p => p.id === rightActive.value))

// ---------- right-click ----------

function showCtxMenu(e: MouseEvent, entry: FileEntry, panelId: string, isLocal: boolean) {
  e.preventDefault()
  e.stopPropagation()
  selectEntry(isLocal, entry)
  ctxMenuEntry.value = entry
  ctxMenuPanelId.value = panelId
  ctxMenuIsLocal.value = isLocal
  ctxMenuX.value = e.clientX
  ctxMenuY.value = e.clientY
  ctxMenuItems.value = [
    {
      label: '\u91cd\u547d\u540d',
      icon: 'M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z',
      action: () => handleRename(),
    },
    {
      label: '\u5220\u9664',
      icon: 'M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16',
      action: () => handleDelete(),
    },
    {
      label: '\u65b0\u5efa\u6587\u4ef6\u5939',
      icon: 'M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z',
      action: () => handleNewFolder(),
    },
    {
      label: '\u5237\u65b0',
      icon: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15',
      action: () => refreshActive(isLocal),
    },
  ]
  ctxMenuVisible.value = true
}

function showPanelCtxMenu(e: MouseEvent, isLocal: boolean) {
  e.preventDefault()
  ctxMenuEntry.value = null
  ctxMenuPanelId.value = ''
  ctxMenuIsLocal.value = isLocal
  ctxMenuX.value = e.clientX
  ctxMenuY.value = e.clientY
  ctxMenuItems.value = [
    {
      label: '\u65b0\u5efa\u6587\u4ef6\u5939',
      icon: 'M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z',
      action: () => handleNewFolder(),
    },
    {
      label: '\u5237\u65b0',
      icon: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15',
      action: () => refreshActive(isLocal),
    },
  ]
  ctxMenuVisible.value = true
}

function closeCtxMenu() {
  ctxMenuVisible.value = false
  ctxMenuEntry.value = null
}

async function handleRename() {
  const entry = ctxMenuEntry.value
  if (!entry) return
  const newName = prompt('\u65b0\u540d\u79f0:', entry.name)
  if (!newName || newName === entry.name) return
  closeCtxMenu()
  try {
    const p = findPanel(ctxMenuIsLocal.value, ctxMenuPanelId.value)
    if (!p) return
    const dir = p.path.replace(/[\/\\]+$/, '')
    const sep = ctxMenuIsLocal.value ? '\\' : '/'
    const newPath = (dir ? dir + sep : '') + newName
    if (ctxMenuIsLocal.value) {
      await SFTP.RenameLocal(entry.path, newPath)
    } else {
      await SFTP.Rename(props.sftpSessionId, entry.path, newPath)
    }
    loadDir(p)
    showToast('\u5df2\u91cd\u547d\u540d\uff1a' + newName, 'success')
  } catch (e: any) { showToast('\u91cd\u547d\u540d\u5931\u8d25: ' + (e?.message ?? e)) }
}

async function handleDelete() {
  const entry = ctxMenuEntry.value
  if (!entry) return
  if (!confirm('\u786e\u5b9a\u5220\u9664 ' + entry.name + ' \uff1f')) return
  closeCtxMenu()
  try {
    const p = findPanel(ctxMenuIsLocal.value, ctxMenuPanelId.value)
    if (!p) return
    if (ctxMenuIsLocal.value) {
      await SFTP.DeleteLocal(entry.path)
    } else {
      await SFTP.Delete(props.sftpSessionId, entry.path)
    }
    loadDir(p)
    showToast('\u5df2\u5220\u9664\uff1a' + entry.name, 'success')
  } catch (e: any) { showToast('\u5220\u9664\u5931\u8d25: ' + (e?.message ?? e)) }
}

async function handleNewFolder() {
  const name = prompt('\u6587\u4ef6\u5939\u540d\u79f0:')
  if (!name) return
  closeCtxMenu()
  try {
    if (ctxMenuIsLocal.value) {
      const p = findPanel(true, activePanel(true))
      if (!p) return
      const dir = p.path.replace(/[\/\\]+$/, '')
      await SFTP.MkdirLocal((dir ? dir + '\\' : '') + name)
      loadDir(p)
    } else {
      const p = findPanel(false, activePanel(false))
      if (!p) return
      const dir = p.path.replace(/[\/\\]+$/, '')
      await SFTP.Mkdir(props.sftpSessionId, (dir ? dir + '/' : '') + name)
      loadDir(p)
    }
    showToast('\u5df2\u521b\u5efa\u6587\u4ef6\u5939\uff1a' + name, 'success')
  } catch (e: any) { showToast('\u521b\u5efa\u5931\u8d25: ' + (e?.message ?? e)) }
}

// ---------- drag & drop ----------

function onDragStart(e: DragEvent, entry: FileEntry, panel: Panel) {
  if (!e.dataTransfer) return
  e.dataTransfer.setData('text/json', JSON.stringify({
    path: entry.path, name: entry.name, isDir: entry.isDir, fromLocal: panel.isLocal,
  }))
  e.dataTransfer.effectAllowed = 'copyMove'
}

function generateSuffix(name: string, entries: FileEntry[]): string {
  const dot = name.lastIndexOf('.')
  const base = dot > 0 ? name.substring(0, dot) : name
  const ext = dot > 0 ? name.substring(dot) : ''
  let i = 1
  let candidate: string
  do {
    candidate = base + ' (' + i + ')' + ext
    i++
  } while (entries.some(e => e.name === candidate))
  return candidate
}

function addTransfer(name: string, side: 'left' | 'right', dir: 'up' | 'down', transferId: string) {
  transfers.value.push({
    id: transferId, name, size: '--', dir, side, pct: 0, speed: '--',
    done: false, err: '', startedAt: Date.now(), bytesDone: 0, bytesTotal: 0,
  })
}

function cancelTransfer(transferId: string) {
  SFTP.CancelTransfer(transferId)
}

function finishTransfer(id: string, err: string) {
  const t = transfers.value.find(x => x.id === id)
  if (t) {
    t.done = true
    t.err = err
    if (!err) t.pct = 100
  }
  setTimeout(() => {
    transfers.value = transfers.value.filter(x => x.id !== id)
  }, 3000)
}

async function doTransfer(dt: DragTransfer, dest: string, onDone: () => void) {
  const side = dt.targetIsLocal ? 'left' : 'right'
  const transferId = 't' + Date.now() + '_' + Math.random().toString(36).slice(2, 6)
  addTransfer(dt.name, side, dt.fromLocal ? 'up' : 'down', transferId)
  try {
    if (dt.fromLocal && !dt.targetIsLocal) {
      await SFTP.UploadFile(props.sftpSessionId, dt.path, dest, transferId)
    } else if (!dt.fromLocal && dt.targetIsLocal) {
      await SFTP.DownloadFile(props.sftpSessionId, dt.path, dest, transferId)
    }
    finishTransfer(transferId, '')
    showToast('\u4f20\u8f93\u5b8c\u6210: ' + dt.name, 'success')
    onDone()
  } catch (e: any) {
    if ((e?.message ?? '').includes('cancelled')) {
      finishTransfer(transferId, '')
    } else {
      finishTransfer(transferId, e?.message ?? '')
      showToast('\u4f20\u8f93\u5931\u8d25: ' + (e?.message ?? e))
    }
  }
}

function onTransferProgress(payload: any) {
  const p = payload.data
  if (!p) return
  const t = transfers.value.find(x => x.id === p.TransferId)
  if (!t) return
  t.bytesDone = p.BytesDone
  t.bytesTotal = p.BytesTotal
  if (t.bytesTotal > 0) t.pct = Math.round((t.bytesDone / t.bytesTotal) * 100)
  const elapsed = (Date.now() - t.startedAt) / 1000
  if (elapsed > 0 && t.bytesDone > 0) {
    const bps = t.bytesDone / elapsed
    t.speed = formatSize(bps) + '/s'
  }
  if (p.Done) {
    t.pct = 100
    if (p.Err && p.Err !== 'cancelled') t.err = p.Err
  }
}

function resolveConflict(destName: string) {
  if (!conflict.value) return
  const dt = conflict.value
  conflict.value = null
  const target = findPanel(dt.targetIsLocal, dt.targetId)
  if (!target) return
  const dir = target.path.replace(/[\/\\]+$/, '')
  const sep = dt.targetIsLocal ? '\\' : '/'
  const dest = (dir ? dir + sep : '') + destName
  doTransfer(dt, dest, () => loadDir(target))
}

async function onDrop(e: DragEvent, targetPanelId: string, isLocal: boolean) {
  e.preventDefault()
  if (!e.dataTransfer) return
  const raw = e.dataTransfer.getData('text/json')
  if (!raw) return
  try {
    const { path, name, isDir, fromLocal } = JSON.parse(raw)
    if (fromLocal === isLocal) return
    const target = findPanel(isLocal, targetPanelId)
    if (!target) return
    const dir = target.path.replace(/[\/\\]+$/, '')
    const sep = isLocal ? '\\' : '/'
    const dest = (dir ? dir + sep : '') + name

    const existing = target.entries.find(e => e.name === name)
    if (existing) {
      const dt: DragTransfer = { path, name, isDir, fromLocal, targetId: targetPanelId, targetIsLocal: isLocal }
      conflict.value = dt
      conflictNewName.value = generateSuffix(name, target.entries)
      return
    }

    doTransfer(
      { path, name, isDir, fromLocal, targetId: targetPanelId, targetIsLocal: isLocal },
      dest,
      () => loadDir(target),
    )
  } catch (e: any) { showToast('\u4f20\u8f93\u5931\u8d25: ' + (e?.message ?? e)) }
}

const leftTransfers = computed(() => transfers.value.filter(t => t.side === 'left'))
const rightTransfers = computed(() => transfers.value.filter(t => t.side === 'right'))

// path bar

function onPathKeydown(e: KeyboardEvent, isLocal: boolean) {
  if (e.key === 'Enter') {
    e.preventDefault()
    const p = findPanel(isLocal, activePanel(isLocal))
    if (!p) return
    const newPath = isLocal ? pathInputLeft.value : pathInputRight.value
    p.path = newPath
    loadDir(p)
    ;(e.target as HTMLInputElement).blur()
  } else if (e.key === 'Escape') {
    const p = findPanel(isLocal, activePanel(isLocal))
    if (p) syncPathInput(p)
    ;(e.target as HTMLInputElement).blur()
  }
}

function onPathBlur(isLocal: boolean) {
  const p = findPanel(isLocal, activePanel(isLocal))
  if (!p) return
  const val = isLocal ? pathInputLeft.value : pathInputRight.value
  if (val !== p.path) {
    p.path = val
    loadDir(p)
  }
}
</script>

<template>
  <div class="flex flex-col h-full relative" :class="{ hidden: !active }">

    <!-- Toast -->
    <div v-if="toastMsg" class="absolute bottom-4 left-1/2 -translate-x-1/2 z-50 px-4 py-2 rounded-lg text-xs shadow-lg border transition-all"
      :class="toastType === 'error'
        ? 'bg-[var(--danger)]/10 border-[var(--danger)]/30 text-[var(--danger)]'
        : 'bg-[var(--success)]/10 border-[var(--success)]/30 text-[var(--success)]'">
      {{ toastMsg }}
    </div>

    <!-- Conflict dialog -->
    <div v-if="conflict" class="absolute inset-0 z-40 flex items-center justify-center bg-black/40">
      <div class="bg-[var(--bg-elevated)] border border-[var(--border)] rounded-lg shadow-2xl p-5 w-80">
        <div class="text-xs font-semibold mb-3">文件已存在</div>
        <div class="text-xs text-[var(--text-muted)] mb-1">目标位置已存在同名文件：</div>
        <div class="text-xs font-mono bg-[var(--bg-surface)] rounded px-2 py-1.5 mb-4 break-all">{{ conflict.name }}</div>
        <div class="flex gap-2 mb-3">
          <button @click="resolveConflict(conflict.name)"
            class="flex-1 px-3 py-1.5 text-xs rounded bg-[var(--danger)]/10 text-[var(--danger)] border border-[var(--danger)]/20 hover:bg-[var(--danger)]/20 transition-colors cursor-pointer">
            替换
          </button>
          <button @click="resolveConflict(conflictNewName)"
            class="flex-1 px-3 py-1.5 text-xs rounded bg-[var(--accent)]/10 text-[var(--accent)] border border-[var(--accent)]/20 hover:bg-[var(--accent)]/20 transition-colors cursor-pointer">
            保留两者
          </button>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-[11px] text-[var(--text-muted)] flex-shrink-0">新名称:</span>
          <input v-model="conflictNewName" @keyup.enter="resolveConflict(conflictNewName)"
            class="flex-1 px-2 py-1 bg-[var(--bg-surface)] border border-[var(--border)] rounded text-xs focus:outline-none focus:border-[var(--accent)]" />
          <button @click="conflict = null"
            class="px-3 py-1 text-xs rounded border border-[var(--border)] text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors cursor-pointer">
            取消
          </button>
        </div>
      </div>
    </div>

    <div class="h-full flex-1 flex divide-x divide-[var(--border)]">
      <!-- Left side -->
      <div class=" flex-1 flex flex-col min-w-0" @dragover.prevent @drop="(e) => activeLeft && onDrop(e, activeLeft.id, true)"
        @contextmenu="showPanelCtxMenu($event, true)">
        <div class="flex items-center bg-[var(--bg-surface)] border-b border-[var(--border)] h-8 px-1 gap-0.5">
          <button v-for="p in leftPanels" :key="p.id" @click="leftActive = p.id"
            class="flex items-center gap-1 px-3 py-1 text-xs rounded-t cursor-pointer transition-colors"
            :class="leftActive === p.id ? 'bg-[var(--bg-base)] text-[var(--text-primary)]' : 'text-[var(--text-muted)] hover:text-[var(--text-primary)]'">
            {{ p.label }}
            <span v-if="leftPanels.length > 1" @click.stop="closePanel(p.id, true)" class="ml-0.5 hover:text-[var(--danger)] cursor-pointer">&times;</span>
          </button>
        </div>
        <div v-if="activeLeft" class=" flex-1 flex flex-col overflow-hidden relative">
          <div class="flex items-center gap-1 px-2 py-1 bg-[var(--bg-surface)] border-b border-[var(--border-muted)] text-xs">
            <button @click="goUp(activeLeft.id, true)" class="p-0.5 hover:text-[var(--text-primary)] text-[var(--text-muted)] cursor-pointer" title="上级目录">
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
            </button>
            <input v-model="pathInputLeft" @keydown="onPathKeydown($event, true)" @blur="onPathBlur(true)"
              class="flex-1 bg-transparent text-[var(--text-muted)] font-mono text-[11px] border border-transparent focus:border-[var(--border)] rounded px-1.5 py-0.5 focus:outline-none focus:text-[var(--text-primary)]" />
          </div>
          <div class="flex-1 overflow-y-auto text-xs">
            <div v-if="activeLeft.loading" class="p-4 text-[var(--text-muted)]">加载中...</div>
            <div v-for="e in activeLeft.entries" :key="e.name" draggable="true" 
              @dragstart="onDragStart($event, e, activeLeft)"
              @click="selectEntry(true, e)"
              @dblclick="navigate(activeLeft.id, true, e)"
              @contextmenu="showCtxMenu($event, e, activeLeft.id, true)"
              :class="[
                'flex items-center gap-2 px-3 py-1 hover:bg-[var(--bg-hover)] cursor-pointer transition-colors',
                selectedLeft === e.path ? 'bg-[var(--accent)]/10 border-l-2 border-[var(--accent)]' : ''
              ]">
              <svg v-if="e.isDir" class="w-4 h-4 text-[var(--warning)] flex-shrink-0" fill="currentColor" viewBox="0 0 24 24"><path d="M10 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"/></svg>
              <svg v-else class="w-4 h-4 text-[var(--text-muted)] flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"/></svg>
              <span class="flex-1 truncate">{{ e.name }}</span>
              <span v-if="!e.isDir" class="text-[var(--text-muted)] flex-shrink-0">{{ formatSize(e.size) }}</span>
              <span class="text-[var(--text-muted)] flex-shrink-0 w-16 text-right text-[10px]">{{ e.modTime.split(' ')[0] }}</span>
            </div>
          </div>
          <!-- Transfers below left panel -->
          <div v-if="leftTransfers.length" class="absolute bottom-0 left-0 right-0 border-t border-[var(--border)] bg-[var(--bg-surface)] max-h-24 overflow-y-auto z-10">
            <div v-for="t in leftTransfers" :key="t.id"
              class="flex items-center gap-2 px-3 py-1.5 text-[11px]">
              <span class="w-20 truncate flex-shrink-0" :class="t.done ? (t.err ? 'text-[var(--danger)]' : 'text-[var(--success)]') : 'text-[var(--text-primary)]'">{{ t.name }}</span>
              <div class="flex-1 h-1.5 bg-[var(--border)] rounded overflow-hidden">
                <div class="h-full rounded transition-all duration-300" :class="t.err ? 'bg-[var(--danger)]' : 'bg-[var(--accent)]'" :style="{ width: t.pct + '%' }" />
              </div>
              <span class="flex-shrink-0 w-14 text-right text-[var(--text-muted)]">{{ t.speed }}</span>
              <button v-if="!t.done" @click="cancelTransfer(t.id)" class="flex-shrink-0 text-[var(--danger)] hover:text-[var(--danger-hover)] cursor-pointer" title="取消">&times;</button>
              <span v-else-if="t.err" class="flex-shrink-0 text-[var(--danger)] text-[10px]">✗</span>
              <span v-else class="flex-shrink-0 text-[var(--success)]">✓</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right side -->
      <div class="flex-1 flex flex-col min-w-0" @dragover.prevent @drop="(e) => activeRight && onDrop(e, activeRight.id, false)"
        @contextmenu="showPanelCtxMenu($event, false)">
        <div class="flex items-center bg-[var(--bg-surface)] border-b border-[var(--border)] h-8 px-1 gap-0.5">
          <button v-for="p in rightPanels" :key="p.id" @click="rightActive = p.id"
            class="flex items-center gap-1 px-3 py-1 text-xs rounded-t cursor-pointer transition-colors"
            :class="rightActive === p.id ? 'bg-[var(--bg-base)] text-[var(--text-primary)]' : 'text-[var(--text-muted)] hover:text-[var(--text-primary)]'">
            {{ p.label }}
            <span v-if="rightPanels.length > 1" @click.stop="closePanel(p.id, false)" class="ml-0.5 hover:text-[var(--danger)] cursor-pointer">&times;</span>
          </button>
        </div>
        <div v-if="activeRight" class="flex-1 flex flex-col overflow-hidden relative">
          <div class="flex items-center gap-1 px-2 py-1 bg-[var(--bg-surface)] border-b border-[var(--border-muted)] text-xs">
            <button @click="goUp(activeRight.id, false)" class="p-0.5 hover:text-[var(--text-primary)] text-[var(--text-muted)] cursor-pointer" title="上级目录">
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
            </button>
            <input v-model="pathInputRight" @keydown="onPathKeydown($event, false)" @blur="onPathBlur(false)"
              class="flex-1 bg-transparent text-[var(--text-muted)] font-mono text-[11px] border border-transparent focus:border-[var(--border)] rounded px-1.5 py-0.5 focus:outline-none focus:text-[var(--text-primary)]" />
          </div>
          <div class="flex-1 overflow-y-auto text-xs">
            <div v-if="activeRight.loading" class="p-4 text-[var(--text-muted)]">加载中...</div>
            <div v-for="e in activeRight.entries" :key="e.name" draggable="true"
              @dragstart="onDragStart($event, e, activeRight)"
              @click="selectEntry(false, e)"
              @dblclick="navigate(activeRight.id, false, e)"
              @contextmenu="showCtxMenu($event, e, activeRight.id, false)"
              :class="[
                'flex items-center gap-2 px-3 py-1 hover:bg-[var(--bg-hover)] cursor-pointer transition-colors',
                selectedRight === e.path ? 'bg-[var(--accent)]/10 border-l-2 border-[var(--accent)]' : ''
              ]">
              <svg v-if="e.isDir" class="w-4 h-4 text-[var(--warning)] flex-shrink-0" fill="currentColor" viewBox="0 0 24 24"><path d="M10 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"/></svg>
              <svg v-else class="w-4 h-4 text-[var(--text-muted)] flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"/></svg>
              <span class="flex-1 truncate">{{ e.name }}</span>
              <span v-if="!e.isDir" class="text-[var(--text-muted)] flex-shrink-0">{{ formatSize(e.size) }}</span>
              <span class="text-[var(--text-muted)] flex-shrink-0 w-16 text-right text-[10px]">{{ e.modTime.split(' ')[0] }}</span>
            </div>
          </div>
          <!-- Transfers below right panel -->
          <div v-if="rightTransfers.length" class="absolute bottom-0 left-0 right-0 border-t border-[var(--border)] bg-[var(--bg-surface)] max-h-24 overflow-y-auto z-10">
            <div v-for="t in rightTransfers" :key="t.id"
              class="flex items-center gap-2 px-3 py-1.5 text-[11px]">
              <span class="w-20 truncate flex-shrink-0" :class="t.done ? (t.err ? 'text-[var(--danger)]' : 'text-[var(--success)]') : 'text-[var(--text-primary)]'">{{ t.name }}</span>
              <div class="flex-1 h-1.5 bg-[var(--border)] rounded overflow-hidden">
                <div class="h-full rounded transition-all duration-300" :class="t.err ? 'bg-[var(--danger)]' : 'bg-[var(--accent)]'" :style="{ width: t.pct + '%' }" />
              </div>
              <span class="flex-shrink-0 w-14 text-right text-[var(--text-muted)]">{{ t.speed }}</span>
              <button v-if="!t.done" @click="cancelTransfer(t.id)" class="flex-shrink-0 text-[var(--danger)] hover:text-[var(--danger-hover)] cursor-pointer" title="取消">&times;</button>
              <span v-else-if="t.err" class="flex-shrink-0 text-[var(--danger)] text-[10px]">✗</span>
              <span v-else class="flex-shrink-0 text-[var(--success)]">✓</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Context menu -->
    <ContextMenu v-if="ctxMenuVisible" :x="ctxMenuX" :y="ctxMenuY" :items="ctxMenuItems" @close="closeCtxMenu" />
  </div>
</template>
