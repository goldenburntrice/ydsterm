<script setup lang="ts">
import { ref } from 'vue'
import type { Key } from '../../bindings/ydsterm/internal/types/models'

defineProps<{ visible: boolean; keys: Key[] }>()
const emit = defineEmits<{ close: []; import: [data: { name: string; privateKey: string; passphrase?: string }]; deleteKey: [id: string] }>()

const importName = ref('')
const importKey = ref('')
const importPass = ref('')

function doImport() {
  if (!importName.value || !importKey.value) return
  emit('import', { name: importName.value, privateKey: importKey.value, passphrase: importPass.value || undefined })
  importName.value = ''; importKey.value = ''; importPass.value = ''
}

const inputClass = 'w-full px-3 py-2 bg-[var(--bg-base)] border border-[var(--border)] rounded-md text-[13px] focus:outline-none focus:border-[var(--accent)] transition-colors placeholder:text-[var(--text-muted)]'
</script>

<template>
  <Teleport to="body">
    <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="emit('close')" />
      <div class="relative w-full max-w-md mx-4 bg-[var(--bg-elevated)] border border-[var(--border)] rounded-xl shadow-2xl">
        <div class="flex items-center justify-between px-5 py-3.5 border-b border-[var(--border)]">
          <h2 class="text-sm font-semibold">SSH 密钥管理</h2>
          <button @click="emit('close')" class="p-1 rounded hover:bg-[var(--bg-hover)] text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors cursor-pointer">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>

        <div class="p-5 space-y-4">
          <div class="p-3 rounded-lg bg-[var(--bg-base)] border border-[var(--border)]">
            <label class="block text-[11px] font-medium text-[var(--text-muted)] mb-2">导入私钥</label>
            <input v-model="importName" type="text" :class="inputClass + ' mb-2'" placeholder="密钥名称" />
            <textarea v-model="importKey" rows="5" :class="inputClass + ' font-mono text-[11px] resize-none'" placeholder="-----BEGIN OPENSSH PRIVATE KEY-----&#10;...&#10;-----END OPENSSH PRIVATE KEY-----" />
            <input v-model="importPass" type="password" :class="inputClass + ' mt-2'" placeholder="Passphrase（可选）" />
            <button @click="doImport" class="mt-2 px-4 py-1.5 rounded-md text-xs bg-[var(--accent)] text-white hover:opacity-90 transition-opacity font-medium cursor-pointer">导入密钥</button>
          </div>

          <div>
            <div class="text-[11px] font-medium text-[var(--text-muted)] mb-2">已有密钥 ({{ keys.length }})</div>
            <div v-if="keys.length === 0" class="text-xs text-[var(--text-muted)] py-2">暂无密钥</div>
            <div v-for="k in keys" :key="k.id" class="flex items-center justify-between px-3 py-2 rounded-md bg-[var(--bg-base)] mb-1 border border-[var(--border)]">
              <div class="flex-1 min-w-0">
                <div class="text-[13px] font-medium truncate">{{ k.name }}</div>
                <div v-if="k.publicKey" class="text-[10px] text-[var(--text-muted)] truncate font-mono mt-0.5">{{ k.publicKey.split(' ').slice(0, 2).join(' ') }}...</div>
              </div>
              <button @click="emit('deleteKey', k.id)" class="ml-2 p-1 rounded hover:bg-[var(--danger)]/15 text-[var(--text-muted)] hover:text-[var(--danger)] transition-colors cursor-pointer">
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
