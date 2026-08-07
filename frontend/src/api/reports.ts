import { client } from './client'
import type { InventoryReportRow, ProductPerformanceRow, SalesReportRow } from '@/types'

export const getSalesReport = (from = '', to = '') =>
  client.get<{ data: SalesReportRow[] }>('/reports/sales', { params: { from, to } })

export const getInventoryReport = () => client.get<{ data: InventoryReportRow[] }>('/reports/inventory')

export const getProductPerformance = (from = '', to = '') =>
  client.get<{ data: ProductPerformanceRow[] }>('/reports/product-performance', { params: { from, to } })

// downloadReport fetches the export as a blob (so the auth header is sent)
// and triggers a browser save-as, since a plain <a href> can't carry it.
export async function downloadReport(
  report: 'sales' | 'inventory' | 'product-performance',
  format: 'csv' | 'pdf',
  params: Record<string, string> = {},
): Promise<void> {
  const response = await client.get(`/reports/${report}`, {
    params: { ...params, format },
    responseType: 'blob',
  })
  const url = URL.createObjectURL(response.data as Blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `${report}-report.${format}`
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}
