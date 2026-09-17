import { Router } from 'express'
import { db } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware } from '../middleware/auth.js'

export const categoriesRouter = Router()
categoriesRouter.use(authMiddleware)

categoriesRouter.get('/categories', (_req, res) => {
  const rows = db.prepare('SELECT * FROM categories ORDER BY name ASC').all()
  return ok(res, rows)
})

categoriesRouter.get('/categories/:id', (req, res) => {
  const cat = db.prepare('SELECT * FROM categories WHERE id = ?').get(req.params.id)
  if (!cat) return err(res, 'category not found', 404)
  return ok(res, cat)
})

categoriesRouter.post('/categories', (req, res) => {
  const { name } = req.body
  if (!name || !name.trim()) {
    return err(res, 'name is required', 400)
  }

  try {
    const stmt = db.prepare('INSERT INTO categories (name) VALUES (?)')
    const result = stmt.run(name.trim())
    const inserted = db.prepare('SELECT * FROM categories WHERE id = ?').get(result.lastInsertRowid)
    return ok(res, inserted, 201)
  } catch (e: any) {
    if (e.message?.includes('UNIQUE')) {
      return err(res, 'category name already exists', 409)
    }
    return err(res, 'failed to create category', 500)
  }
})

categoriesRouter.put('/categories/:id', (req, res) => {
  const { name } = req.body
  if (!name || !name.trim()) {
    return err(res, 'name is required', 400)
  }

  try {
    db.prepare('UPDATE categories SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?').run(name.trim(), req.params.id)
    const updated = db.prepare('SELECT * FROM categories WHERE id = ?').get(req.params.id)
    if (!updated) return err(res, 'category not found', 404)
    return ok(res, updated)
  } catch (e: any) {
    if (e.message?.includes('UNIQUE')) {
      return err(res, 'category name already exists', 409)
    }
    return err(res, 'failed to update category', 500)
  }
})

categoriesRouter.delete('/categories/:id', (req, res) => {
  // Check if any product is using this category
  const count = db.prepare('SELECT COUNT(*) as count FROM products WHERE category_id = ?').get(req.params.id) as any
  if (count && count.count > 0) {
    return err(res, 'cannot delete category with existing products', 400)
  }

  db.prepare('DELETE FROM categories WHERE id = ?').run(req.params.id)
  return ok(res, { message: 'category deleted successfully' })
})
