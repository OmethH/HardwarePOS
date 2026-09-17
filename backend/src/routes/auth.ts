import { Router } from 'express'
import bcrypt from 'bcryptjs'
import jwt from 'jsonwebtoken'
import { db } from '../database.js'
import { config } from '../config.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware } from '../middleware/auth.js'

export const authRouter = Router()

authRouter.post('/auth/login', (req, res) => {
  const { username, password } = req.body
  if (!username || !password) {
    return err(res, 'username and password are required', 400)
  }

  const query = db.prepare(`
    SELECT u.id, u.name, u.username, u.password_hash, u.role_id, u.active, r.name as role
    FROM users u
    JOIN roles r ON u.role_id = r.id
    WHERE u.username = ? AND u.active = 1
  `)
  const user = query.get(username) as any

  if (!user || !bcrypt.compareSync(password, user.password_hash)) {
    return err(res, 'invalid username or password', 401)
  }

  const token = jwt.sign(
    {
      user_id: user.id,
      username: user.username,
      role: user.role,
      role_id: user.role_id,
    },
    config.jwtSecret,
    { expiresIn: '24h' }
  )

  const { password_hash, ...userProfile } = user
  return ok(res, {
    token,
    user: userProfile,
  })
})

authRouter.get('/auth/me', authMiddleware, (req, res) => {
  const query = db.prepare(`
    SELECT u.id, u.name, u.username, u.role_id, u.active, r.name as role
    FROM users u
    JOIN roles r ON u.role_id = r.id
    WHERE u.id = ?
  `)
  const user = query.get(req.user!.id)
  if (!user) {
    return err(res, 'user not found', 404)
  }
  return ok(res, user)
})
