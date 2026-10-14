# PRD - Sistem Point of Sale (POS) Makanan

## Ringkasan Produk
Aplikasi kasir web untuk usaha makanan/minuman, dengan backend Go dan antarmuka yang cepat untuk mengelola menu, transaksi, stok, shift, serta ringkasan penjualan.

## Masalah
Kasir membutuhkan checkout yang cepat dan minim salah input. Pemilik usaha perlu mengetahui ketersediaan bahan/menu, pergerakan stok, dan penjualan secara terukur.

## Pengguna
- **Kasir:** membuka/menutup shift, memilih menu, menerima pembayaran, mencetak struk, dan merekonsiliasi kas.
- **Admin/Manager:** mengelola menu, kategori, harga, stok, hak akses, dan laporan.

## Fitur Inti
1. **Kasir dan checkout**
   - Pencarian menu cepat serta dukungan pemindai barcode bila produk menggunakannya.
   - Keranjang dengan kuantitas, catatan pesanan, varian, dan add-on.
   - Diskon per item/transaksi dan pajak.
   - Pembayaran Cash, QRIS, Transfer, atau EDC; untuk pembayaran non-tunai aplikasi mencatat metode dan status, bukan menyimpan data kartu atau kredensial gateway.
   - Tampilkan kembalian untuk pembayaran tunai, cetak struk thermal, dan sediakan opsi struk digital.
2. **Manajemen menu dan inventaris**
   - Kelola menu makanan/minuman, kategori, harga, ketersediaan, varian, dan add-on.
   - Catat stok masuk/keluar dan beri peringatan stok menipis.
3. **Shift kasir**
   - Buka/tutup shift dengan kas awal dan kas akhir serta ringkasan selisih.
4. **Laporan**
   - Ringkasan penjualan harian/bulanan, omzet, dan produk terlaris.

## Alur Kasir
Kasir membuka shift -> mencari/memilih menu -> menentukan varian/add-on dan kuantitas -> memeriksa keranjang/diskon -> memilih metode pembayaran -> server memvalidasi dan menyimpan checkout -> kasir melihat hasil dan mencetak/mengirim struk.

## Prinsip Data Transaksi
Backend memvalidasi harga, diskon, pajak, pembayaran, dan stok. Penyimpanan transaksi serta perubahan stok harus atomik agar stok tidak berkurang sebagian ketika checkout gagal.
