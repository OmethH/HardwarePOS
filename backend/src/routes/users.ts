import { Router } from 'express'
import bcrypt from 'bcryptjs'
import { db } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware, requireRole } from '../middleware/auth.js'

export const usersRouter = Router()
usersRouter.use(authMiddleware)

function formatUser(row: any) {
  if (!row) return null
  const { password_hash, ...rest } = row
  return {
    ...rest,
    Role: row.role_name
      ? {
          id: row.role_id,
          name: row.role_name,
        }
      : null,
  }
}

usersRouter.get('/roles', requireRole('admin'), (_req, res) => {
  const roles = db.prepare('SELECT * FROM roles ORDER BY id ASC').all()
  return ok(res, roles)
})

usersRouter.get('/users', requireRole('admin'), (_req, res) => {
  const rows = db.prepare(`
    SELECT u.*, r.name as role_name
    FROM users u
    LEFT JOIN roles r ON u.role_id = r.id
    ORDER BY u.name ASC
  `).all()
  return ok(res, rows.map(formatUser))
})

usersRouter.get('/users/:id', requireRole('admin'), (req, res) => {
  const row = db.prepare(`
    SELECT u.*, r.name as role_name
    FROM users u
    LEFT JOIN roles r ON u.role_id = r.id
    WHERE u.id = ?
  `).get(req.params.id)

  if (!row) return err(res, 'user not found', 404)
  return ok(res, formatUser(row))
})

usersRouter.post('/users', requireRole('admin'), (req, res) => {
  const { name, username, password, role_id } = req.body
  if (!name || !username || !password || !role_id || password.length < 6) {
    return err(res, 'name, username, password (min 6 chars) and role_id are required', 400)
  }

  const hash = bcrypt.hashSync(password, 10)
  try {
    const stmt = db.prepare(`
      INSERT INTO users (name, username, password_hash, role_id, active)
      VALUES (?, ?, ?, ?, 1)
    `)
    const result = stmt.run(name.trim(), username.trim(), hash, Number(role_id))
    const inserted = db.prepare(`
      SELECT u.*, r.name as role_name
      FROM users u
      LEFT JOIN roles r ON u.role_id = r.id
      WHERE u.id = ?
    `).get(result.lastInsertRowid)

    return ok(res, formatUser(inserted), 201)
  } catch (e: any) {
    if (e.message?.includes('UNIQUE')) {
      return err(res, 'username already exists', 409)
    }
    return err(res, 'could not create user', 500)
  }
})

usersRouter.put('/users/:id', requireRole('admin'), (req, res) => {
  const { name, role_id, active, password } = req.body
  if (!name || !role_id) {
    return err(res, 'name and role_id are required', 400)
  }

  const existing = db.prepare('SELECT * FROM users WHERE id = ?').get(req.params.id) as any
  if (!existing) return err(res, 'user not found', 404)

  const finalHash = password && password.length >= 6 ? bcrypt.hashSync(password, 10) : existing.password_hash
  const isActive = active !== undefined ? (active ? 1 : 0) : existing.active

  db.prepare(`
    UPDATE users SET
      name = ?, role_id = ?, active = ?, password_hash = ?, updated_at = CURRENT_TIMESTAMP
    WHERE id = ?
  `).run(name.trim(), Number(role_id), isActive, finalHash, req.params.id)

  const updated = db.prepare(`
    SELECT u.*, r.name as role_name
    FROM users u
    LEFT JOIN roles r ON u.role_id = r.id
    WHERE u.id = ?
  `).get(req.params.id)

  return ok(res, formatUser(updated))
})

usersRouter.delete('/users/:id', requireRole('admin'), (req, res) => {
  // Soft deactivate so history remains intact
  db.prepare('UPDATE users SET active = 0 WHERE id = ?').run(req.params.id)
  return ok(res, { deactivated: true })
})
