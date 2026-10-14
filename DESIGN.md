# Design System POS Makanan

## Arah Visual
Tampilan bersih, kontras tinggi, dan cepat dipindai. Utamakan target sentuh besar serta alur keyboard agar kasir tidak perlu banyak berpindah fokus.

## Fokus UI Kasir
- **Navigasi kategori menu:** kategori makanan/minuman mudah dijangkau dan pencarian menu cepat; mendukung input dari barcode scanner bila diperlukan.
- **Grid produk:** kartu menu menampilkan foto, nama, harga, dan status ketersediaan. Produk habis tidak bisa dimasukkan ke keranjang.
- **Varian dan add-on:** pilihan seperti ukuran, level, atau tambahan ditampilkan jelas sebelum masuk keranjang.
- **Panel keranjang:** tetap terlihat di sisi kanan pada layar desktop/tablet; item, kuantitas, catatan pesanan, diskon, pajak, dan total mudah diperiksa.
- **Pembayaran:** nominal diterima, kembalian, serta metode Cash, QRIS, Transfer, atau EDC tampil jelas. Tombol aksi utama besar dan terpisah dari aksi berisiko.
- **Numpad:** tombol angka besar untuk input nominal tunai.
- **Status:** badge Paid, Pending, Failed, dan Refunded dengan teks yang terbaca, bukan warna saja.

## Area Admin
- Sediakan tabel menu yang dapat dipindai cepat, dengan harga, jumlah stok, batas stok menipis, dan status aktif/arsip terlihat jelas.
- Form produk dan penyesuaian stok harus ringkas; perubahan stok meminta alasan dan status stok menipis/habis tidak hanya dibedakan dengan warna.
- Akses Admin memiliki keadaan masuk, gagal autentikasi, memuat, sukses, dan kosong yang jelas.

## Responsif dan Perangkat
- Prioritaskan resolusi monitor kasir desktop/tablet; tetap usable di layar kecil.
- Sediakan fokus keyboard yang terlihat, urutan tab masuk akal, dan area klik yang cukup besar.
- Format cetak struk dengan `@media print` agar sesuai kertas thermal 58 mm dan 80 mm; sembunyikan elemen UI non-struk saat mencetak.
