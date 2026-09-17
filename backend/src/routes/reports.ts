import { Router, type Response } from 'express'
import PDFDocument from 'pdfkit'
import { db } from '../database.js'
import { ok, err } from '../httpresp.js'
import { authMiddleware, requireRole } from '../middleware/auth.js'

export const reportsRouter = Router()
reportsRouter.use(authMiddleware)
reportsRouter.use(requireRole('admin', 'manager'))

function exportCSV(res: Response, filename: string, headers: string[], rows: (string | number)[][]) {
  res.setHeader('Content-Disposition', `attachment; filename="${filename}"`)
  res.setHeader('Content-Type', 'text/csv')

  const csvContent = [
    headers.join(','),
    ...rows.map((row) =>
      row
        .map((cell) => {
          const str = String(cell ?? '')
          return str.includes(',') || str.includes('"') || str.includes('\n')
            ? `"${str.replace(/"/g, '""')}"`
            : str
        })
        .join(',')
    ),
  ].join('\n')

  return res.send(csvContent)
}

function exportPDF(
  res: Response,
  filename: string,
  title: string,
  headers: string[],
  rows: (string | number)[][]
) {
  res.setHeader('Content-Disposition', `attachment; filename="${filename}"`)
  res.setHeader('Content-Type', 'application/pdf')

  const doc = new PDFDocument({ margin: 40, size: 'A4' })
  doc.pipe(res)

  doc.fontSize(18).text(title, { align: 'center' })
  doc.fontSize(10).text(`Generated on ${new Date().toLocaleString()}`, { align: 'center' })
  doc.moveDown(1.5)

  // Table setup
  const startX = 40
  let currentY = doc.y
  const colWidth = (doc.page.width - 80) / headers.length

  // Header row
  doc.font('Helvetica-Bold').fontSize(10)
  headers.forEach((h, i) => {
    doc.text(h, startX + i * colWidth, currentY, { width: colWidth, align: 'left' })
  })

  currentY += 20
  doc.moveTo(startX, currentY - 5).lineTo(doc.page.width - 40, currentY - 5).stroke('#cccccc')

  // Body rows
  doc.font('Helvetica').fontSize(9)
  for (const row of rows) {
    if (currentY > doc.page.height - 60) {
      doc.addPage()
      currentY = 40
    }
    row.forEach((cell, i) => {
      doc.text(String(cell ?? ''), startX + i * colWidth, currentY, {
        width: colWidth,
        align: 'left',
      })
    })
    currentY += 18
  }

  doc.end()
}

reportsRouter.get('/reports/sales', (req, res) => {
  const from = typeof req.query.from === 'string' ? req.query.from.trim() : ''
  const to = typeof req.query.to === 'string' ? req.query.to.trim() : ''
  const format = req.query.format

  let sql = `
    SELECT 
      date(s.created_at) as date,
      COUNT(DISTINCT s.id) as invoices,
      COALESCE(SUM(s.total), 0) as revenue,
      COALESCE(SUM(s.discount), 0) as discount,
      COALESCE((
        SELECT SUM(si.profit) 
        FROM sale_items si 
        WHERE si.sale_id IN (
          SELECT id FROM sales s2 WHERE date(s2.created_at) = date(s.created_at)
        )
      ), 0) as profit
    FROM sales s
    WHERE 1=1
  `
  const params: any[] = []
  if (from) {
    sql += ` AND date(s.created_at) >= ?`
    params.push(from)
  }
  if (to) {
    sql += ` AND date(s.created_at) <= ?`
    params.push(to)
  }

  sql += ` GROUP BY date(s.created_at) ORDER BY date ASC`

  try {
    const rows = db.prepare(sql).all(...params) as any[]

    if (format === 'csv') {
      const headers = ['Date', 'Invoices', 'Revenue', 'Discount', 'Profit']
      const dataRows = rows.map((r) => [
        r.date,
        r.invoices,
        Number(r.revenue).toFixed(2),
        Number(r.discount).toFixed(2),
        Number(r.profit).toFixed(2),
      ])
      return exportCSV(res, 'sales-report.csv', headers, dataRows)
    }

    if (format === 'pdf') {
      const headers = ['Date', 'Invoices', 'Revenue', 'Discount', 'Profit']
      const dataRows = rows.map((r) => [
        r.date,
        r.invoices,
        `$${Number(r.revenue).toFixed(2)}`,
        `$${Number(r.discount).toFixed(2)}`,
        `$${Number(r.profit).toFixed(2)}`,
      ])
      return exportPDF(res, 'sales-report.pdf', 'Sales Report', headers, dataRows)
    }

    return ok(res, rows)
  } catch (e: any) {
    return err(res, 'failed to generate sales report', 500)
  }
})

reportsRouter.get('/reports/inventory', (req, res) => {
  const format = req.query.format

  const sql = `
    SELECT 
      name as product, 
      sku,
      current_stock,
      (current_stock * purchase_price) as stock_value,
      CASE
        WHEN current_stock <= 0 THEN 'out_of_stock'
        WHEN current_stock <= minimum_stock THEN 'low_stock'
        ELSE 'in_stock'
      END as status
    FROM products
    ORDER BY name ASC
  `

  try {
    const rows = db.prepare(sql).all() as any[]

    if (format === 'csv') {
      const headers = ['Product', 'SKU', 'Current Stock', 'Stock Value', 'Status']
      const dataRows = rows.map((r) => [
        r.product,
        r.sku,
        r.current_stock,
        Number(r.stock_value).toFixed(2),
        r.status,
      ])
      return exportCSV(res, 'inventory-report.csv', headers, dataRows)
    }

    if (format === 'pdf') {
      const headers = ['Product', 'SKU', 'Current Stock', 'Stock Value', 'Status']
      const dataRows = rows.map((r) => [
        r.product,
        r.sku,
        r.current_stock,
        `$${Number(r.stock_value).toFixed(2)}`,
        r.status,
      ])
      return exportPDF(res, 'inventory-report.pdf', 'Inventory Report', headers, dataRows)
    }

    return ok(res, rows)
  } catch (e: any) {
    return err(res, 'failed to generate inventory report', 500)
  }
})

reportsRouter.get('/reports/product-performance', (req, res) => {
  const from = typeof req.query.from === 'string' ? req.query.from.trim() : ''
  const to = typeof req.query.to === 'string' ? req.query.to.trim() : ''
  const format = req.query.format

  let sql = `
    SELECT 
      p.name as product,
      COALESCE(SUM(si.quantity), 0) as quantity_sold,
      COALESCE(SUM(si.quantity * si.selling_price), 0) as revenue,
      COALESCE(SUM(si.profit), 0) as profit
    FROM sale_items si
    JOIN sales s ON s.id = si.sale_id
    JOIN products p ON p.id = si.product_id
    WHERE 1=1
  `
  const params: any[] = []
  if (from) {
    sql += ` AND date(s.created_at) >= ?`
    params.push(from)
  }
  if (to) {
    sql += ` AND date(s.created_at) <= ?`
    params.push(to)
  }

  sql += ` GROUP BY p.id, p.name ORDER BY revenue DESC`

  try {
    const rows = db.prepare(sql).all(...params) as any[]

    if (format === 'csv') {
      const headers = ['Product', 'Quantity Sold', 'Revenue', 'Profit']
      const dataRows = rows.map((r) => [
        r.product,
        r.quantity_sold,
        Number(r.revenue).toFixed(2),
        Number(r.profit).toFixed(2),
      ])
      return exportCSV(res, 'product-performance.csv', headers, dataRows)
    }

    if (format === 'pdf') {
      const headers = ['Product', 'Quantity Sold', 'Revenue', 'Profit']
      const dataRows = rows.map((r) => [
        r.product,
        r.quantity_sold,
        `$${Number(r.revenue).toFixed(2)}`,
        `$${Number(r.profit).toFixed(2)}`,
      ])
      return exportPDF(res, 'product-performance.pdf', 'Product Performance', headers, dataRows)
    }

    return ok(res, rows)
  } catch (e: any) {
    return err(res, 'failed to generate product performance report', 500)
  }
})
