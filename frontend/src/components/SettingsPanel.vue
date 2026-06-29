<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as SettingsService from '../../bindings/ydsterm/internal/service/settingsserviceimpl'
import * as UserService from '../../bindings/ydsterm/internal/service/userserviceimpl'
import { useSyncProgress } from '../composables/useSyncProgress'

const emit = defineEmits<{ close: []; refresh: [] }>()

const { syncing, currentModule, syncError } = useSyncProgress()

const activeTab = ref<'general' | 'cloud'>('general')

const serverAddr = ref('')
const serverKey = ref('')
const username = ref('')
const password = ref('')
const verifying = ref(false)
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null)
const savedUsers = ref<{ id: string; username: string; isActive: boolean }[]>([])
const currentUser = ref<string>('')

async function loadSettings() {
  try {
    serverAddr.value = await SettingsService.Get('server_addr') ?? ''
    serverKey.value = await SettingsService.Get('server_key') ?? ''
    await loadUsers()
  } catch {
  }
}

async function loadUsers() {
  try {
    const users = await UserService.ListUsers()
    savedUsers.value = users ?? []
    const current = await UserService.GetCurrentUser()
    currentUser.value = current?.username ?? 'LOCALUSER'
  } catch {
    savedUsers.value = []
  }
}

async function saveServerAddr() {
  await SettingsService.Set('server_addr', serverAddr.value)
}

async function saveServerKey() {
  await SettingsService.Set('server_key', serverKey.value)
}

async function handleVerify() {
  if (!serverAddr.value || !serverKey.value) {
    message.value = { type: 'error', text: '请先配置服务器地址和密钥' }
    return
  }
  if (!username.value || !password.value) {
    message.value = { type: 'error', text: '请输入用户名和密码' }
    return
  }

  verifying.value = true
  message.value = null
  try {
    const result = await SettingsService.VerifyUser(serverAddr.value, serverKey.value, username.value, password.value)
    await UserService.RegisterVerifiedUser(username.value, result?.passwordHash ?? '')
    message.value = { type: 'success', text: '验证成功，正在同步数据...' }
    await loadUsers()
    setTimeout(() => {
      emit('refresh')
      emit('close')
    }, 1500)
  } catch (e: any) {
    message.value = { type: 'error', text: e?.message ?? '验证失败' }
  } finally {
    verifying.value = false
  }
}

async function handleSwitchUser(userId: string) {
  try {
    await UserService.SwitchUser(userId)
    await loadUsers()
    message.value = { type: 'success', text: '切换成功，正在同步数据...' }
    setTimeout(() => {
      emit('refresh')
      emit('close')
    }, 1500)
  } catch (e: any) {
    message.value = { type: 'error', text: e?.message ?? '切换失败' }
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
}

onMounted(() => {
  loadSettings()
})
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="emit('close')" @keydown="handleKeydown">
    <div class="w-[640px] h-[560px] bg-[var(--bg-surface)] border border-[var(--border)] rounded-xl shadow-2xl flex overflow-hidden" @click.stop>
      <div class="w-[180px] border-r border-[var(--border)] py-4 px-2 flex flex-col gap-1">
        <button
          @click="activeTab = 'general'"
          class="w-full text-left px-3 py-2 rounded-lg text-sm transition-colors cursor-pointer"
          :class="activeTab === 'general' ? 'bg-[var(--accent)]/10 text-[var(--accent)] font-medium' : 'text-[var(--text-secondary)] hover:bg-[var(--bg-hover)]'"
        >
          通用
        </button>
        <button
          @click="activeTab = 'cloud'"
          class="w-full text-left px-3 py-2 rounded-lg text-sm transition-colors cursor-pointer"
          :class="activeTab === 'cloud' ? 'bg-[var(--accent)]/10 text-[var(--accent)] font-medium' : 'text-[var(--text-secondary)] hover:bg-[var(--bg-hover)]'"
        >
          云配置
        </button>
      </div>
      <div class="flex-1 flex flex-col">
        <div class="flex items-center justify-between px-5 py-3 border-b border-[var(--border)]">
          <h2 class="text-sm font-medium text-[var(--text-primary)]">
            {{ activeTab === 'general' ? '通用设置' : '云配置' }}
          </h2>
          <button @click="emit('close')" class="p-1 rounded-lg hover:bg-[var(--bg-hover)] text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors cursor-pointer">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>
        <div class="flex-1 overflow-y-auto p-5">
          <div v-if="activeTab === 'general'">
            <div class="text-sm text-[var(--text-secondary)]">
              <div class="flex items-center justify-between py-2">
                <span>版本号</span>
                <span class="font-mono text-[var(--text-muted)]">{{ 'dev' }}</span>
              </div>
              <div class="flex items-center justify-between py-2">
                <span>当前用户</span>
                <span class="font-mono text-[var(--accent)]">{{ currentUser }}</span>
              </div>
            </div>
          </div>
          <div v-if="activeTab === 'cloud'" class="flex flex-col gap-4">
            <div v-if="syncing" class="flex items-center gap-2 px-3 py-2 bg-blue-500/10 border border-blue-500/30 rounded-lg">
              <svg class="w-4 h-4 animate-spin text-blue-500" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span class="text-sm text-blue-500">正在同步 {{ currentModule }}...</span>
            </div>
            <div v-if="syncError" class="flex items-center gap-2 px-3 py-2 bg-red-500/10 border border-red-500/30 rounded-lg">
              <svg class="w-4 h-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <span class="text-sm text-red-500">{{ syncError }}</span>
            </div>
            <div>
              <label class="block text-xs text-[var(--text-muted)] mb-1.5">服务器地址</label>
              <input
                v-model="serverAddr"
                type="text"
                placeholder="http://localhost:7788"
                @blur="saveServerAddr"
                class="w-full h-9 px-3 bg-[var(--bg-base)] border border-[var(--border)] rounded-lg text-sm focus:outline-none focus:border-[var(--accent)] transition-colors placeholder:text-[var(--text-muted)]"
              />
            </div>
            <div>
              <label class="block text-xs text-[var(--text-muted)] mb-1.5">服务器密钥</label>
              <input
                v-model="serverKey"
                type="password"
                placeholder="输入服务器密钥"
                @blur="saveServerKey"
                class="w-full h-9 px-3 bg-[var(--bg-base)] border border-[var(--border)] rounded-lg text-sm focus:outline-none focus:border-[var(--accent)] transition-colors placeholder:text-[var(--text-muted)]"
              />
            </div>
            <div class="border-t border-[var(--border)] pt-4">
              <div class="text-xs text-[var(--text-muted)] mb-3">已保存的账户</div>
              <div class="flex flex-col gap-2 mb-4">
                <div
                  v-for="user in savedUsers"
                  :key="user.id"
                  class="flex items-center justify-between px-3 py-2 rounded-lg transition-colors"
                  :class="user.isActive ? 'bg-[var(--accent)]/10 border border-[var(--accent)]/30' : 'bg-[var(--bg-base)] border border-[var(--border)] hover:border-[var(--accent)]/50'"
                >
                  <div class="flex items-center gap-2">
                    <svg v-if="user.isActive" class="w-4 h-4 text-[var(--accent)]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                    <span class="text-sm" :class="user.isActive ? 'text-[var(--accent)] font-medium' : 'text-[var(--text-primary)]'">{{ user.username }}</span>
                  </div>
                  <button
                    v-if="!user.isActive"
                    @click="handleSwitchUser(user.id)"
                    class="px-2 py-1 text-xs rounded bg-[var(--bg-hover)] hover:bg-[var(--accent)] hover:text-white transition-colors cursor-pointer"
                  >
                    切换
                  </button>
                  <span v-else class="text-xs text-[var(--accent)]">当前</span>
                </div>
              </div>
            </div>
            <div class="border-t border-[var(--border)] pt-4">
              <div class="text-xs text-[var(--text-muted)] mb-3">添加新账户</div>
              <div class="flex flex-col gap-3">
                <input
                  v-model="username"
                  type="text"
                  placeholder="用户名"
                  class="w-full h-9 px-3 bg-[var(--bg-base)] border border-[var(--border)] rounded-lg text-sm focus:outline-none focus:border-[var(--accent)] transition-colors placeholder:text-[var(--text-muted)]"
                />
                <input
                  v-model="password"
                  type="password"
                  placeholder="密码"
                  @keydown.enter="handleVerify"
                  class="w-full h-9 px-3 bg-[var(--bg-base)] border border-[var(--border)] rounded-lg text-sm focus:outline-none focus:border-[var(--accent)] transition-colors placeholder:text-[var(--text-muted)]"
                />
                <button
                  @click="handleVerify"
                  :disabled="verifying"
                  class="w-full h-9 bg-[var(--accent)] text-white rounded-lg text-sm font-medium hover:opacity-90 transition-opacity cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ verifying ? '验证中...' : '验证并登录' }}
                </button>
                <div v-if="message" class="text-xs text-center py-1" :class="message.type === 'success' ? 'text-emerald-500' : 'text-[var(--danger)]'">
                  {{ message.text }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
