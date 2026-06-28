import { ref } from 'vue'

export function useModal() {
  const visible = ref(false)

  function open() { visible.value = true }
  function close() { visible.value = false }

  return { visible, open, close }
}
