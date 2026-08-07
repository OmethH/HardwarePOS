import { client } from './client'
import type { Purchase } from '@/types'

export interface PurchaseItemInput {
  product_id: number
  quantity: number
  price: number
}

export interface PurchaseInput {
  supplier_id: number
  items: PurchaseItemInput[]
}

export const listPurchases = (limit = 0) =>
  client.get<{ data: Purchase[] }>('/purchases', { params: limit ? { limit } : {} })
export const getPurchase = (id: number) => client.get<{ data: Purchase }>(`/purchases/${id}`)
export const createPurchase = (input: PurchaseInput) => client.post<{ data: Purchase }>('/purchases', input)
