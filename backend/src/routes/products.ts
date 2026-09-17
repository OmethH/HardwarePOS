import { Router } from 'express'
import { db } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware } from '../middleware/auth.js'

export const productsRouter = Router()
productsRouter.use(authMiddleware)

function formatProduct(row: any) {
  if (!row) return null
  return {
    ...row,
    Category: row.category_id
      ? {
          id: row.category_id,
          name: row.category_name,
        }
      : null,
  }
}

productsRouter.get('/products', (req, res) => {
  const search = typeof req.query.search === 'string' ? req.query.search.trim() : ''
  const categoryId = req.query.category_id ? Number(req.query.category_id) : 0
  const status = typeof req.query.status === 'string' ? req.query.status.trim() : ''
  const lowStock = req.query.low_stock === 'true'

  let sql = `
    SELECT p.*, c.name as category_name
    FROM products p
    LEFT JOIN categories c ON p.category_id = c.id
    WHERE 1=1
  `
  const params: any[] = []

  if (search) {
    sql += ` AND (p.name LIKE ? OR p.sku LIKE ? OR p.barcode LIKE ?)`
    const like = `%${search}%`
    params.push(like, like, like)
  }
  if (categoryId) {
    sql += ` AND p.category_id = ?`
    params.push(categoryId)
  }
  if (status) {
    sql += ` AND p.status = ?`
    params.push(status)
  }
  if (lowStock) {
    sql += ` AND p.current_stock <= p.minimum_stock`
  }

  sql += ` ORDER BY p.name ASC`

  const rows = db.prepare(sql).all(...params)
  return ok(res, rows.map(formatProduct))
})

productsRouter.get('/products/:id(\\d+)', (req, res) => {
  const row = db.prepare(`
    SELECT p.*, c.name as category_name
    FROM products p
    LEFT JOIN categories c ON p.category_id = c.id
    WHERE p.id = ?
  `).get(req.params.id)

  if (!row) return err(res, 'product not found', 404)
  return ok(res, formatProduct(row))
})

productsRouter.get('/products/barcode/:barcode', (req, res) => {
  const row = db.prepare(`
    SELECT p.*, c.name as category_name
    FROM products p
    LEFT JOIN categories c ON p.category_id = c.id
    WHERE p.barcode = ?
  `).get(req.params.barcode)

  if (!row) return err(res, 'product not found', 404)
  return ok(res, formatProduct(row))
})

productsRouter.post('/products', (req, res) => {
  const {
    sku,
    barcode = null,
    name,
    category_id = null,
    brand = '',
    description = '',
    unit_type = 'Piece',
    purchase_price = 0,
    selling_price = 0,
    minimum_stock = 0,
    status = 'active',
  } = req.body

  if (!sku || !name) {
    return err(res, 'sku and name are required', 400)
  }

  try {
    const stmt = db.prepare(`
      INSERT INTO products (
        sku, barcode, name, category_id, brand, description, unit_type,
        purchase_price, selling_price, minimum_stock, current_stock, status
      ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?)
    `)
    const result = stmt.run(
      sku.trim(),
      barcode ? String(barcode).trim() : null,
      name.trim(),
      category_id || null,
      brand,
      description,
      unit_type,
      purchase_price,
      selling_price,
      minimum_stock,
      status
    )

    const inserted = db.prepare(`
      SELECT p.*, c.name as category_name
      FROM products p
      LEFT JOIN categories c ON p.category_id = c.id
      WHERE p.id = ?
    `).get(result.lastInsertRowid)

    return ok(res, formatProduct(inserted), 201)
  } catch (e: any) {
    if (e.message?.includes('UNIQUE')) {
      return err(res, 'product with this sku or barcode already exists', 409)
    }
    return err(res, 'failed to create product', 500)
  }
})

productsRouter.put('/products/:id', (req, res) => {
  const {
    sku,
    barcode = null,
    name,
    category_id = null,
    brand = '',
    description = '',
    unit_type = 'Piece',
    purchase_price,
    selling_price = 0,
    minimum_stock = 0,
    status = 'active',
  } = req.body

  if (!sku || !name) {
    return err(res, 'sku and name are required', 400)
  }

  const existing = db.prepare('SELECT * FROM products WHERE id = ?').get(req.params.id) as any
  if (!existing) return err(res, 'product not found', 404)

  const isPrivileged = req.user?.role === 'admin' || req.user?.role === 'manager'
  if (!isPrivileged && purchase_price !== undefined && purchase_price !== existing.purchase_price) {
    return err(res, 'only admins and managers may change the purchase price', 403)
  }

  const finalPurchasePrice = isPrivileged && purchase_price !== undefined ? purchase_price : existing.purchase_price

  try {
    db.prepare(`
      UPDATE products SET
        sku = ?, barcode = ?, name = ?, category_id = ?, brand = ?, description = ?,
        unit_type = ?, purchase_price = ?, selling_price = ?, minimum_stock = ?,
        status = ?, updated_at = CURRENT_TIMESTAMP
      WHERE id = ?
    `).run(
      sku.trim(),
      barcode ? String(barcode).trim() : null,
      name.trim(),
      category_id || null,
      brand,
      description,
      unit_type,
      finalPurchasePrice,
      selling_price,
      minimum_stock,
      status,
      req.params.id
    )

    const updated = db.prepare(`
      SELECT p.*, c.name as category_name
      FROM products p
      LEFT JOIN categories c ON p.category_id = c.id
      WHERE p.id = ?
    `).get(req.params.id)

    return ok(res, formatProduct(updated))
  } catch (e: any) {
    if (e.message?.includes('UNIQUE')) {
      return err(res, 'product with this sku or barcode already exists', 409)
    }
    return err(res, 'failed to update product', 500)
  }
})

productsRouter.delete('/products/:id', (req, res) => {
  db.prepare('DELETE FROM products WHERE id = ?').run(req.params.id)
  return ok(res, { message: 'product deleted successfully' })
})
