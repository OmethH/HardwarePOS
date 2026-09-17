import { Router } from 'express'
import { db, transaction } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware } from '../middleware/auth.js'
import { applyMovement } from '../inventory.js'

export const salesRouter = Router()
salesRouter.use(authMiddleware)

export function getFullSale(idOrInvoice: number | string) {
  const isInvoice = typeof idOrInvoice === 'string' && isNaN(Number(idOrInvoice))
  const sql = `
    SELECT s.*, 
           c.name as customer_name, c.phone as customer_phone, c.address as customer_address,
           u.name as user_name, u.username as user_username
    FROM sales s
    LEFT JOIN customers c ON s.customer_id = c.id
    LEFT JOIN users u ON s.created_by = u.id
    WHERE ${isInvoice ? 's.invoice_number = ?' : 's.id = ?'}
  `
  const s = db.prepare(sql).get(idOrInvoice) as any
  if (!s) return null

  const items = db.prepare(`
    SELECT si.*, 
           p.name as product_name, p.sku as product_sku, p.unit_type as product_unit_type
    FROM sale_items si
    LEFT JOIN products p ON si.product_id = p.id
    WHERE si.sale_id = ?
  `).all(s.id) as any[]

  return {
    id: s.id,
    invoice_number: s.invoice_number,
    customer_id: s.customer_id,
    subtotal: s.subtotal,
    discount: s.discount,
    total: s.total,
    payment_method: s.payment_method,
    status: s.status,
    created_by: s.created_by,
    created_at: s.created_at,
    Customer: s.customer_name
      ? {
          id: s.customer_id,
          name: s.customer_name,
          phone: s.customer_phone,
          address: s.customer_address,
        }
      : null,
    CreatedByUser: s.user_name
      ? {
          id: s.created_by,
          name: s.user_name,
          username: s.user_username,
        }
      : null,
    Items: items.map((it) => ({
      id: it.id,
      sale_id: it.sale_id,
      product_id: it.product_id,
      quantity: it.quantity,
      returned_quantity: it.returned_quantity,
      selling_price: it.selling_price,
      purchase_price: it.purchase_price,
      profit: it.profit,
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

salesRouter.get('/sales', (req, res) => {
  const search = typeof req.query.search === 'string' ? req.query.search.trim() : ''
  const from = typeof req.query.from === 'string' ? req.query.from.trim() : ''
  const to = typeof req.query.to === 'string' ? req.query.to.trim() : ''
  const customerId = req.query.customer_id ? Number(req.query.customer_id) : 0

  let sql = `SELECT id FROM sales WHERE 1=1`
  const params: any[] = []

  if (search) {
    sql += ` AND invoice_number LIKE ?`
    params.push(`%${search}%`)
  }
  if (from) {
    sql += ` AND created_at >= ?`
    params.push(from)
  }
  if (to) {
    sql += ` AND created_at <= ?`
    params.push(to)
  }
  if (customerId) {
    sql += ` AND customer_id = ?`
    params.push(customerId)
  }

  // Cashiers only see their own sales
  if (req.user?.role === 'cashier') {
    sql += ` AND created_by = ?`
    params.push(req.user.id)
  }

  sql += ` ORDER BY created_at DESC`

  const rows = db.prepare(sql).all(...params) as any[]
  const list = rows.map((r) => getFullSale(r.id)).filter(Boolean)
  return ok(res, list)
})

salesRouter.get('/sales/:id(\\d+)', (req, res) => {
  const sale = getFullSale(Number(req.params.id))
  if (!sale) return err(res, 'sale not found', 404)
  return ok(res, sale)
})

salesRouter.get('/sales/invoice/:invoice', (req, res) => {
  const sale = getFullSale(req.params.invoice)
  if (!sale) return err(res, 'sale not found', 404)
  return ok(res, sale)
})

salesRouter.post('/sales', (req, res) => {
  const { customer_id = 1, discount = 0, payment_method = 'cash', items } = req.body
  if (!Array.isArray(items) || items.length === 0) {
    return err(res, 'a sale must include at least one item', 400)
  }

  try {
    const saleId = transaction(() => {
      // Validate all items & stock
      let subtotal = 0
      const validatedItems: Array<{ product: any; quantity: number }> = []

      for (const it of items) {
        const product = db.prepare('SELECT * FROM products WHERE id = ?').get(it.product_id) as any
        if (!product) {
          throw new Error(`product ${it.product_id} not found`)
        }
        const qty = Number(it.quantity)
        if (qty <= 0) {
          throw new Error(`invalid quantity for ${product.name}`)
        }
        if (product.current_stock < qty) {
          throw new Error(
            `insufficient stock for ${product.name}: available ${product.current_stock}, requested ${qty}`
          )
        }
        subtotal += Number(product.selling_price) * qty
        validatedItems.push({ product, quantity: qty })
      }

      const disc = Math.min(Number(discount) || 0, subtotal)
      const total = subtotal - disc

      // Insert sale
      const tempInvoice = `INV-PENDING-${Date.now()}`
      const stmt = db.prepare(`
        INSERT INTO sales (
          invoice_number, customer_id, subtotal, discount, total, payment_method, status, created_by
        ) VALUES (?, ?, ?, ?, ?, ?, 'COMPLETED', ?)
      `)
      const res = stmt.run(tempInvoice, Number(customer_id) || 1, subtotal, disc, total, payment_method, req.user!.id)
      const sId = Number(res.lastInsertRowid)

      // Set clean sequential invoice number
      const invoiceNumber = `INV${String(sId).padStart(5, '0')}`
      db.prepare('UPDATE sales SET invoice_number = ? WHERE id = ?').run(invoiceNumber, sId)

      // Insert sale items & deduct stock
      for (const { product, quantity } of validatedItems) {
        const profit = (Number(product.selling_price) - Number(product.purchase_price)) * quantity
        db.prepare(`
          INSERT INTO sale_items (
            sale_id, product_id, quantity, returned_quantity, selling_price, purchase_price, profit
          ) VALUES (?, ?, ?, 0, ?, ?, ?)
        `).run(sId, product.id, quantity, product.selling_price, product.purchase_price, profit)

        applyMovement(
          product.id,
          -quantity,
          'SALE',
          'sale',
          sId,
          `Sale ${invoiceNumber}`,
          req.user!.id
        )
      }

      return sId
    })

    const created = getFullSale(saleId)
    return ok(res, created, 201)
  } catch (e: any) {
    return err(res, e.message || 'checkout failed', 400)
  }
})
