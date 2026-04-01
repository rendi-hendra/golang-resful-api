# Issue: Pembuatan Rangkaian Unit Test untuk Semua Endpoint API

Dokumen ini berisi panduan untuk mengimplementasikan unit test (struktur integrasi HTTP) bagi semua endpoint API yang berjalan di aplikasi ini. Tugas Anda adalah menulis seluruh kode tes berdasarkan instruksi berikut menggunakan framework testing `go test` bawaan Golang.

## Instruksi & Standar Pengerjaan

1. **Lokasi File**: Seluruh rangkaian test wajib diletakkan secara terpisah di dalam direktori `test` di *root level* (misal: `test/user_test.go`).
2. **Inisialisasi & Reset Data Otomatis**: 
   - Anda harus menyiapkan fungsi *setup* untuk melakukan *bootstraping* aplikasi saat *testing* dimulai (membuat instance Fiber app dan mengkonekikannya ke DB testing).
   - **WAJIB**: Sebelum masing-masing *skenario* (fungsi test tunggal) dijalankan, Anda harus menjalankan skrip untuk menghapus/mengosongkan (*truncate*/*delete*) data di tabel `users` terlebih dahulu. Ini berguna agar semua skenario berjalan dalam keadaan bersih dan tidak saling tumpang tindih (*isolated & idempotent*).
3. **Pengecekan Komprehensif**: Di tiap skenario, pastikan assertion memeriksa *HTTP Status Code* dan *Response JSON Structure*.

---

## Daftar Skenario per API

Terdapat 5 endpoint utama pada aplikasi ini. Silakan buat fungsi tes per skenario selengkap mungkin dengan rincian berikut:

### 1. API Register (`POST /api/users`)
*   **Skenario Sukses**: Mengirim JSON request berisi properti `id`, `password`, `name`, dan `email` yang valid dan unik. Ekspektasi: Status 200 OK dan me-return profil yang divalidasi.
*   **Skenario Gagal - ID Duplikat**: Mengirim JSON dengan *ID* atau *Email* yang sudah ditambahkan pihak lain ke database di awal. Ekspektasi: Status 409 Conflict.
*   **Skenario Gagal - Data Kosong**: Mengirim *payload* kosong atau tidak melampirkan *mandatory fields*. Ekspektasi: Status 400 Bad Request.
*   **Skenario Gagal - Format Email Invalid**: Sengaja mengirimkan data `email` berformat bebas tanpa `@`. Ekspektasi: Status 400 Bad Request.

### 2. API Login (`POST /api/users/_login`)
*   **Skenario Sukses**: Mengirim kredensial username (`id`) dan `password` murni yang terdaftar. Ekspektasi: Status 200 OK dan mengembalikan token berupa `access_token` dan `refresh_token`.
*   **Skenario Gagal - User Tidak Ada**: Mengirim ID yang tidak ada di database. Ekspektasi: Status 404 Not Found.
*   **Skenario Gagal - Password Salah**: Mengirim `id` benar dengan kata sandi acak. Ekspektasi: Status 401 Unauthorized.
*   **Skenario Gagal - Validasi**: Mengirim request keliru/kosong. Ekspektasi: Status 400 Bad Request.

### 3. API Refresh Token (`POST /refresh-token`)
*   **Skenario Sukses**: Mengambil dan melampirkan `refresh_token` asli di dalam JSON *request body*. Ekspektasi: Status 200 OK dan menerima `access_token` edisi terbaru.
*   **Skenario Gagal - Token Salah/Palsu**: Mengirimkan string acak sebagai token, atau token milik entitas user lain. Ekspektasi: Status 401 Unauthorized.
*   **Skenario Gagal - Unprocessable Entity**: Mengirim payload kosong atau tipe data keliru. Ekspektasi: Status 400 Bad Request.

### 4. API Get Current User (`GET /api/users/_current`)
*   **Skenario Sukses**: Request ke endpoint dengan melampirkan _Header_ berisi `Authorization: Bearer <access_token>`. Ekspektasi: Status 200 OK dan mengembalikan data User Response spesifik milik user yang disematkan token-nya.
*   **Skenario Gagal - Token Kadaluwarsa/Dimanipulasi**: Mengirim Bearer token bodong (atau masa berlaku sudah habis / *tampered*). Ekspektasi: Status 401 Unauthorized.
*   **Skenario Gagal - Tanpa Token**: Melakukan *hit URL* sama sekali tanpa header `Authorization`. Ekspektasi: Status 401 Unauthorized.

### 5. API Update Current User (`PATCH /api/users/_current`)
*   **Skenario Sukses - Ganti Nama Panggilan**: Memperbarui satu properti (`name`) dengan diiringi *Header Bearer token*. Ekspektasi: Status 200 OK. Nama dalam response (dan Database) berubah.
*   **Skenario Sukses - Ganti Password**: Memperbarui kata sandi dan membuktikan nilainya sudah dienkripsi ulang (Opsional: tes untuk memastikan login gagal lewat kata sandi lama). Ekspektasi: Status 200 OK.
*   **Skenario Gagal - Parameter Typo/Salah Tipe**: Mengubah email menjadi tanpa karakter `@`. Ekspektasi: Status 400 Bad Request akibat gagal divalidasi _Struct Validator_.
*   **Skenario Gagal - Tidak Diotorisasi (Kosong)**: Hit URL merubah data namun tanpa header `Authorization` valid. Ekspektasi: Status 401 Unauthorized.

---

> **Note untuk implementator:** 
> - Tidak perlu mendesain arsitektur test. Langsung kerjakan fungsionalitas pengujiannya di dalam `.go` file.
> - Hindari menggunakan mock DB jika Anda bisa menggunakan DB test langsung yang siap pakai. Fokus saja pada alur *Black Box* HTTP Handler Testing (`ResponseRecorder` milik Go Fiber).
