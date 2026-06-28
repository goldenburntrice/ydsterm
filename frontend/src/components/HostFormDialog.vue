<script setup lang="ts">
import { reactive, watch, computed, ref } from 'vue'
import type { Host, HostCreateInput, HostUpdateInput, HostGroup, Key } from '../../bindings/ydsterm/internal/types/models'
import { CreateGroup } from '../../bindings/ydsterm/internal/service/hostserviceimpl'

const props = defineProps<{ visible: boolean; host: Host | null; groups: HostGroup[]; keys: Key[] }>()
const emit = defineEmits<{ close: []; save: [data: HostCreateInput | HostUpdateInput]; groupsChanged: [] }>()

const isEdit = computed(() => !!props.host)
const colors = ['#2f81f7', '#f85149', '#3fb950', '#d29922', '#bc8cff', '#39c5cf', '#ff7b72', '#79c0ff']

const form = reactive({
  name: '', hostname: '', port: 22, username: 'root',
  authMethod: 'password', password: '', keyId: '', groupId: '', color: '',
})

const showGroupDropdown = ref(false)
const showNewGroupInput = ref(false)
const newGroupName = ref('')

watch(() => props.visible, (val) => {
  if (!val) return
  showGroupDropdown.value = false
  showNewGroupInput.value = false
  newGroupName.value = ''
  if (props.host) {
    Object.assign(form, {
      name: props.host.name, hostname: props.host.hostname, port: props.host.port,
      username: props.host.username, authMethod: props.host.authMethod,
      password: '', keyId: props.host.keyId, groupId: props.host.groupId, color: props.host.color,
    })
  } else {
    Object.assign(form, {
      name: '', hostname: '', port: 22, username: 'root',
      authMethod: 'password', password: '', keyId: '', groupId: '', color: '',
    })
  }
})

async function handleCreateGroup() {
  const name = newGroupName.value.trim()
  if (!name) return
  await CreateGroup({ name, parentId: '' })
  newGroupName.value = ''
  showNewGroupInput.value = false
  emit('groupsChanged')
}

function selectGroup(id: string) {
  form.groupId = id
  showGroupDropdown.value = false
}

const selectedGroupName = computed(() => {
  if (!form.groupId) return '\u672a\u5206\u7ec4'
  const g = props.groups.find(x => x.id === form.groupId)
  return g ? g.name : '\u672a\u5206\u7ec4'
})

function handleSubmit() {
  if (!form.name || !form.hostname) return
  if (isEdit.value) {
    const data: HostUpdateInput = {
      id: props.host!.id, name: form.name, hostname: form.hostname, port: form.port,
      username: form.username, authMethod: form.authMethod, keyId: form.keyId,
      groupId: form.groupId, color: form.color,
    }
    if (form.authMethod === 'password' && form.password) data.password = form.password
    emit('save', data)
  } else {
    emit('save', {
      name: form.name, hostname: form.hostname, port: form.port, username: form.username,
      authMethod: form.authMethod, password: form.password || '', keyId: form.keyId,
      groupId: form.groupId, color: form.color,
    })
  }
}

const inputClass = 'w-full px-3 py-2 bg-[var(--bg-base)] border border-[var(--border)] rounded-md text-sm focus:outline-none focus:border-[var(--accent)] transition-colors placeholder:text-[var(--text-muted)]'
</script>

<template>
  <Teleport to="body">
    <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="emit('close')" />
      <div class="relative w-full max-w-md mx-4 bg-[var(--bg-elevated)] border border-[var(--border)] rounded-xl shadow-2xl">
        <div class="flex items-center justify-between px-5 py-3.5 border-b border-[var(--border)]">
          <h2 class="text-sm font-semibold">{{ isEdit ? '编辑主机' : '新建主机' }}</h2>
          <button @click="emit('close')" class="p-1 rounded hover:bg-[var(--bg-hover)] text-[var(--text-muted)] hover:text-[var(--text-primary)] transition-colors cursor-pointer">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>

        <form @submit.prevent="handleSubmit" class="p-5 space-y-4">
          <div class="flex items-center gap-1.5">
            <div v-for="c in colors" :key="c" @click="form.color = c"
              class="w-6 h-6 rounded-full cursor-pointer border-2 transition-all hover:scale-110"
              :class="form.color === c ? 'border-white scale-110 ring-2 ring-white/20' : 'border-transparent'"
              :style="{ backgroundColor: c }" />
          </div>

          <div>
            <label class="block text-[11px] font-medium text-[var(--text-muted)] mb-1.5 uppercase tracking-wide">名称</label>
            <input v-model="form.name" type="text" required :class="inputClass" placeholder="生产服务器" />
          </div>
          <div class="grid grid-cols-3 gap-3">
            <div class="col-span-2">
              <label class="block text-[11px] font-medium text-[var(--text-muted)] mb-1.5 uppercase tracking-wide">地址</label>
              <input v-model="form.hostname" type="text" required :class="inputClass" placeholder="192.168.1.100" />
            </div>
            <div>
              <label class="block text-[11px] font-medium text-[var(--text-muted)] mb-1.5 uppercase tracking-wide">端口</label>
              <input v-model.number="form.port" type="number" :class="inputClass" />
            </div>
          </div>
          <div>
            <label class="block text-[11px] font-medium text-[var(--text-muted)] mb-1.5 uppercase tracking-wide">用户名</label>
            <input v-model="form.username" type="text" :class="inputClass" placeholder="root" />
          </div>

          <div>
            <label class="block text-[11px] font-medium text-[var(--text-muted)] mb-1.5 uppercase tracking-wide">认证方式</label>
            <div class="flex gap-1.5">
              <button type="button" @click="form.authMethod = 'password'"
                class="flex-1 py-2 rounded-md text-xs font-medium transition-colors cursor-pointer"
                :class="form.authMethod === 'password' ? 'bg-[var(--accent)] text-white' : 'bg-[var(--bg-base)] text-[var(--text-secondary)] border border-[var(--border)]'">密码</button>
              <button type="button" @click="form.authMethod = 'private_key'"
                class="flex-1 py-2 rounded-md text-xs font-medium transition-colors cursor-pointer"
                :class="form.authMethod === 'private_key' ? 'bg-[var(--accent)] text-white' : 'bg-[var(--bg-base)] text-[var(--text-secondary)] border border-[var(--border)]'">密钥</button>
            </div>
          </div>

          <div v-if="form.authMethod === 'password'">
            <input v-model="form.password" type="password" :class="inputClass" :placeholder="isEdit ? '留空保持不变' : '输入 SSH 密码'" />
          </div>
          <div v-else>
            <select v-model="form.keyId" :class="inputClass">
              <option value="">-- 选择密钥 --</option>
              <option v-for="k in keys" :key="k.id" :value="k.id">{{ k.name }}</option>
            </select>
          </div>

          <div class="relative">
            <label class="block text-[11px] font-medium text-[var(--text-muted)] mb-1.5 uppercase tracking-wide">分组</label>
            <button type="button" @click="showGroupDropdown = !showGroupDropdown"
              class="w-full flex items-center justify-between px-3 py-2 bg-[var(--bg-base)] border border-[var(--border)] rounded-md text-sm hover:border-[var(--text-muted)] transition-colors cursor-pointer">
              <span :class="form.groupId ? 'text-[var(--text-primary)]' : 'text-[var(--text-muted)]'">{{ selectedGroupName }}</span>
              <svg class="w-3.5 h-3.5 text-[var(--text-muted)]" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
            </button>
            <div v-if="showGroupDropdown"
              class="absolute top-full left-0 right-0 mt-1 bg-[var(--bg-elevated)] border border-[var(--border)] rounded-md shadow-lg z-10 max-h-48 overflow-y-auto">
              <button type="button" @click="selectGroup('')"
                class="w-full flex items-center px-3 py-2 text-xs hover:bg-[var(--bg-hover)] transition-colors cursor-pointer"
                :class="!form.groupId ? 'text-[var(--accent)]' : 'text-[var(--text-primary)]'">
                <svg class="w-3 h-3 mr-2 opacity-40" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
                未分组
              </button>
              <button type="button" v-for="g in groups" :key="g.id" @click="selectGroup(g.id)"
                class="w-full flex items-center px-3 py-2 text-xs hover:bg-[var(--bg-hover)] transition-colors cursor-pointer"
                :class="form.groupId === g.id ? 'text-[var(--accent)]' : 'text-[var(--text-primary)]'">
                <svg class="w-3 h-3 mr-2 opacity-40" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
                {{ g.name }}
              </button>
              <div class="border-t border-[var(--border)]">
                <div v-if="showNewGroupInput" class="flex items-center gap-1 px-2 py-1.5">
                  <input v-model="newGroupName" @keyup.enter="handleCreateGroup" @keyup.escape="showNewGroupInput = false"
                    class="flex-1 px-2 py-1 bg-[var(--bg-base)] border border-[var(--accent)] rounded text-xs focus:outline-none"
                    placeholder="分组名称..." />
                  <button type="button" @click="handleCreateGroup"
                    class="px-2 py-1 bg-[var(--accent)] text-white rounded text-[10px] hover:opacity-90 cursor-pointer">确定</button>
                </div>
                <button v-else type="button" @click="showNewGroupInput = true"
                  class="w-full flex items-center px-3 py-2 text-xs text-[var(--accent)] hover:bg-[var(--bg-hover)] transition-colors cursor-pointer">
                  <svg class="w-3 h-3 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
                  新建分组
                </button>
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <button type="button" @click="emit('close')"
              class="px-4 py-2 rounded-md text-xs bg-[var(--bg-base)] border border-[var(--border)] text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors cursor-pointer">取消</button>
            <button type="submit"
              class="px-4 py-2 rounded-md text-xs bg-[var(--accent)] text-white hover:opacity-90 transition-opacity font-medium cursor-pointer">{{ isEdit ? '保存' : '创建' }}</button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>
