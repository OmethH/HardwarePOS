import { Router } from 'express'
import { db, transaction } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware, requireRole } from '../middleware/auth.js'
import { applyMovement } from '../inventory.js'
import { getFullSale } from './sales.js'

export const returnsRouter = Router()
returnsRouter.use(authMiddleware)

function getFullReturn(id: number | bigint) {
  const r = db.prepare(`
    SELECT ret.*, 
           u.name as user_name, u.username as user_username
    FROM returns ret
    LEFT JOIN users u ON ret.created_by = u.id
    WHERE ret.id = ?
  `).get(id) as any

  if (!r) return null

  const sale = getFullSale(r.sale_id)
  const items = db.prepare(`
    SELECT ri.*, 
           p.name as product_name, p.sku as product_sku, p.unit_type as product_unit_type
    FROM return_items ri
    LEFT JOIN products p ON ri.product_id = p.id
    WHERE ri.return_id = ?
  `).all(id) as any[]

  return {
    id: r.id,
    return_number: r.return_number,
    sale_id: r.sale_id,
    total_refund: r.total_refund,
    created_by: r.created_by,
    created_at: r.created_at,
    Sale: sale,
    CreatedByUser: r.user_name
      ? {
          id: r.created_by,
          name: r.user_name,
          username: r.user_username,
        }
      : null,
    Items: items.map((it) => ({
      id: it.id,
      return_id: it.return_id,
      sale_item_id: it.sale_item_id,
      product_id: it.product_id,
      quantity: it.quantity,
      refund_amount: it.refund_amount,
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

returnsRouter.get('/returns', requireRole('admin', 'manager'), (req, res) => {
  const limit = req.query.limit ? Number(req.query.limit) : 0
  let sql = 'SELECT id FROM returns ORDER BY created_at DESC'
  if (limit > 0) sql += ` LIMIT ${limit}`

  const rows = db.prepare(sql).all() as any[]
  const list = rows.map((r) => getFullReturn(r.id)).filter(Boolean)
  return ok(res, list)
})

returnsRouter.get('/returns/:id', requireRole('admin', 'manager'), (req, res) => {
  const ret = getFullReturn(Number(req.params.id))
  if (!ret) return err(res, 'return not found', 404)
  return ok(res, ret)
})

returnsRouter.post('/returns', requireRole('admin', 'manager'), (req, res) => {
  const { sale_id, items } = req.body
  if (!sale_id || !Array.isArray(items) || items.length === 0) {
    return err(res, 'sale_id and items are required', 400)
  }

  try {
    const returnId = transaction(() => {
      const sale = db.prepare('SELECT * FROM sales WHERE id = ?').get(sale_id) as any
      if (!sale) {
        throw new Error('sale not found')
      }

      const saleItems = db.prepare('SELECT * FROM sale_items WHERE sale_id = ?').all(sale_id) as any[]
      const saleItemById = new Map<number, any>(saleItems.map((si) => [si.id, si]))

      let totalRefund = 0
      const returnNumber = `RET-${Math.floor(Date.now() / 1000)}-${Math.floor(Math.random() * 1000)}`

      const stmt = db.prepare(`
        INSERT INTO returns (return_number, sale_id, total_refund, created_by)
        VALUES (?, ?, 0, ?)
      `)
      const retRes = stmt.run(returnNumber, Number(sale_id), req.user!.id)
      const retId = Number(retRes.lastInsertRowid)

      for (const it of items) {
        const saleItem = saleItemById.get(Number(it.sale_item_id))
        if (!saleItem || saleItem.sale_id !== sale.id) {
          throw new Error(`sale item ${it.sale_item_id} does not belong to sale ${sale_id}`)
        }

        const qty = Number(it.quantity)
        const remaining = Number(saleItem.quantity) - Number(saleItem.returned_quantity)
        if (qty <= 0 || qty > remaining) {
          throw new Error(
            `return quantity exceeds available quantity sold (remaining ${remaining}, requested ${qty})`
          )
        }

        const refund = qty * Number(saleItem.selling_price)
        totalRefund += refund

        db.prepare(`
          INSERT INTO return_items (return_id, sale_item_id, product_id, quantity, refund_amount)
          VALUES (?, ?, ?, ?, ?)
        `).run(retId, saleItem.id, saleItem.product_id, qty, refund)

        // Update returned_quantity on sale item
        const newReturnedQty = Number(saleItem.returned_quantity) + qty
        db.prepare('UPDATE sale_items SET returned_quantity = ? WHERE id = ?').run(newReturnedQty, saleItem.id)
        saleItem.returned_quantity = newReturnedQty

        // Restock
        applyMovement(
          saleItem.product_id,
          qty,
          'RETURN',
          'return',
          retId,
          `Return ${returnNumber}`,
          req.user!.id
        )
      }

      // Update total refund on return record
      db.prepare('UPDATE returns SET total_refund = ? WHERE id = ?').run(totalRefund, retId)

      // Update sale status
      let fullyReturned = true
      for (const si of saleItems) {
        if (Number(si.returned_quantity) < Number(si.quantity)) {
          fullyReturned = false
          break
        }
      }
      const newStatus = fullyReturned ? 'RETURNED' : 'PARTIALLY_RETURNED'
      db.prepare('UPDATE sales SET status = ? WHERE id = ?').run(newStatus, sale.id)

      return retId
    })

    const created = getFullReturn(returnId)
    return ok(res, created, 201)
  } catch (e: any) {
    return err(res, e.message || 'failed to process return', 400)
  }
})
