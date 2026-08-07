import { client } from './client'
import type { Settings } from '@/types'

export interface SettingsInput {
  shop_name: string
  address: string
  phone: string
  receipt_footer: string
  tax_percentage: number
}

export const getSettings = () => client.get<{ data: Settings }>('/settings')
export const updateSettings = (input: SettingsInput) => client.put<{ data: Settings }>('/settings', input)
