import axios from 'axios'

const http = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Background-poll endpoints — упасть с 401 для них нормально, перенаправлять
// юзера на /login из-за фонового запроса = плохо (он мог залогиниваться,
// читать /login или просто открыть страницу). Тихо отклоняем.
const SILENT_401_PATTERNS = [
  '/feed/seen',
  '/feed/unread-count',
]

http.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const url = error.config?.url || ''
      const isSilent = SILENT_401_PATTERNS.some((p) => url.includes(p))
      const isAuthEndpoint = url.includes('/invites/validate/') || url.includes('/auth/')
      const alreadyOnLogin = window.location.pathname === '/login' ||
                             window.location.pathname.startsWith('/invite/')

      if (!isSilent && !isAuthEndpoint && !alreadyOnLogin) {
        localStorage.removeItem('token')
        window.location.href = '/login'
      } else if (isSilent && !alreadyOnLogin) {
        // Фоновый poll упал — токен невалиден, чистим. Но не делаем
        // hard-reload: пусть текущий рендер продолжается, юзер увидит
        // 401 при следующем активном действии и тогда уйдёт на /login.
        localStorage.removeItem('token')
      }
    }
    return Promise.reject(error)
  }
)

export default http
