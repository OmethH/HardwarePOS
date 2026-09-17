import { Router } from 'express'
import { db } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware, requireRole } from '../middleware/auth.js'

export const settingsRouter = Router()
settingsRouter.use(authMiddleware)

settingsRouter.get('/settings', (_req, res) => {
  const settings = db.prepare('SELECT * FROM settings WHERE id = 1').get()
  if (!settings) {
    return ok(res, {
      id: 1,
      shop_name: 'HardwarePOS',
      address: '',
      phone: '',
      receipt_footer: 'Thank You',
      tax_percentage: 0,
    })
  }
  return ok(res, settings)
})

settingsRouter.put('/settings', requireRole('admin', 'manager'), (req, res) => {
  const {
    shop_name = 'HardwarePOS',
    address = '',
    phone = '',
    receipt_footer = 'Thank You',
    tax_percentage = 0,
  } = req.body

  db.prepare(`
    INSERT INTO settings (id, shop_name, address, phone, receipt_footer, tax_percentage, updated_at)
    VALUES (1, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
    ON CONFLICT(id) DO UPDATE SET
      shop_name = excluded.shop_name,
      address = excluded.address,
      phone = excluded.phone,
      receipt_footer = excluded.receipt_footer,
      tax_percentage = excluded.tax_percentage,
      updated_at = CURRENT_TIMESTAMP
  `).run(shop_name, address, phone, receipt_footer, Number(tax_percentage))

  const updated = db.prepare('SELECT * FROM settings WHERE id = 1').get()
  return ok(res, updated)
})
