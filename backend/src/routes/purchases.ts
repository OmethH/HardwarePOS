import { Router } from 'express'
import { db, transaction } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware, requireRole } from '../middleware/auth.js'
import { applyMovement } from '../inventory.js'

export const purchasesRouter = Router()
purchasesRouter.use(authMiddleware)

function getFullPurchase(id: number | bigint) {
  const p = db.prepare(`
    SELECT pu.*, 
           s.name as supplier_name, s.phone as supplier_phone, s.email as supplier_email,
           u.name as user_name, u.username as user_username
    FROM purchases pu
    LEFT JOIN suppliers s ON pu.supplier_id = s.id
    LEFT JOIN users u ON pu.created_by = u.id
    WHERE pu.id = ?
  `).get(id) as any

  if (!p) return null

  const items = db.prepare(`
    SELECT pi.*, 
           p.name as product_name, p.sku as product_sku, p.unit_type as product_unit_type
    FROM purchase_items pi
    LEFT JOIN products p ON pi.product_id = p.id
    WHERE pi.purchase_id = ?
  `).all(id) as any[]

  return {
    id: p.id,
    reference_number: p.reference_number,
    supplier_id: p.supplier_id,
    total_amount: p.total_amount,
    created_by: p.created_by,
    created_at: p.created_at,
    Supplier: p.supplier_name
      ? {
          id: p.supplier_id,
          name: p.supplier_name,
          phone: p.supplier_phone,
          email: p.supplier_email,
        }
      : null,
    CreatedByUser: p.user_name
      ? {
          id: p.created_by,
          name: p.user_name,
          username: p.user_username,
        }
      : null,
    Items: items.map((it) => ({
      id: it.id,
      purchase_id: it.purchase_id,
      product_id: it.product_id,
      quantity: it.quantity,
      price: it.price,
      Product: it.product_name
        ? {
            id: it.product_id,
            name: it.product_name,
            sku: it.product_sku,
            unit_type: it.product_unit_type,
          }
        : null,
    })),
  }
}

purchasesRouter.get('/purchases', requireRole('admin', 'manager'), (req, res) => {
  const limit = req.query.limit ? Number(req.query.limit) : 0
  let sql = `SELECT id FROM purchases ORDER BY created_at DESC`
  if (limit > 0) sql += ` LIMIT ${limit}`

  const rows = db.prepare(sql).all() as any[]
  const purchases = rows.map((r) => getFullPurchase(r.id)).filter(Boolean)
  return ok(res, purchases)
})

purchasesRouter.get('/purchases/:id', requireRole('admin', 'manager'), (req, res) => {
  const p = getFullPurchase(Number(req.params.id))
  if (!p) return err(res, 'purchase not found', 404)
  return ok(res, p)
})

purchasesRouter.post('/purchases', requireRole('admin', 'manager'), (req, res) => {
  const { supplier_id, items } = req.body
  if (!supplier_id || !Array.isArray(items) || items.length === 0) {
    return err(res, 'supplier_id and items are required', 400)
  }

  try {
    const purchaseId = transaction(() => {
      let total = 0
      for (const it of items) {
        total += Number(it.quantity) * Number(it.price)
      }

      const refNum = `PO-${Math.floor(Date.now() / 1000)}-${Math.floor(Math.random() * 1000)}`
      const stmt = db.prepare(`
        INSERT INTO purchases (reference_number, supplier_id, total_amount, created_by)
        VALUES (?, ?, ?, ?)
      `)
      const res = stmt.run(refNum, Number(supplier_id), total, req.user!.id)
      const pId = Number(res.lastInsertRowid)

      for (const it of items) {
        db.prepare(`
          INSERT INTO purchase_items (purchase_id, product_id, quantity, price)
          VALUES (?, ?, ?, ?)
        `).run(pId, Number(it.product_id), Number(it.quantity), Number(it.price))

        // Adjust stock
        applyMovement(
          Number(it.product_id),
          Number(it.quantity),
          'PURCHASE',
          'purchase',
          pId,
          `Purchase ${refNum}`,
          req.user!.id
        )

        // Update product purchase_price
        db.prepare('UPDATE products SET purchase_price = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?').run(
          Number(it.price),
          Number(it.product_id)
        )
      }

      return pId
    })

    const created = getFullPurchase(purchaseId)
    return ok(res, created, 201)
  } catch (e: any) {
    return err(res, e.message || 'failed to create purchase', 400)
  }
})
