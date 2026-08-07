import { client } from './client'
import type { StockMovement } from '@/types'

export const listMovements = (limit = 0) =>
  client.get<{ data: StockMovement[] }>('/inventory/movements', { params: limit ? { limit } : {} })
export const listMovementsByProduct = (productId: number) =>
  client.get<{ data: StockMovement[] }>(`/inventory/movements/product/${productId}`)
export const createAdjustment = (input: { product_id: number; quantity: number; note: string }) =>
  client.post<{ data: StockMovement }>('/inventory/adjustments', input)
