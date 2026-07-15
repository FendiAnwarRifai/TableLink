# Frontend Website (Vue 3 + JavaScript Murni + Vite)

Halaman ini berisi petunjuk instalasi, cara menjalankan, dan penjelasan struktur folder untuk bagian **Frontend (Website)** dari project TableLink Backoffice Portal.

---

## 🚀 Cara Instalasi & Menjalankan Frontend

### 1. Prasyarat
Pastikan Anda sudah menginstal **Node.js** di komputer Anda.

### 2. Langkah Instalasi (Hanya sekali saja)
Buka terminal di dalam folder `/frontend`, lalu jalankan perintah berikut untuk menginstal semua library pendukung (seperti Axios, TailwindCSS, Lucide Icons, dll):
```bash
npm install
```

### 3. Cara Menjalankan Website
Setelah proses instalasi selesai, jalankan perintah berikut untuk menyalakan server lokal website:
```bash
npm run dev
```
Setelah berjalan, buka alamat berikut di browser Anda:
👉 **`http://localhost:5173/`**

---

## 📁 Struktur Folder dan Penjelasan File

Berikut adalah isi folder `/frontend` yang sangat sederhana dan bersih:

```
frontend/
├── src/
│   ├── presentation/
│   │   └── views/
│   │       ├── IngredientsView.vue   # Halaman kelola Bahan Makanan (CRUD + Axios + Pagination)
│   │       └── ItemsView.vue         # Halaman kelola Menu Makanan (CRUD + Axios + Relasi + Pagination)
│   ├── App.vue                       # Layout utama (Sidebar + Navigasi antar menu)
│   ├── main.js                       # Pintu masuk (Entry point) inisialisasi Vue 3
│   └── style.css                     # Desain minimalis putih-hitam & TailwindCSS
├── index.html                        # Halaman HTML dasar pembungkus Vue
├── package.json                      # Daftar library (dependencies) & script run
└── vite.config.js                    # Setelan server Vite & TailwindCSS plugin
```

### Penjelasan File Utama:
1.  **`IngredientsView.vue`**: Halaman untuk mengelola bahan makanan. Seluruh aksi seperti menambah bahan, mengubah data, pagination, dan menghapus (*soft-delete*) dikerjakan di file ini menggunakan Axios langsung ke URL `http://localhost:3000/api/ingredients`.
2.  **`ItemsView.vue`**: Halaman untuk mengelola menu makanan. Anda bisa menambah item baru, memasukkan harga, serta memilih bahan makanan apa saja yang digunakan untuk menu tersebut (menggunakan multi-select checkbox).
3.  **`App.vue`**: Mengatur tampilan navigasi samping (*sidebar*) untuk memudahkan Anda berpindah antara halaman *Ingredients* dan *Items*.
4.  **`main.js`**: File JavaScript pertama yang dibaca browser untuk me-mount aplikasi Vue 3 ke dalam `index.html`.
5.  **`style.css`**: Desain visual tema putih-hitam minimalis dengan garis siku-siku tegas.
