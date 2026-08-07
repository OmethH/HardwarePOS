import { client } from './client'
import type { Customer } from '@/types'

export interface CustomerInput {
  name: string
  phone: string
  address: string
}

export const listCustomers = (search = '') =>
  client.get<{ data: Customer[] }>('/customers', { params: search ? { search } : {} })
export const getCustomer = (id: number) => client.get<{ data: Customer }>(`/customers/${id}`)
export const createCustomer = (input: CustomerInput) => client.post<{ data: Customer }>('/customers', input)
export const updateCustomer = (id: number, input: CustomerInput) =>
  client.put<{ data: Customer }>(`/customers/${id}`, input)
export const deleteCustomer = (id: number) => client.delete(`/customers/${id}`)
