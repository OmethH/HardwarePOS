import { Router } from 'express'
import { db, transaction } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware, requireRole } from '../middleware/auth.js'
import { applyMovement } from '../inventory.js'

export const inventoryRouter = Router()
inventoryRouter.use(authMiddleware)

function formatMovement(row: any) {
  if (!row) return null
  return {
    id: row.id,
    product_id: row.product_id,
    quantity: row.quantity,
    type: row.type,
    reference_type: row.reference_type,
    reference_id: row.reference_id,
    note: row.note,
    created_by: row.created_by,
    created_at: row.created_at,
    Product: row.product_name
      ? {
          id: row.product_id,
          sku: row.product_sku,
          name: row.product_name,
          unit_type: row.product_unit_type,
        }
      : null,
    CreatedByUser: row.user_name
      ? {
          id: row.created_by,
          name: row.user_name,
          username: row.user_username,
        }
      : null,
  }
}

inventoryRouter.get('/inventory/movements', (req, res) => {
  const limit = req.query.limit ? Number(req.query.limit) : 0
  let sql = `
    SELECT sm.*, 
           p.name as product_name, p.sku as product_sku, p.unit_type as product_unit_type,
           u.name as user_name, u.username as user_username
    FROM stock_movements sm
    LEFT JOIN products p ON sm.product_id = p.id
    LEFT JOIN users u ON sm.created_by = u.id
    ORDER BY sm.created_at DESC
  `
  if (limit > 0) {
    sql += ` LIMIT ${limit}`
  }

  const rows = db.prepare(sql).all()
  return ok(res, rows.map(formatMovement))
})

inventoryRouter.get('/inventory/movements/product/:id', (req, res) => {
  const sql = `
    SELECT sm.*, 
           p.name as product_name, p.sku as product_sku, p.unit_type as product_unit_type,
           u.name as user_name, u.username as user_username
    FROM stock_movements sm
    LEFT JOIN products p ON sm.product_id = p.id
    LEFT JOIN users u ON sm.created_by = u.id
    WHERE sm.product_id = ?
    ORDER BY sm.created_at DESC
  `
  const rows = db.prepare(sql).all(req.params.id)
  return ok(res, rows.map(formatMovement))
})

inventoryRouter.post('/inventory/adjustments', requireRole('admin', 'manager'), (req, res) => {
  const { product_id, quantity, note = '' } = req.body
  if (!product_id || quantity === undefined || quantity === 0) {
    return err(res, 'valid product_id and non-zero quantity are required', 400)
  }

  try {
    const movementId = transaction(() => {
      return applyMovement(
        Number(product_id),
        Number(quantity),
        'ADJUSTMENT',
        'manual',
        null,
        note,
        req.user!.id
      )
    })

    const row = db.prepare(`
      SELECT sm.*, 
             p.name as product_name, p.sku as product_sku, p.unit_type as product_unit_type,
             u.name as user_name, u.username as user_username
      FROM stock_movements sm
      LEFT JOIN products p ON sm.product_id = p.id
      LEFT JOIN users u ON sm.created_by = u.id
      WHERE sm.id = ?
    `).get(movementId)

    return ok(res, formatMovement(row), 201)
  } catch (e: any) {
    return err(res, e.message || 'failed to apply adjustment', 400)
  }
})
