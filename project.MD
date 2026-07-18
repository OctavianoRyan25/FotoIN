# Project Context: Event & Sunmori Photography SaaS

## 1. Project Overview

Platform SaaS (Software as a Service) yang menghubungkan fotografer event (fokus pada olahraga dan Sunmori) dengan peserta event. Platform memungkinkan fotografer mengunggah foto, membuat galeri, dan menjual foto secara digital. Peserta event dapat mencari dan membeli foto mereka.

## 2. Tech Stack (MVP Phase)

- **Backend:** Golang (REST API)
- **Frontend:** React.js (Web App untuk Dashboard Fotografer/Admin & Halaman Klien)
- **Database:** PostgreSQL (Relational Data)
- **Storage:** AWS S3 (atau S3-compatible Object Storage)
- **Out of Scope for MVP:** Message Broker/Queue (RabbitMQ/Kafka), Caching layer (Redis), dan Custom Machine Learning (Face/BIB scan). Ini akan diimplementasikan pada Fase 2.

## 3. Design System & UI/UX Guidelines

Platform ini memiliki **dua wajah visual yang berbeda**, disesuaikan dengan target user masing-masing.

### A. Dashboard (Photographer & Super Admin) — shadcn/ui Style

Dashboard adalah tool kerja (utility-first), jadi harus terasa rapi, fungsional, dan familiar bagi power user.

- **Component Library:** `shadcn/ui` (di atas Radix UI + Tailwind CSS).
- **Styling base:** gunakan preset default shadcn — `Inter` sebagai font, radius `0.5rem`, palet netral (`zinc` atau `slate`) sebagai base color.
- **Layout pattern:**
  - Sidebar navigation (collapsible) + topbar berisi breadcrumb, user menu, notifikasi.
  - Gunakan `Card` untuk metric widgets (Total Penjualan, Sisa Storage, dsb).
  - Gunakan `Table` (dengan `DataTable` pattern dari shadcn — sorting, pagination, filter) untuk daftar galeri, foto, transaksi, dan withdrawal.
  - Gunakan `Dialog`/`Sheet` untuk form create/edit (Create Gallery, Upload Photo, Request Withdrawal), bukan halaman terpisah, agar workflow tetap cepat.
  - Gunakan `Tabs` untuk pemisahan section di halaman detail galeri (Photos, Sales, Settings).
  - Gunakan `Badge` untuk status (Pending, Approved, Free Tier, Premium).
  - Gunakan `Toast` (via `sonner` atau `use-toast` shadcn) untuk notifikasi aksi (upload sukses, withdrawal diajukan, dsb).
- **Dark mode:** wajib didukung sejak awal (shadcn punya built-in theming via CSS variables) — banyak fotografer bekerja malam hari saat proses seleksi foto.
- **Data density:** dashboard boleh lebih padat informasi dibanding client app — prioritaskan efisiensi kerja, bukan estetika minimalis.
- **Role differentiation:** Super Admin dan Photographer memakai shell/komponen yang sama, tapi sidebar menu & permission berbeda (Super Admin punya menu tambahan: User Management, Platform Revenue, Subscription Management).

### B. Client / Buyer App — Simple & Elegant Style

Client app diakses oleh peserta event awam (bukan tech-savvy), sering dari mobile saat event berlangsung. Fokus: **kecepatan, kejelasan, dan kepercayaan (trust)** untuk transaksi pembayaran.

- **Visual tone:** minimalis, banyak whitespace, tipografi besar dan jelas, foto sebagai hero content (bukan UI chrome yang ramai).
- **Layout pattern:**
  - Grid galeri foto full-bleed dengan lazy loading & skeleton loading state.
  - Watermark halus (tidak mengganggu preview, tapi cukup untuk mencegah screenshot-abuse).
  - Sticky cart/bottom bar di mobile (menampilkan jumlah foto terpilih + total harga + tombol checkout) agar user tidak kehilangan konteks saat scroll galeri panjang.
  - Checkout flow **satu halaman** (email → payment) — hindari multi-step form yang membuat drop-off.
  - Empty state & error state didesain khusus (bukan default browser alert), karena user awam mudah bingung.
- **Color & typography:** palet netral + satu accent color (mis. warna brand platform), font yang mudah dibaca di layar kecil (mis. `Inter` atau `Plus Jakarta Sans`).
- **Mobile-first:** wajib, karena besar kemungkinan akses lewat link WhatsApp/Instagram Story dari HP.
- **Tidak perlu shadcn di sini** — komponen boleh custom-built lebih ringan (tanpa Radix), fokus ke performa loading (LCP) karena berisi banyak gambar.

## 4. Business Logic & Monetization (Freemium Model)

Aplikasi menggunakan model berlangganan untuk fotografer dengan dua tier:

### A. Basic Plan (Gratis)

- **Storage Limit:** Maksimal 3 GB.
- **Data Retention:** Galeri dan foto otomatis dihapus setelah 14 hari dari tanggal pembuatan.
- **Platform Fee/Komisi:** 20% dari harga jual foto.
- **Feature:** Standar galeri, manual tag/search.

### B. Premium Plan (Berbayar: Rp 149.000/bulan)

- **Storage Limit:** 100 GB.
- **Data Retention:** Selama langganan aktif.
- **Platform Fee/Komisi:** 4% dari harga jual foto (mencakup 2% PG Fee + 2% Platform Margin).
- **Feature:** Prioritas support, unmetered bandwidth limit.

### C. Transaksi & Aturan Harga

- Harga jual foto ditentukan secara bebas oleh Fotografer per galeri/foto.
- **CRITICAL RULE:** Harga jual minimal (Floor Price) dibatasi secara hard-coded di backend: **Minimum Rp 10.000 per foto**. Hal ini untuk mencegah kerugian platform akibat biaya egress S3.

## 5. MVP Core Features

### Role: Photographer (Seller)

1. **Authentication:** Register/Login (Email & Password).
2. **Dashboard Management:**
   - Melihat metrik sederhana (Total Penjualan, Sisa Storage S3).
   - Subscription Menu: Upgrade ke Premium (Simulasi/Integrasi PG).
3. **Gallery Management:**
   - Create Event/Gallery (Judul, Tanggal, Lokasi, Harga per foto).
   - Upload Foto langsung ke S3 (via Presigned URL atau Backend Handler).
   - Auto-compression di backend: Ukuran file yang disajikan ke publik/diunduh maksimal dikompresi agar hemat egress cost (target: ~2MB–5MB per file).
4. **Wallet/Withdrawal:**
   - Melihat saldo pendapatan bersih (setelah dipotong komisi platform sesuai Tier aktif).
   - Request Withdrawal (Fase MVP bisa dicatat statusnya menjadi 'Pending' untuk diproses manual).

### Role: Client/Buyer (Guest)

1. **Gallery Access:** Mengakses galeri menggunakan URL unik.
2. **Browsing:** Melihat foto dengan watermark platform (di-generate saat upload atau on-the-fly via CSS/Canvas di frontend untuk MVP).
3. **Checkout Flow:**
   - Menambahkan foto ke Cart.
   - Mengisi email (untuk pengiriman link download).
   - Membayar via Payment Gateway.
4. **Post-Purchase:** Menerima email berisi Presigned URL AWS S3 berbatas waktu untuk mengunduh foto beresolusi penuh tanpa watermark.

### Role: Super Admin (tambahan, mengikuti style shadcn dashboard)

1. **User Management:** Melihat/mengelola daftar fotografer (suspend, verify, dsb).
2. **Platform Revenue Overview:** Total komisi terkumpul, breakdown per tier.
3. **Subscription Management:** Melihat status langganan aktif, riwayat pembayaran Premium.
4. **Withdrawal Approval:** Antrian withdrawal request dari fotografer untuk diproses manual (approve/reject).

## 6. Development Guidelines for AI

- Terapkan **Clean Architecture** pada Golang (Handler/Controller, Usecase/Service, Repository) dengan pendekatan domain-driven.
- Gunakan **GORM** untuk interaksi ke PostgreSQL.
- Semua kalkulasi finansial (potongan komisi, total harga) **WAJIB** dilakukan di Backend (Golang) — jangan pernah mempercayai input nominal dari Frontend.
- Gunakan **DB Transactions** saat memproses status pembayaran dan update saldo dompet fotografer.
- Rancang skema database (DDL) dengan mempertimbangkan struktur untuk MVP, namun sediakan ruang berekspansi (misal: penambahan kolom `face_vector_id` pada tabel `photos` untuk keperluan ML di Fase 2).
- Untuk frontend dashboard, gunakan struktur folder yang memisahkan komponen shadcn (`components/ui`) dari komponen bisnis (`components/features`), mengikuti konvensi standar shadcn CLI (`npx shadcn add ...`).
- Untuk client app, pertimbangkan Next.js (jika belum fixed di React CRA/Vite) demi SEO galeri publik dan performa image loading (`next/image`), meskipun ini opsional tergantung keputusan tim.
