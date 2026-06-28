<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import { Events } from '@wailsio/runtime'
import '@xterm/xterm/css/xterm.css'
import { Write, Resize } from '../../bindings/ydsterm/internal/service/terminalserviceimpl'
import { useTheme, xtermDark, xtermLight } from '../composables/useTheme'

const props = defineProps<{
  tabId: string; sessionId: string; hostName: string; hostId: string; active: boolean
}>()

const { isDark } = useTheme()
const containerEl = ref<HTMLDivElement | null>(null)
let term: Terminal | null = null
let fitAddon: FitAddon | null = null
let ro: ResizeObserver | null = null
const connected = ref(true)

function getXtermTheme() {
  return isDark() ? xtermDark : xtermLight
}

function doFit() {
  if (!term || !fitAddon) return
  try { fitAddon.fit() } catch (_) {}
  term.scrollToBottom()
  if (connected.value && props.sessionId) {
    Resize(props.sessionId, term.cols, term.rows)
  }
}

function initTerminal() {
  const el = containerEl.value
  if (!el) return

  term = new Terminal({
    cursorBlink: true,
    cursorStyle: 'bar',
    fontSize: 14,
    fontFamily: '"JetBrains Mono", "Cascadia Code", Consolas, monospace',
    lineHeight: 1.4,
    allowTransparency: true,
    allowProposedApi: true,
    theme: getXtermTheme(),
  })

  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.loadAddon(new WebLinksAddon())
  term.open(el)
  doFit()
  term.focus()

  el.addEventListener('click', () => term?.focus())
  el.addEventListener('mousedown', () => term?.focus())

  term.onData(data => {
    if (connected.value && props.sessionId) Write(props.sessionId, data)
  })
  term.onResize(({ cols, rows }) => {
    if (connected.value && props.sessionId) Resize(props.sessionId, cols, rows)
  })

  ro = new ResizeObserver(() => {
    requestAnimationFrame(() => doFit())
  })
  ro.observe(el)
}

const unwatchTheme = watch(isDark, () => {
  if (term) {
    term.options.theme = getXtermTheme()
    if (containerEl.value) {
      const vp = containerEl.value.querySelector('.xterm-viewport') as HTMLElement
      if (vp) vp.style.backgroundColor = getXtermTheme().background
    }
  }
})

onMounted(() => {
  initTerminal()

  const onTermOutput = (payload: any) => {
    const inner = payload.data
    if (!inner || inner.SessionId !== props.sessionId) return
    if (inner.Data) {
      const binary = atob(inner.Data)
      const bytes = new Uint8Array(binary.length)
      for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i)
      term?.write(bytes)
      term?.scrollToBottom()
    }
  }
  const onTermDisconnected = (payload: any) => {
    const inner = payload.data
    if (inner?.SessionId !== props.sessionId) return
    connected.value = false
    term?.write('\r\n\x1b[31m[连接已断开]\x1b[0m\r\n')
  }

  Events.On('term-output', onTermOutput)
  Events.On('term-disconnected', onTermDisconnected)

  onBeforeUnmount(() => {
    unwatchTheme()
    ro?.disconnect()
    term?.dispose()
  })
})

watch(() => props.active, (isActive) => {
  if (isActive) {
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {
        doFit()
        term?.focus()
      })
    })
  }
})
</script>

<template>
  <div class="flex flex-col h-full" :class="{ hidden: !active }">
    <div ref="containerEl" class="h-full" />
  </div>
</template>
