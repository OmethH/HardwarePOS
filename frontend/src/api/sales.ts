import { client } from './client'
import type { Sale } from '@/types'

export interface SaleItemInput {
  product_id: number
  quantity: number
}

export interface CheckoutInput {
  customer_id: number
  discount: number
  payment_method: string
  items: SaleItemInput[]
}

export interface SaleListParams {
  search?: string
  from?: string
  to?: string
  customer_id?: number
}

export const checkout = (input: CheckoutInput) => client.post<{ data: Sale }>('/sales', input)
export const listSales = (params: SaleListParams = {}) => client.get<{ data: Sale[] }>('/sales', { params })
export const getSale = (id: number) => client.get<{ data: Sale }>(`/sales/${id}`)
export const getSaleByInvoice = (invoice: string) => client.get<{ data: Sale }>(`/sales/invoice/${invoice}`)
