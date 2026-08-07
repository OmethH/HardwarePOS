import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '@/api/auth'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('hardwarepos_token'))
  const user = ref<User | null>(JSON.parse(localStorage.getItem('hardwarepos_user') || 'null'))

  const isAuthenticated = computed(() => !!token.value)
  const role = computed(() => user.value?.role.name ?? '')

  function persist() {
    if (token.value) localStorage.setItem('hardwarepos_token', token.value)
    else localStorage.removeItem('hardwarepos_token')

    if (user.value) localStorage.setItem('hardwarepos_user', JSON.stringify(user.value))
    else localStorage.removeItem('hardwarepos_user')
  }

  async function login(username: string, password: string) {
    const res = await authApi.login(username, password)
    token.value = res.data.data.token
    user.value = res.data.data.user
    persist()
  }

  function logout() {
    token.value = null
    user.value = null
    persist()
  }

  function hasRole(...roles: string[]) {
    return !!user.value && roles.includes(user.value.role.name)
  }

  return { token, user, isAuthenticated, role, login, logout, hasRole }
})
