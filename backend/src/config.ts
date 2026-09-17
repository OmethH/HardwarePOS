import path from 'node:path'
import { fileURLToPath } from 'node:url'
import dotenv from 'dotenv'

dotenv.config()

const __dirname = path.dirname(fileURLToPath(import.meta.url))

export const config = {
  port: parseInt(process.env.PORT || '8971', 10),
  dbPath: process.env.DB_PATH || path.resolve(__dirname, '../data/hardwarepos.db'),
  jwtSecret: process.env.JWT_SECRET || 'hardwarepos-dev-secret-change-in-production',
  migrationsPath: process.env.MIGRATIONS_PATH || path.resolve(__dirname, '../migrations'),
}
