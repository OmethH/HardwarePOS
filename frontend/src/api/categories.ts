import { client } from './client'
import type { Category } from '@/types'

export const listCategories = () => client.get<{ data: Category[] }>('/categories')
export const getCategory = (id: number) => client.get<{ data: Category }>(`/categories/${id}`)
export const createCategory = (name: string) => client.post<{ data: Category }>('/categories', { name })
export const updateCategory = (id: number, name: string) =>
  client.put<{ data: Category }>(`/categories/${id}`, { name })
export const deleteCategory = (id: number) => client.delete(`/categories/${id}`)
