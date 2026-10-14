# Arsitektur Sistem POS Makanan

## Stack yang Dipilih
- **Frontend:** Vue 3, TypeScript strict, dan Vite.
- **Backend:** Go dengan HTTP API berbasis `net/http` agar fondasi mudah dipelajari dan tidak bergantung pada framework besar.
- **Database:** PostgreSQL.
- **Komunikasi:** REST API berformat JSON. Saat development, Vite meneruskan `/api` ke backend Go.
- **State keranjang:** Disimpan sementara di frontend; harga dan stok tetap diverifikasi ulang oleh backend saat checkout.

## Struktur Project yang Dituju
```text
backend/
  go.mod
  cmd/server/       # titik masuk aplikasi API
  internal/         # konfigurasi, domain, layanan, dan repository
frontend/
  src/
    components/pos/ # ProductGrid, Cart, Numpad, Receipt
    services/       # pemanggilan API, transaksi, dan printer
    types/          # Product, Transaction, CartItem, User
    views/          # halaman kasir dan pengelolaan
```

## Alur Checkout
1. Frontend mengambil menu dari API dan menyimpan pilihan sementara di keranjang.
2. Kasir memilih produk, varian/add-on, diskon, dan metode pembayaran.
3. Frontend mengirim item, kuantitas, pembayaran, dan kunci idempotensi ke API checkout.
4. Backend memvalidasi input, mengambil harga/stok terbaru, lalu memakai transaksi PostgreSQL untuk menyimpan transaksi, pembayaran, dan perubahan stok secara atomik.
5. Backend mengembalikan hasil transaksi dan data struk; frontend menampilkan status dan menyediakan cetak struk.

Backend adalah sumber kebenaran untuk harga, diskon, pajak, total, dan stok. Nilai uang disimpan sebagai bilangan bulat dalam rupiah (bukan float). Rahasia koneksi database dibaca dari environment, tidak disimpan di Git.

## Menjalankan Saat Project Sudah Diinisialisasi
Buka dua terminal dari folder project:

```powershell
# Terminal 1
cd backend
go run ./cmd/server
```

```powershell
# Terminal 2
cd frontend
npm install
npm run dev
```

URL development: frontend `http://localhost:5173`, API `http://localhost:8080`. Database `pos` dan role lokal `pos_app` sudah dibuat pada PostgreSQL Laragon; migrasi awal membuat tabel `products` dan `stock_movements`, lalu memasukkan menu contoh hanya jika tabel produk masih kosong. Role `pos_app` memakai aturan `trust` localhost bawaan Laragon, jadi hanya gunakan untuk development lokal; jangan expose PostgreSQL ke jaringan. API admin dilindungi Basic Auth dari environment dan server hanya bind ke `127.0.0.1`.

Atur environment backend pada Terminal 1 sebelum menjalankan server. Ganti placeholder koneksi dan password lokal, jangan commit nilai rahasia:

```powershell
cd backend
$env:DATABASE_URL = "postgres://pos_app@127.0.0.1:5432/pos?sslmode=disable"
$env:ADMIN_USERNAME = "admin"
$env:ADMIN_PASSWORD = "GantiDenganPasswordKuat"
go run ./cmd/server
```

Pada Terminal 2:

```powershell
cd frontend
npm run dev
```

Halaman Admin menyediakan katalog, arsip menu, batas stok menipis, dan penyesuaian stok dengan alasan yang tercatat. Stok MVP dihitung per menu/porsi, belum sebagai bahan baku/resep. Checkout dan laporan belum tersedia. Set `ADMIN_USERNAME` dan `ADMIN_PASSWORD` sendiri sebelum server dijalankan; endpoint admin sengaja gagal tertutup jika keduanya belum ada. Go pada lingkungan saat ini terdeteksi versi 1.19.1; perbarui ke versi yang masih didukung sebelum penggunaan operasional.
