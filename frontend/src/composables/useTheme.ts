import { ref } from 'vue'

const THEME_KEY = 'ydsterm-theme'

const theme = ref<string>((typeof document !== 'undefined' ? document.documentElement.dataset.theme : null) || 'dark')

function applyTheme(t: string) {
  theme.value = t
  if (typeof document !== 'undefined') {
    document.documentElement.dataset.theme = t
    try { localStorage.setItem(THEME_KEY, t) } catch (_) {}
  }
}

// restore from localStorage on init
if (typeof document !== 'undefined') {
  const saved = localStorage.getItem(THEME_KEY)
  if (saved === 'light' || saved === 'dark') applyTheme(saved)
}

function toggleTheme() {
  applyTheme(theme.value === 'dark' ? 'light' : 'dark')
}

export interface ThemeState {
  theme: Readonly<typeof theme>
  toggleTheme: () => void
  isDark: () => boolean
  isLight: () => boolean
}

// singleton
let inst: ThemeState | null = null

export function useTheme(): ThemeState {
  if (inst) return inst
  inst = {
    theme: theme as Readonly<typeof theme>,
    toggleTheme,
    isDark: () => theme.value === 'dark',
    isLight: () => theme.value === 'light',
  }
  return inst
}

// Light theme xterm colors (export for reuse)
export const xtermDark = {
  background: '#0d1117',
  foreground: '#e6edf3',
  cursor: '#2f81f7',
  selectionBackground: '#2f81f740',
  black: '#484f58', red: '#ff7b72', green: '#3fb950', yellow: '#d29922',
  blue: '#58a6ff', magenta: '#bc8cff', cyan: '#39c5cf', white: '#b1bac4',
  brightBlack: '#6e7681', brightRed: '#ffa198', brightGreen: '#56d364',
  brightYellow: '#e3b341', brightBlue: '#79c0ff', brightMagenta: '#d2a8ff',
  brightCyan: '#56d4dd', brightWhite: '#f0f6fc',
}

export const xtermLight = {
  background: '#ffffff',
  foreground: '#1f2328',
  cursor: '#0969da',
  selectionBackground: '#0969da30',
  black: '#d0d7de', red: '#cf222e', green: '#1a7f37', yellow: '#9a6700',
  blue: '#0969da', magenta: '#8250df', cyan: '#1b7c83', white: '#656d76',
  brightBlack: '#8b949e', brightRed: '#a40e26', brightGreen: '#116329',
  brightYellow: '#633c01', brightBlue: '#0550ae', brightMagenta: '#572398',
  brightCyan: '#0e626b', brightWhite: '#1f2328',
}
