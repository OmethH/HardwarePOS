import axios from 'axios'

// The Vite dev server proxies /api to localhost:8971 (see vite.config.ts),
// but in production/Vercel hosting it connects to VITE_API_URL if configured,
// falling back to http://localhost:8971/api for the desktop Tauri sidecar.
const apiBase = (import.meta.env.VITE_API_URL as string) || (import.meta.env.DEV ? '/api' : 'http://localhost:8971/api')

export const client = axios.create({
  baseURL: apiBase,
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('hardwarepos_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('hardwarepos_token')
      localStorage.removeItem('hardwarepos_user')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)

export function apiErrorMessage(err: unknown, fallback = 'Something went wrong'): string {
  if (axios.isAxiosError(err)) {
    return err.response?.data?.error ?? fallback
  }
  return fallback
}
