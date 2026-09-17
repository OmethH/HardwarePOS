import { Router } from 'express'
import { db } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware, requireRole } from '../middleware/auth.js'

export const dashboardRouter = Router()
dashboardRouter.use(authMiddleware)
dashboardRouter.use(requireRole('admin', 'manager'))

dashboardRouter.get('/dashboard/summary', (_req, res) => {
  try {
    const todaySales = db.prepare(`
      SELECT 
        COALESCE(SUM(total), 0) as sales_amount,
        COUNT(*) as transaction_count
      FROM sales
      WHERE date(created_at) = date('now')
    `).get() as any

    const todayItems = db.prepare(`
      SELECT 
        COALESCE(SUM(si.profit), 0) as profit,
        COALESCE(SUM(si.quantity), 0) as items_sold
      FROM sale_items si
      JOIN sales s ON s.id = si.sale_id
      WHERE date(s.created_at) = date('now')
    `).get() as any

    const totalProducts = db.prepare('SELECT COUNT(*) as count FROM products').get() as any
    const lowStock = db.prepare(`
      SELECT COUNT(*) as count FROM products 
      WHERE current_stock <= minimum_stock AND current_stock > 0
    `).get() as any
    const outOfStock = db.prepare('SELECT COUNT(*) as count FROM products WHERE current_stock <= 0').get() as any

    const summary = {
      today: {
        sales_amount: Number(todaySales?.sales_amount || 0),
        profit: Number(todayItems?.profit || 0),
        transaction_count: Number(todaySales?.transaction_count || 0),
        items_sold: Number(todayItems?.items_sold || 0),
      },
      inventory: {
        total_products: Number(totalProducts?.count || 0),
        low_stock_products: Number(lowStock?.count || 0),
        out_of_stock_products: Number(outOfStock?.count || 0),
      },
    }

    return ok(res, summary)
  } catch (e: any) {
    return err(res, 'failed to load dashboard summary', 500)
  }
})

dashboardRouter.get('/dashboard/sales-over-time', (req, res) => {
  const days = Number(req.query.days) || 30
  try {
    const rows = db.prepare(`
      SELECT date(created_at) as date, COALESCE(SUM(total), 0) as total
      FROM sales
      WHERE created_at >= date('now', '-' || ? || ' days')
      GROUP BY date(created_at)
      ORDER BY date ASC
    `).all(days) as any[]

    return ok(res, rows)
  } catch (e: any) {
    return err(res, 'failed to load sales over time', 500)
  }
})

dashboardRouter.get('/dashboard/profit-trend', (req, res) => {
  const days = Number(req.query.days) || 30
  try {
    const rows = db.prepare(`
      SELECT date(s.created_at) as date, COALESCE(SUM(si.profit), 0) as profit
      FROM sale_items si
      JOIN sales s ON s.id = si.sale_id
      WHERE s.created_at >= date('now', '-' || ? || ' days')
      GROUP BY date(s.created_at)
      ORDER BY date ASC
    `).all(days) as any[]

    return ok(res, rows)
  } catch (e: any) {
    return err(res, 'failed to load profit trend', 500)
  }
})

dashboardRouter.get('/dashboard/top-products', (req, res) => {
  const days = Number(req.query.days) || 30
  const limit = Number(req.query.limit) || 10
  try {
    const rows = db.prepare(`
      SELECT 
        p.id as product_id,
        p.name as product_name,
        COALESCE(SUM(si.quantity), 0) as quantity_sold,
        COALESCE(SUM(si.quantity * si.selling_price), 0) as revenue
      FROM sale_items si
      JOIN sales s ON s.id = si.sale_id
      JOIN products p ON p.id = si.product_id
      WHERE s.created_at >= date('now', '-' || ? || ' days')
      GROUP BY p.id, p.name
      ORDER BY quantity_sold DESC
      LIMIT ?
    `).all(days, limit) as any[]

    return ok(res, rows)
  } catch (e: any) {
    return err(res, 'failed to load top products', 500)
  }
})
