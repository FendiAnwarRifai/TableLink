# Panduan Instalasi & Menjalankan Project

Panduan ringkas untuk menginstal dan menjalankan backend (Golang) dan frontend (Vue 3) project **TableLink Portal**.

---

## 🚀 Langkah 1: Persiapan Database (PostgreSQL)

1.  Nyalakan PostgreSQL di komputer Anda (misalnya menggunakan **Laragon**).
2.  Buat database baru bernama **`testdb`**.
3.  Import file database **`tablelink-test.sql`** (ada di folder root utama ini) ke database `testdb` tersebut.
4.  Pastikan konfigurasi koneksi di file **`backend/.env`** sudah benar.

---

## ⚙️ Langkah 2: Menjalankan Backend (Golang)

Backend terdiri dari dua server yang harus dinyalakan secara bersamaan:

### 1. Jalankan gRPC Server (Layanan Relasi Internal)
Buka terminal baru, masuk ke folder `/backend`, lalu jalankan:
```bash
go run cmd/grpc/main.go
```
*Server ini berjalan di port `50051`.*

### 2. Jalankan HTTP REST API Server (Fiber)
Buka terminal baru lagi, masuk ke folder `/backend`, lalu jalankan:
```bash
go run cmd/api/main.go
```
*Server ini berjalan di port `3000`.*

---

## 💻 Langkah 3: Menjalankan Frontend (Vue 3)

1.  Buka terminal baru, masuk ke folder `/frontend`.
2.  Install paket/library pendukung (hanya sekali saat pertama kali setup):
    ```bash
    npm install
    ```
3.  Jalankan server website lokal:
    ```bash
    npm run dev
    ```
4.  Buka browser Anda dan akses:
    👉 **`http://localhost:5173/`**

---

## 📁 Ringkasan Struktur Folder

*   **`/backend`**: Kode Go Server. Handler API REST ada di `cmd/api/main.go` dan Handler gRPC ada di `cmd/grpc/main.go`.
*   **`/frontend`**: Kode Vue 3 (JavaScript murni). Halaman UI berada di folder `src/presentation/views/`.
*   **`tablelink-test.sql`**: Data backup awal database PostgreSQL.
