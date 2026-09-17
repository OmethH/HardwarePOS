import { Router } from 'express'
import { db } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware } from '../middleware/auth.js'

export const suppliersRouter = Router()
suppliersRouter.use(authMiddleware)

suppliersRouter.get('/suppliers', (req, res) => {
  const search = typeof req.query.search === 'string' ? req.query.search.trim() : ''
  if (search) {
    const like = `%${search}%`
    const rows = db.prepare(`
      SELECT * FROM suppliers 
      WHERE name LIKE ? OR company LIKE ? OR phone LIKE ? 
      ORDER BY name ASC
    `).all(like, like, like)
    return ok(res, rows)
  }
  const rows = db.prepare('SELECT * FROM suppliers ORDER BY name ASC').all()
  return ok(res, rows)
})

suppliersRouter.get('/suppliers/:id', (req, res) => {
  const row = db.prepare('SELECT * FROM suppliers WHERE id = ?').get(req.params.id)
  if (!row) return err(res, 'supplier not found', 404)
  return ok(res, row)
})

suppliersRouter.post('/suppliers', (req, res) => {
  const { name, phone = '', email = '', address = '', company = '' } = req.body
  if (!name || !name.trim()) return err(res, 'name is required', 400)

  const stmt = db.prepare(`
    INSERT INTO suppliers (name, phone, email, address, company)
    VALUES (?, ?, ?, ?, ?)
  `)
  const result = stmt.run(name.trim(), phone, email, address, company)
  const inserted = db.prepare('SELECT * FROM suppliers WHERE id = ?').get(result.lastInsertRowid)
  return ok(res, inserted, 201)
})

suppliersRouter.put('/suppliers/:id', (req, res) => {
  const { name, phone = '', email = '', address = '', company = '' } = req.body
  if (!name || !name.trim()) return err(res, 'name is required', 400)

  db.prepare(`
    UPDATE suppliers
    SET name = ?, phone = ?, email = ?, address = ?, company = ?, updated_at = CURRENT_TIMESTAMP
    WHERE id = ?
  `).run(name.trim(), phone, email, address, company, req.params.id)

  const updated = db.prepare('SELECT * FROM suppliers WHERE id = ?').get(req.params.id)
  if (!updated) return err(res, 'supplier not found', 404)
  return ok(res, updated)
})

suppliersRouter.delete('/suppliers/:id', (req, res) => {
  db.prepare('DELETE FROM suppliers WHERE id = ?').run(req.params.id)
  return ok(res, { message: 'supplier deleted successfully' })
})
