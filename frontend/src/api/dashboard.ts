import { client } from './client'
import type { DashboardSummary } from '@/types'

export const getSummary = () => client.get<{ data: DashboardSummary }>('/dashboard/summary')

export const getSalesOverTime = (days = 30) =>
  client.get<{ data: { date: string; total: number }[] }>('/dashboard/sales-over-time', { params: { days } })

export const getProfitTrend = (days = 30) =>
  client.get<{ data: { date: string; profit: number }[] }>('/dashboard/profit-trend', { params: { days } })

export const getTopProducts = (days = 30, limit = 10) =>
  client.get<{
    data: { product_id: number; product_name: string; quantity_sold: number; revenue: number }[]
  }>('/dashboard/top-products', { params: { days, limit } })
