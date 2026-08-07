import { client } from './client'
import type { Supplier } from '@/types'

export interface SupplierInput {
  name: string
  phone: string
  email: string
  address: string
  company: string
}

export const listSuppliers = (search = '') =>
  client.get<{ data: Supplier[] }>('/suppliers', { params: search ? { search } : {} })
export const getSupplier = (id: number) => client.get<{ data: Supplier }>(`/suppliers/${id}`)
export const createSupplier = (input: SupplierInput) => client.post<{ data: Supplier }>('/suppliers', input)
export const updateSupplier = (id: number, input: SupplierInput) =>
  client.put<{ data: Supplier }>(`/suppliers/${id}`, input)
export const deleteSupplier = (id: number) => client.delete(`/suppliers/${id}`)
