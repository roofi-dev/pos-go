# AGENTS.md

## Sebelum Mulai
- Baca `PRD.md` untuk alur dan batas fitur POS makanan.
- Baca `ARCHITECTURE.md` untuk stack Go, Vue, PostgreSQL, struktur project, dan alur checkout.
- Baca `DESIGN.md` dan periksa komponen UI yang sudah ada sebelum membuat komponen baru.

## Aturan Umum
- Frontend menggunakan Vue 3 + TypeScript strict; hindari `any`.
- Backend menggunakan Go dan validasi transaksi/pembayaran di server.
- Buat komponen POS modular, cepat, dan ramah layar sentuh maupun keyboard.
- Backend menjadi sumber kebenaran harga, total, pembayaran, dan stok.
- Jangan mengubah skema database yang sudah dipakai tanpa konfirmasi.
- Simpan jumlah uang sebagai integer rupiah; jangan memakai floating point untuk perhitungan uang.

## Keamanan
- Jangan hardcode API key, password, atau kredensial payment gateway.
- Baca konfigurasi rahasia dari environment; jangan commit file `.env` berisi rahasia.
- Gunakan transaksi database untuk checkout dan perubahan stok terkait.

## Perintah Pengembangan
Jalankan perintah dari folder bagian terkait (setelah source code dibuat):
- Backend: `cd backend`, `go run ./cmd/server`, `go test ./...`
- Frontend: `cd frontend`, `npm install`, `npm run dev`, `npm run build`

API kasir dan admin memakai PostgreSQL. Tabel produk dan riwayat perubahan stok dibuat oleh migrasi awal yang idempotent; data contoh ditambahkan hanya jika tabel produk kosong. Admin menggunakan Basic Auth dari environment dan server development hanya bind ke loopback. Checkout belum tersedia. Jangan mengekspos server development ke jaringan publik.
