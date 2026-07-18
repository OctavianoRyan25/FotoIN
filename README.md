# FotoIN — Platform Fotografi Event Indonesia

Platform SaaS yang menghubungkan fotografer event (fokus olahraga & Sunmori) dengan peserta event. Fotografer unggah foto ke galeri, pembeli browse & beli via link unik.

## Struktur Proyek

```
FotoIN/
├── backend/          # Go REST API (Echo + GORM + PostgreSQL)
├── dashboard/        # React + Vite (Fotografer & Admin Dashboard) 
└── client/           # Next.js (Galeri publik & Checkout)
```

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go, Echo, GORM, PostgreSQL |
| Dashboard | React 19, Vite, Tailwind CSS v4, Radix UI |
| Client App | Next.js 16, Tailwind CSS |
| Storage | AWS S3 (S3-compatible) |
| Auth | JWT |
| Email | SMTP |
| Payment | Midtrans/Duitku QR Code Direct |

## Cara Menjalankan

### Prerequisites
- Go 1.24+
- Node.js 20+
- PostgreSQL

### 1. Backend
```bash
cd backend
cp .env.example .env
# Edit .env sesuai konfigurasi Anda

go run ./cmd/api/main.go
# Server berjalan di http://localhost:8080
```

### 2. Dashboard (Fotografer)
```bash
cd dashboard
npm install
npm run dev
# Buka http://localhost:5173
```

### 3. Client App (Galeri Publik)
```bash
cd client
npm install
npm run dev
# Buka http://localhost:3000
```

## Fitur MVP

### Fotografer
- ✅ Register/Login
- ✅ Dashboard dengan metrics (saldo, storage, galeri)
- ✅ Create/Edit/Delete galeri
- ✅ Upload foto langsung ke S3 via presigned URL
- ✅ Dompet & request penarikan dana

### Pembeli (Client)
- ✅ Akses galeri via URL unik (`/gallery/:slug`)
- ✅ Browse foto dengan watermark
- ✅ Cart & checkout (email → payment)
- ✅ Email download link setelah bayar

### Super Admin
- ✅ Manajemen pengguna (suspend/aktif)
- ✅ Revenue platform overview
- ✅ Approval/rejection withdrawal request

## Business Rules

| Rule | Value |
|---|---|
| Floor Price | Min Rp 10.000/foto (hard-coded backend) |
| Basic Storage | 3 GB |
| Basic Retention | 14 hari |
| Basic Commission | 20% |
| Premium Storage | 100 GB |
| Premium Commission | 4% |
| Premium Price | Rp 149.000/bulan |

## API Endpoints

```
POST /api/v1/auth/register
POST /api/v1/auth/login

GET  /api/v1/public/galleries/:slug
GET  /api/v1/public/events
GET  /api/v1/public/events/:slug
POST /api/v1/public/orders
POST /api/v1/webhooks/payment

GET  /api/v1/events              # Photographer auth
POST /api/v1/events
GET  /api/v1/events/:id
PUT  /api/v1/events/:id
DELETE /api/v1/events/:id

GET  /api/v1/galleries           # Photographer auth
POST /api/v1/galleries
GET  /api/v1/galleries/:id
PUT  /api/v1/galleries/:id
DELETE /api/v1/galleries/:id
POST /api/v1/galleries/:id/photos/upload-url
POST /api/v1/galleries/:id/photos/confirm
DELETE /api/v1/photos/:id

GET  /api/v1/wallet
POST /api/v1/wallet/withdrawal

GET  /api/v1/admin/users         # Super Admin only
GET  /api/v1/admin/revenue
GET  /api/v1/admin/withdrawals
PATCH /api/v1/admin/withdrawals/:id
```
