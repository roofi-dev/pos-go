# Rasa POS

Aplikasi Point of Sale (POS) berbasis web untuk usaha makanan dan minuman. Proyek ini menggunakan Vue 3 dan TypeScript pada frontend, Go untuk REST API, serta PostgreSQL untuk data produk dan stok.

> **Status:** MVP dalam pengembangan. Pengelolaan menu dan stok sudah tersedia. Keranjang masih bersifat sementara; checkout dan pembayaran belum aktif.

## Fitur saat ini

- **Kasir:** melihat menu aktif dari API, mencari menu, memfilter kategori, menambahkan produk ke keranjang, mengubah jumlah, menghapus item, dan melihat subtotal.
- **Admin:** autentikasi Basic Auth, melihat katalog, menambah dan mengubah produk, mengarsipkan produk, serta menyesuaikan stok dengan alasan yang dicatat.
- **Inventaris:** penanda stok menipis berdasarkan batas tiap produk dan riwayat perubahan stok.
- **Data awal:** migrasi database dijalankan saat API mulai. Produk contoh ditambahkan hanya jika tabel produk masih kosong.

## Fitur yang belum tersedia

Checkout, pencatatan transaksi dan pembayaran, shift kasir, pencetakan struk, serta laporan penjualan belum diimplementasikan. Tombol checkout di keranjang saat ini memang belum aktif.

## Teknologi

- Frontend: Vue 3, TypeScript, Vite
- Backend: Go (`net/http`)
- Database: PostgreSQL
- Integrasi: REST API berformat JSON; Vite meneruskan request `/api` ke API Go saat development

## Persyaratan

- Go versi yang didukung; modul backend mendeklarasikan Go 1.19.
- Node.js yang kompatibel dengan Vite 7 (Node.js 20.19+ atau 22.12+) dan npm.
- PostgreSQL yang dapat diakses dari komputer lokal.

## Menjalankan secara lokal

### 1. Siapkan database

Siapkan database PostgreSQL bernama `pos` dan role yang dapat terhubung, menjalankan migrasi, serta membaca dan mengubah data produk dan stok. Migrasi awal akan membuat tabel `products` dan `stock_movements` saat backend dijalankan.

### 2. Jalankan backend

Buka terminal pertama dari direktori proyek:

```powershell
cd backend
$env:DATABASE_URL = "postgres://pos_app:PASSWORD_DATABASE@127.0.0.1:5432/pos?sslmode=disable"
$env:ADMIN_USERNAME = "admin"
$env:ADMIN_PASSWORD = "GANTI_DENGAN_PASSWORD_KUAT"
go run ./cmd/server
```

Ganti nilai koneksi dan password dengan konfigurasi lokal Anda. Jika password database berisi karakter khusus, encode password tersebut untuk URL. `ADMIN_USERNAME` dan `ADMIN_PASSWORD` diperlukan untuk menggunakan API Admin. Jangan commit kredensial ke repository.

API hanya bind ke `127.0.0.1:8080` dan menyediakan pemeriksaan kesehatan di `http://127.0.0.1:8080/api/health`.

### 3. Jalankan frontend

Buka terminal kedua dari direktori proyek:

```powershell
cd frontend
npm install
npm run dev
```

Buka URL yang ditampilkan Vite, biasanya `http://localhost:5173`. Frontend membutuhkan backend aktif di `127.0.0.1:8080` agar menu dapat dimuat.

## Pengujian dan build

Jalankan pengujian backend:

```powershell
cd backend
go test ./...
```

Periksa TypeScript dan build frontend:

```powershell
cd frontend
npm run build
```

## Endpoint API utama

| Method | Endpoint | Keterangan |
| --- | --- | --- |
| `GET` | `/api/health` | Status API |
| `GET` | `/api/products` | Daftar produk aktif untuk tampilan kasir |
| `GET` | `/api/admin/products` | Daftar produk untuk admin (Basic Auth) |
| `POST` | `/api/admin/products` | Membuat produk (Basic Auth) |
| `PUT` | `/api/admin/products/{id}` | Mengubah produk (Basic Auth) |
| `DELETE` | `/api/admin/products/{id}` | Mengarsipkan produk (Basic Auth) |
| `POST` | `/api/admin/products/{id}/stock` | Menyesuaikan stok dan mencatat alasan (Basic Auth) |

## Struktur proyek

```text
.
├── backend/
│   ├── cmd/server/                 # entry point API Go
│   └── internal/
│       ├── httpapi/                # handler REST API
│       ├── model/                  # model produk
│       └── store/                  # PostgreSQL, migrasi, dan akses data
├── frontend/
│   └── src/
│       ├── components/admin/        # pengelolaan katalog dan stok
│       ├── components/pos/          # grid menu dan keranjang
│       ├── services/                # pemanggilan API dan utilitas
│       └── types/                   # tipe TypeScript
├── ARCHITECTURE.md
├── DESIGN.md
└── PRD.md
```

## Keamanan dan batasan penggunaan

- Server development hanya mendengarkan loopback (`127.0.0.1`); jangan mengeksposnya langsung ke jaringan publik.
- Kredensial admin dibaca dari environment dan API Admin menggunakan Basic Auth. Gunakan hanya untuk development lokal; siapkan HTTPS dan mekanisme autentikasi yang sesuai sebelum deployment.
- Nilai harga disimpan sebagai integer rupiah. Stok MVP dihitung per produk/menu, bukan bahan baku atau resep.
- Untuk detail rancangan dan aturan pengembangan, lihat [ARCHITECTURE.md](ARCHITECTURE.md), [DESIGN.md](DESIGN.md), dan [PRD.md](PRD.md).
