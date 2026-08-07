import { client } from './client'
import type { Product } from '@/types'

export interface ProductInput {
  sku: string
  barcode: string | null
  name: string
  category_id: number | null
  brand: string
  description: string
  unit_type: string
  purchase_price: number
  selling_price: number
  minimum_stock: number
  status: string
}

export interface ProductListParams {
  search?: string
  category_id?: number
  low_stock?: boolean
  status?: string
}

export const listProducts = (params: ProductListParams = {}) =>
  client.get<{ data: Product[] }>('/products', { params })
export const getProduct = (id: number) => client.get<{ data: Product }>(`/products/${id}`)
export const getProductByBarcode = (barcode: string) =>
  client.get<{ data: Product }>(`/products/barcode/${barcode}`)
export const createProduct = (input: ProductInput) => client.post<{ data: Product }>('/products', input)
export const updateProduct = (id: number, input: ProductInput) =>
  client.put<{ data: Product }>(`/products/${id}`, input)
export const deleteProduct = (id: number) => client.delete(`/products/${id}`)
