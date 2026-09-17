import { execSync } from 'node:child_process'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

// Resolve paths from this file's own location, not process.cwd() — Tauri
// runs beforeDevCommand/beforeBuildCommand from wherever `tauri` itself was
// invoked (repo root, via the npm scripts), not from src-tauri/.
const srcTauriDir = path.dirname(fileURLToPath(import.meta.url))
const backendDir = path.join(srcTauriDir, '..', 'backend')

// Tauri sets TAURI_ENV_TARGET_TRIPLE while running before-dev/before-build
// commands; fall back to the local host triple when run manually.
const triple = process.env.TAURI_ENV_TARGET_TRIPLE || 'x86_64-pc-windows-msvc'
const ext = triple.includes('windows') ? '.exe' : ''
const outPath = path.join(srcTauriDir, 'binaries', `hardwarepos-backend-${triple}${ext}`)

execSync(`go build -o "${outPath}" ./cmd/server`, {
  cwd: backendDir,
  stdio: 'inherit',
})
