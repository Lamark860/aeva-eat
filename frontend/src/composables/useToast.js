import { ref } from 'vue'

const toasts = ref([])
let nextId = 0

export function useToast() {
  function show(message, type = 'success', duration = 3000) {
    const id = ++nextId
    toasts.value.push({ id, message, type })
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, duration)
  }

  function success(message) { show(message, 'success') }
  function error(message)   { show(message, 'error') }
  function info(message)    { show(message, 'info') }
  function warning(message) { show(message, 'warning') }

  return { toasts, show, success, error, info, warning }
}
