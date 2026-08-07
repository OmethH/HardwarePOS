import { client } from './client'
import type { ReturnRecord } from '@/types'

export interface ReturnItemInput {
  sale_item_id: number
  quantity: number
}

export interface ReturnInput {
  sale_id: number
  items: ReturnItemInput[]
}

export const listReturns = (limit = 0) =>
  client.get<{ data: ReturnRecord[] }>('/returns', { params: limit ? { limit } : {} })
export const getReturn = (id: number) => client.get<{ data: ReturnRecord }>(`/returns/${id}`)
export const createReturn = (input: ReturnInput) => client.post<{ data: ReturnRecord }>('/returns', input)
