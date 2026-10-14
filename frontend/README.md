# Frontend Rasa POS

Antarmuka kasir dan admin dibangun dengan Vue 3, TypeScript, dan Vite. Panduan fitur, persiapan PostgreSQL, konfigurasi environment, dan cara menjalankan seluruh aplikasi tersedia di [README utama](../README.md).

Untuk menjalankan frontend secara terpisah:

```powershell
npm install
npm run dev
```

Saat development, Vite meneruskan request `/api` ke backend Go di `http://127.0.0.1:8080`. Jalankan backend terlebih dahulu agar data menu dapat dimuat.
