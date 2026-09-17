import { Router } from 'express'
import { db } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware } from '../middleware/auth.js'

export const customersRouter = Router()
customersRouter.use(authMiddleware)

customersRouter.get('/customers', (req, res) => {
  const search = typeof req.query.search === 'string' ? req.query.search.trim() : ''
  if (search) {
    const like = `%${search}%`
    const rows = db.prepare(`
      SELECT * FROM customers 
      WHERE name LIKE ? OR phone LIKE ? 
      ORDER BY name ASC
    `).all(like, like)
    return ok(res, rows)
  }
  const rows = db.prepare('SELECT * FROM customers ORDER BY name ASC').all()
  return ok(res, rows)
})

customersRouter.get('/customers/:id', (req, res) => {
  const row = db.prepare('SELECT * FROM customers WHERE id = ?').get(req.params.id)
  if (!row) return err(res, 'customer not found', 404)
  return ok(res, row)
})

customersRouter.post('/customers', (req, res) => {
  const { name, phone = '', address = '' } = req.body
  if (!name || !name.trim()) return err(res, 'name is required', 400)

  const stmt = db.prepare('INSERT INTO customers (name, phone, address) VALUES (?, ?, ?)')
  const result = stmt.run(name.trim(), phone, address)
  const inserted = db.prepare('SELECT * FROM customers WHERE id = ?').get(result.lastInsertRowid)
  return ok(res, inserted, 201)
})

customersRouter.put('/customers/:id', (req, res) => {
  const { name, phone = '', address = '' } = req.body
  if (!name || !name.trim()) return err(res, 'name is required', 400)

  db.prepare(`
    UPDATE customers
    SET name = ?, phone = ?, address = ?, updated_at = CURRENT_TIMESTAMP
    WHERE id = ?
  `).run(name.trim(), phone, address, req.params.id)

  const updated = db.prepare('SELECT * FROM customers WHERE id = ?').get(req.params.id)
  if (!updated) return err(res, 'customer not found', 404)
  return ok(res, updated)
})

customersRouter.delete('/customers/:id', (req, res) => {
  if (req.params.id === '1') {
    return err(res, 'cannot delete default walk-in customer', 400)
  }
  db.prepare('DELETE FROM customers WHERE id = ?').run(req.params.id)
  return ok(res, { message: 'customer deleted successfully' })
})
