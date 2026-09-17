import type { Request, Response, NextFunction } from 'express'
import jwt from 'jsonwebtoken'
import { config } from '../config.js'
import { err } from '../httpresp.js'

export interface AuthUser {
  id: number
  username: string
  role: string
  role_id: number
}

declare global {
  namespace Express {
    interface Request {
      user?: AuthUser
    }
  }
}

export function authMiddleware(req: Request, res: Response, next: NextFunction) {
  const header = req.headers.authorization
  if (!header || !header.startsWith('Bearer ')) {
    return err(res, 'unauthorized', 401)
  }

  const token = header.slice(7).trim()
  try {
    const payload = jwt.verify(token, config.jwtSecret) as any
    req.user = {
      id: payload.user_id || payload.id,
      username: payload.username,
      role: payload.role,
      role_id: payload.role_id,
    }
    next()
  } catch {
    return err(res, 'invalid or expired token', 401)
  }
}

export function requireRole(...allowedRoles: string[]) {
  return (req: Request, res: Response, next: NextFunction) => {
    if (!req.user || !allowedRoles.includes(req.user.role)) {
      return err(res, 'forbidden: insufficient permissions', 403)
    }
    next()
  }
}
