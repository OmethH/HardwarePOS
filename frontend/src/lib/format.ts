export function formatMoney(value: number): string {
  return value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

export function formatQty(value: number): string {
  return Number.isInteger(value) ? String(value) : value.toFixed(2)
}

export function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString()
}

export function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString()
}
