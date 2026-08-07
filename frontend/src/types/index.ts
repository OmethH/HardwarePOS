export interface Role {
  id: number
  name: string
}

export interface User {
  id: number
  name: string
  username: string
  role_id: number
  role: Role
  active: boolean
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  name: string
  created_at: string
  updated_at: string
}

export interface Supplier {
  id: number
  name: string
  phone: string
  email: string
  address: string
  company: string
  created_at: string
  updated_at: string
}

export interface Product {
  id: number
  sku: string
  barcode: string | null
  name: string
  category_id: number | null
  category?: Category
  brand: string
  description: string
  unit_type: string
  purchase_price: number
  selling_price: number
  minimum_stock: number
  current_stock: number
  status: string
  created_at: string
  updated_at: string
}

export interface StockMovement {
  id: number
  product_id: number
  product?: Product
  quantity: number
  type: 'PURCHASE' | 'SALE' | 'RETURN' | 'ADJUSTMENT'
  reference_type: string
  reference_id: number | null
  note: string
  created_by: number | null
  created_by_user?: User
  created_at: string
}

export interface PurchaseItem {
  id: number
  purchase_id: number
  product_id: number
  product?: Product
  quantity: number
  price: number
}

export interface Purchase {
  id: number
  reference_number: string
  supplier_id: number
  supplier?: Supplier
  total_amount: number
  created_by: number | null
  created_by_user?: User
  created_at: string
  items: PurchaseItem[]
}

export interface Customer {
  id: number
  name: string
  phone: string
  address: string
  created_at: string
  updated_at: string
}

export type SaleStatus = 'COMPLETED' | 'PARTIALLY_RETURNED' | 'RETURNED'

export interface SaleItem {
  id: number
  sale_id: number
  product_id: number
  product?: Product
  quantity: number
  returned_quantity: number
  selling_price: number
  purchase_price: number
  profit: number
}

export interface Sale {
  id: number
  invoice_number: string
  customer_id: number
  customer?: Customer
  subtotal: number
  discount: number
  total: number
  payment_method: string
  status: SaleStatus
  created_by: number | null
  created_by_user?: User
  created_at: string
  items: SaleItem[]
}

export interface ReturnItem {
  id: number
  return_id: number
  sale_item_id: number
  product_id: number
  product?: Product
  quantity: number
  refund_amount: number
}

export interface ReturnRecord {
  id: number
  return_number: string
  sale_id: number
  sale?: Sale
  total_refund: number
  created_by: number | null
  created_by_user?: User
  created_at: string
  items: ReturnItem[]
}

export interface Settings {
  id: number
  shop_name: string
  address: string
  phone: string
  receipt_footer: string
  tax_percentage: number
  updated_at: string
}

export interface DashboardSummary {
  today: {
    sales_amount: number
    profit: number
    transaction_count: number
    items_sold: number
  }
  inventory: {
    total_products: number
    low_stock_products: number
    out_of_stock_products: number
  }
}

export interface SalesReportRow {
  date: string
  invoices: number
  revenue: number
  discount: number
  profit: number
}

export interface InventoryReportRow {
  product: string
  sku: string
  current_stock: number
  stock_value: number
  status: string
}

export interface ProductPerformanceRow {
  product: string
  quantity_sold: number
  revenue: number
  profit: number
}

export interface ApiError {
  error: string
}
