# Backend Server (Golang - Fiber + gRPC)

Halaman ini berisi petunjuk instalasi, cara menjalankan, dan penjelasan struktur folder untuk bagian **Backend (Server & Database)** dari project TableLink Backoffice Portal.

---

## 🚀 Cara Instalasi & Menjalankan Backend

### 1. Prasyarat
*   Sudah menginstal **Golang** (Go).
*   Sudah menginstal dan menyalakan database **PostgreSQL** (misalnya menggunakan Laragon).

### 2. Setup Database
1.  Buka Laragon / PostgreSQL client Anda, buat database baru bernama **`testdb`**.
2.  Import skema dan data awal dari file **`tablelink-test.sql`** (berada di folder root project utama) ke dalam database `testdb` tersebut.
3.  Pastikan file **`.env`** di dalam folder `/backend` sudah berisi data koneksi PostgreSQL yang cocok:
    ```env
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=testdb
    DB_SSLMODE=disable
    ```

### 3. Cara Menjalankan Server
Backend menggunakan **dua server** yang harus dinyalakan bersamaan di terminal terpisah:

#### Terminal 1: Jalankan gRPC Server (Layanan Relasi Internal)
Buka terminal baru, masuk ke folder `/backend`, lalu jalankan:
```bash
go run cmd/grpc/main.go
```
*Server gRPC akan berjalan di port `50051`.*

#### Terminal 2: Jalankan HTTP REST Server (Layanan Utama API)
Buka terminal baru satu lagi, masuk ke folder `/backend`, lalu jalankan:
```bash
go run cmd/api/main.go
```
*Server API Fiber akan berjalan di port `3000`.*

---

## 📁 Struktur Folder dan Penjelasan File

Struktur folder backend sengaja dibuat rata (*flat*) agar ramah bagi pemula:

```
backend/
├── cmd/
│   ├── api/
│   │   └── main.go       # HTTP REST Server (Fiber) & CRUD Query Database
│   └── grpc/
│       └── main.go       # gRPC Server (Layanan relasi tabel tm_item_ingredient)
├── pb/
│   ├── relation.pb.go         # File otomatis hasil kompilasi gRPC
│   └── relation_grpc.pb.go    # File otomatis hasil kompilasi gRPC
├── proto/
│   └── relation.proto    # Kontrak perjanjian / rute gRPC
├── .env                  # Konfigurasi koneksi database & port server
├── go.mod                # Daftar modul library Go yang digunakan
└── go.sum                # Catatan checksum keamanan library Go
```

### Penjelasan File Utama:
1.  **`cmd/api/main.go`**: Berisi seluruh logika HTTP API (Fiber). Mulai dari rute, validasi keunikan nama, pemanggilan koneksi gRPC internal, hingga query SQL langsung ke tabel `tm_ingredient` dan `tm_item`.
2.  **`cmd/grpc/main.go`**: Berisi server gRPC yang bertugas menerima perintah internal dari Fiber untuk memasukkan, membaca, dan menghapus data relasi menu makanan di tabel `tm_item_ingredient`.
3.  **`proto/relation.proto`**: File kontrak protobuf yang menentukan format pengiriman data dan nama fungsi yang tersedia di gRPC.
4.  **`pb/`**: Hasil kompilasi yang otomatis dibuat dari file `.proto` di atas agar bahasa Go bisa membaca rute gRPC dengan mudah.
