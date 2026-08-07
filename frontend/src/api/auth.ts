import { client } from './client'
import type { User } from '@/types'

export interface LoginResult {
  token: string
  user: User
}

export function login(username: string, password: string) {
  return client.post<{ data: LoginResult }>('/auth/login', { username, password })
}

export function me() {
  return client.get<{ data: User }>('/auth/me')
}
