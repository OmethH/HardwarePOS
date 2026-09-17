import { db } from './database.js'

export function applyMovement(
  productId: number,
  quantity: number,
  movementType: 'PURCHASE' | 'SALE' | 'RETURN' | 'ADJUSTMENT',
  referenceType: string | null = null,
  referenceId: number | null = null,
  note: string = '',
  userId: number | null = null
) {
  const product = db.prepare('SELECT * FROM products WHERE id = ?').get(productId) as any
  if (!product) {
    throw new Error(`product ${productId} not found`)
  }

  const newStock = Number(product.current_stock) + Number(quantity)
  if (newStock < 0) {
    throw new Error(
      `insufficient stock for ${product.name}: available ${product.current_stock}, requested ${-quantity}`
    )
  }

  db.prepare('UPDATE products SET current_stock = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?').run(
    newStock,
    productId
  )

  const stmt = db.prepare(`
    INSERT INTO stock_movements (
      product_id, quantity, type, reference_type, reference_id, note, created_by
    ) VALUES (?, ?, ?, ?, ?, ?, ?)
  `)
  const res = stmt.run(productId, quantity, movementType, referenceType, referenceId, note, userId)
  return res.lastInsertRowid
}
