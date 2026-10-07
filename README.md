# SPKD MedCredix — Golang + Vue.js 3 + Vite (Versi 2.0)

**Medical Credentialing Intelligence Platform** — PT Sistem Pelayanan Kesehatan dan Data (SPKD)
"From Hospital Data to Action."

Rebuild penuh dari Next.js 16 + Prisma ke stack **Golang (backend) + Vue.js 3 + Vite (frontend)**.
Seluruh fitur, kontrak API, data demo, design system SPKD, dan hasil scoring engine dipertahankan
identik — Health Score **84,4** / BDCS **86,6** Stage 6 LULUS (sama dengan versi asli).

## Struktur Paket

| Path | Isi |
|---|---|
| `backend/` | API Golang — Gin + GORM + SQLite (pure-Go driver), 33 model, 21 endpoint, RBAC 10 peran, analytics engine (scoring), seeder data demo |
| `frontend/` | Vue 3.5 + Vite 7 + Tailwind CSS 4 + Pinia + chart.js — 29 view hash-router SPA, design token SPKD, dark mode |
| `frontend/src/` | Sumber frontend: `main.ts` (entry), `App.vue`, `components/`, `stores/`, `composables/`, `utils/`, `assets/` |
| `frontend/public/fonts/` | Poppins & Open Sans (self-hosted, subset latin) |
| `sandbox/` | Skrip menjalankan layanan di sandbox (opsional — lihat README di dalamnya) |

## Cara Menjalankan (Lokal)

### 1. Backend Golang (port 8081)

```bash
cd backend
go run ./cmd/server          # auto-migrate + seed data demo saat DB kosong
# atau: go build -o medcredix-api ./cmd/server && ./medcredix-api
```

- Database: SQLite di `backend/data/medcredix.db` (dibuat otomatis + seed 2 RS, 12 pengguna,
  38 praktisi, 91 kredensial, 48 klaim + 87 anomali, 12 risiko, 10 peringatan, 16 aksi, dll.)
- Environment opsional: `MEDCREDIX_PORT` (default 8081), `MEDCREDIX_DB` (path SQLite),
  `MEDCREDIX_STATIC` (direktori build frontend untuk mode single-binary).

### 2. Frontend Vue 3 + Vite (port 3000)

```bash
cd frontend
npm install                  # atau bun install / pnpm install
npm run dev                  # http://localhost:3000 — proxy /api → backend Golang :8081
```

- Alamat backend untuk proxy dev dapat diubah lewat env `MEDCREDIX_API_URL`
  (lihat `frontend/.env.example`).
- Mode produksi (single-binary Go): `npm run build` → hasil di `frontend/dist`, lalu jalankan
  backend dengan `MEDCREDIX_STATIC=../frontend/dist` sehingga frontend & API dilayani dari satu port.
- Frontend dan backend beda origin: build dengan `VITE_API_BASE=https://api.contoh.id npm run build`
  (CORS backend sudah mengizinkan header `X-User-Id`).
- Mode demo statis tanpa backend (mis. Vercel/Netlify): `npm run dev:demo` / `npm run build:demo`
  — seluruh `/api` dilayani di browser oleh `src/utils/demoApi.ts` dari `src/data/demo-db.json`.
- Cek tipe: `npm run typecheck` (vue-tsc).

### 3. Login Demo (kata sandi semua akun: `medcredix2026`)

| Peran | Email |
|---|---|
| SPKD Super Admin | superadmin@spkd.co.id |
| Admin RS Nusantara Sehat | admin@nusantarasehat.co.id |
| Manajemen RS (eksekutif) | manajemen@nusantarasehat.co.id |
| Petugas Kredensial | kredensial@nusantarasehat.co.id |
| Komite Medis | komite@nusantarasehat.co.id |
| Tim Klaim JKN | klaim@nusantarasehat.co.id |
| IT / SIMRS | it@nusantarasehat.co.id |
| Auditor Independen | auditor@spkd.co.id |
| Auditor BPJS | auditorbpjs@spkd.co.id |

Layar login menyediakan tombol quick-login per peran.

## Ringkasan Fitur (paritas penuh dengan versi Next.js)

- Dashboard Eksekutif: Hospital Health Score (9 domain) + BDCS Stage 0–7 dengan Trust Multiplier transparan
- Setiap skor drillable: "Mengapa Skor Ini?" → Kalkulasi/Kontributor/Bukti/Risiko/Rekomendasi/Tren
- Risk Radar, Early Warning Center, Action Center (closed-loop improvement)
- Credentialing (STR/SIP, verifikasi, prasyarat P1–P6), Hospital Readiness, Clinical & Claim Intelligence, Data Intelligence
- Evidence Center, Regulatory Compliance Matrix (termasuk Kepmenkes 84/2026 CKG), Audit Trail, Reports (export CSV/print)
- Multi-tenant 2 RS, RBAC 10 peran server-enforced, dark mode, pencarian global Ctrl+K

## Pemetaan Porting

| Asli (Next.js 16) | Rebuild ini |
|---|---|
| `src/app/api/*` (21 route) | `backend/internal/api/*.go` (Gin) |
| `prisma/schema.prisma` (33 model) | `backend/internal/model/model.go` (GORM) |
| `prisma/seed*.ts` | `backend/internal/seed/*.go` |
| `src/lib/scoring.ts` | `backend/internal/scoring/scoring.go` |
| `src/lib/auth.ts` (RBAC) | `backend/internal/auth/auth.go` |
| `src/components/medcredix/*` | `frontend/src/components/medcredix/*` (Vue SFC) |
| `src/components/ui/*` (shadcn) | `frontend/src/components/ui/*` (padanan ringan) |
| Zustand store + hash-router | Pinia store + hash-router (identik) |
| recharts | chart.js via vue-chartjs (ChartLine/ChartBar/ChartRadar) |
| next/font (Poppins/Open Sans) | self-hosted woff2 + @font-face |

## Migrasi Nuxt 4 → Vue 3 + Vite

| Nuxt 4 | Vue 3 + Vite |
|---|---|
| `nuxt.config.ts` (head, css, modules) | `index.html` (title/meta/favicon/font) + `src/main.ts` (Pinia, CSS) + `vite.config.ts` |
| `app/app.vue` | `src/App.vue` (dimount di `#app`) |
| Auto-import Vue API, composables, utils, store | Import eksplisit (`import { ref } from 'vue'`, `~/composables/useApi`, dst.) |
| Auto-register komponen + `resolveComponent()` | Import komponen eksplisit di tiap SFC |
| `@pinia/nuxt` | `pinia` (`createPinia()` di `main.ts`) |
| `nuxt generate` → `.output/public` | `vite build` → `dist/` |
| Mode demo statis (tanpa backend) | Default memanggil backend Golang; mode demo via `VITE_DEMO_MODE=true` |

Template, CSS (`src/assets/css/main.css`), font, design token, dan data demo tidak diubah —
tampilan identik dengan versi Nuxt.

## Catatan Teknis

- Kontrak JSON API identik (envelope `{ok, data}` / `{ok, error}`, field camelCase) sehingga
  frontend Vue dan klien lain dapat berganti backend tanpa perubahan.
- Sesi demo memakai header `X-User-Id` (sesuai versi asli); produksi disarankan JWT/OAuth2.
- Password hashing demo SHA-256+salt (sama dengan asli); produksi: bcrypt/argon2.
- Seluruh data demo fiktif dan berlabel DEMO — tidak ada data pasien nyata.
