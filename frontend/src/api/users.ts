import { client } from './client'
import type { Role, User } from '@/types'

export interface CreateUserInput {
  name: string
  username: string
  password: string
  role_id: number
}

export interface UpdateUserInput {
  name: string
  role_id: number
  active: boolean
  password?: string
}

export const listUsers = () => client.get<{ data: User[] }>('/users')
export const getUser = (id: number) => client.get<{ data: User }>(`/users/${id}`)
export const createUser = (input: CreateUserInput) => client.post<{ data: User }>('/users', input)
export const updateUser = (id: number, input: UpdateUserInput) =>
  client.put<{ data: User }>(`/users/${id}`, input)
export const deactivateUser = (id: number) => client.delete(`/users/${id}`)
export const listRoles = () => client.get<{ data: Role[] }>('/roles')
