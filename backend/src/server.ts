import express from 'express'
import cors from 'cors'
import { config } from './config.js'
import { initDatabase } from './database.js'
import { authRouter } from './routes/auth.js'
import { categoriesRouter } from './routes/categories.js'
import { suppliersRouter } from './routes/suppliers.js'
import { customersRouter } from './routes/customers.js'
import { productsRouter } from './routes/products.js'
import { inventoryRouter } from './routes/inventory.js'
import { purchasesRouter } from './routes/purchases.js'
import { salesRouter } from './routes/sales.js'
import { returnsRouter } from './routes/returns.js'
import { usersRouter } from './routes/users.js'
import { dashboardRouter } from './routes/dashboard.js'
import { reportsRouter } from './routes/reports.js'
import { settingsRouter } from './routes/settings.js'

// Initialize SQLite database schema
initDatabase()

const app = express()

app.use(cors())
app.use(express.json())

// Health check
app.get('/health', (_req, res) => {
  res.json({ status: 'ok' })
})

// Mount API routes
const api = express.Router()
api.use(authRouter)
api.use(categoriesRouter)
api.use(suppliersRouter)
api.use(customersRouter)
api.use(productsRouter)
api.use(inventoryRouter)
api.use(purchasesRouter)
api.use(salesRouter)
api.use(returnsRouter)
api.use(usersRouter)
api.use(dashboardRouter)
api.use(reportsRouter)
api.use(settingsRouter)

app.use('/api', api)

// Global error handler
app.use((err: any, _req: express.Request, res: express.Response, _next: express.NextFunction) => {
  console.error('Unhandled server error:', err)
  res.status(500).json({ error: err.message || 'Internal Server Error' })
})

app.listen(config.port, () => {
  console.log(`HardwarePOS JavaScript/Node API listening on :${config.port}`)
})
