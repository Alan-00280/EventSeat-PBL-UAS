# EventSeat

<p align="center">
  <a href="https://go.dev/" target="_blank"><img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go"></a>
  <a href="https://gofiber.io/" target="_blank"><img src="https://img.shields.io/badge/Fiber-v2-00ACD7?style=flat&logo=gofiber&logoColor=white" alt="Fiber v2"></a>
  <a href="https://www.postgresql.org/" target="_blank"><img src="https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql&logoColor=white" alt="PostgreSQL"></a>
  <a href="https://jwt.io/" target="_blank"><img src="https://img.shields.io/badge/JWT-v5-000000?style=flat&logo=jsonwebtokens&logoColor=white" alt="JWT"></a>
  <a href="https://redis.io/" target="_blank"><img src="https://img.shields.io/badge/Redis-7.x-DC382D?style=flat&logo=redis&logoColor=white" alt="Redis"></a>
  <a href="https://www.mongodb.com/" target="_blank"><img src="https://img.shields.io/badge/MongoDB-7.x-47A248?style=flat&logo=mongodb&logoColor=white" alt="MongoDB"></a>
</p>

Backend REST API berbasis **Go** dan **Fiber v2** untuk manajemen pendaftaran seminar, workshop, dan kompetisi kampus berkapasitas terbatas, beserta pengelolaan tiket event, notifikasi asinkron, dan feedback peserta.

---

## Informasi Akademik

| Kode Kelas | Nama Mata Kuliah | Semester |
| :--- | :--- | :--- |
| **TI-C4** | Pemrograman Backend Lanjut (Praktikum) | Semester 5 |

### Dosen Pengampu
**Christiant Dimas Renggana, S.Kom**

### Anggota Tim
| Nama | NIM | Peran Utama |
| :--- | :--- | :--- |
| **LUTHFI ALAN PERDANA** | 434241052 | Git Master & Core Developer |
| **MUHAMMAD RAKA RAZZANI** | 434241056 | Developer |
| **SYIFA AULIA ALIMATUZZAHRA** | 434241050 | Developer |

---

## Teknologi Utama

- **Bahasa & Framework:** Go 1.23+, Fiber v2
- **Database Relasional:** PostgreSQL (`pgx/v5` connection pool)
- **Database Dokumen (Rencana):** MongoDB (`mongo-driver/v2`) untuk feedback dan rating
- **Antrean & Cache (Rencana):** Redis + Asynq untuk pengiriman e-tiket dan Transactional Outbox
- **Keamanan & Autentikasi:** JWT (`golang-jwt/jwt/v5`), Bcrypt password hashing, Common Password Dictionary Checker, Rate Limiter

---

## Panduan Development Lokal

Ikuti langkah-langkah di bawah ini untuk menjalankan backend di mesin lokal.

### 1. Prasyarat Sistem
Pastikan telah menginstal:
- [Go](https://go.dev/dl/) versi 1.23 atau lebih baru
- [PostgreSQL](https://www.postgresql.org/download/) (v14 atau lebih baru)
- Git

### 2. Kloning Repositori & Masuk ke Direktori
```bash
git clone https://github.com/Alan-00280/EventSeat-PBL-UAS.git
cd EventSeat-PBL-UAS
```

### 3. Konfigurasi Environment (`.env`)
Salin file `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```
Sesuaikan konfigurasi kredensial PostgreSQL dan JWT pada file `.env`:
```env
APP_NAME="app-backend-event-seat"
APP_PORT=3000
APP_ENV=DEVELOPMENT

# Konfigurasi Database PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password_postgres_anda
DB_NAME=eventseat_db
DB_SSLMODE=disable
DB_MAX_CONNS=10

LOG_LEVEL=info

# Konfigurasi JWT
JWT_SECRET=rahasia-jwt-anda-yang-aman-dan-panjang
JWT_ISSUER=backend-event-seat
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_DAYS=7

ALLOWED_ORIGINS=http://localhost:5173
```

> [!NOTE]
> Pastikan database dengan nama yang ditentukan pada `DB_NAME` (contoh: `eventseat_db`) sudah dibuat di PostgreSQL Anda sebelum menjalankan migrasi.

### 4. Eksekusi Migrasi & Seeding Database
Jalankan file migrasi SQL yang berada di direktori `migrations/` secara berurutan pada database Anda:

1. **`migrations/001_users.sql`** (Tabel pengguna dan indeks case-insensitive)
2. **`migrations/002_auth.sql`** (Tabel refresh token dan indeks)
3. **`migrations/003_user_dummy_seed.sql`** (Data pengguna awal untuk testing)

**Contoh eksekusi menggunakan `psql`:**
```bash
psql -U postgres -d eventseat_db -f migrations/001_users.sql
psql -U postgres -d eventseat_db -f migrations/002_auth.sql
psql -U postgres -d eventseat_db -f migrations/003_user_dummy_seed.sql
```
*(Atau Anda bisa mengeksekusi script SQL di atas menggunakan GUI seperti DBeaver, pgAdmin, atau Navicat).*

### 5. Memastikan File Kamus Kata Sandi
Aplikasi membutuhkan file kamus kata sandi lemah `files/common_password.txt` untuk validasi `strongpassword`. Pastikan file tersebut ada pada direktori:
```text
files/common_password.txt
```

### 6. Menjalankan Aplikasi
Unduh dependensi Go dan jalankan server:
```bash
go mod tidy
go run main.go
```
Jika berhasil, server akan aktif pada port yang ditentukan (default: `http://localhost:3000`).

---

## Akun Uji Coba (Dummy Seed)

Untuk mempermudah pengujian lokal, migrasi `003_user_dummy_seed.sql` menyediakan beberapa akun dengan kata sandi bawaan: **`Password123!`**

| Username | Email | Role | Status | Catatan |
| :--- | :--- | :--- | :--- | :--- |
| `staff_event` | `staff@eventseat.com` | `staff` | Aktif | Petugas Operasional Event |
| `staff_ticketing` | `ticketing@eventseat.com` | `staff` | Aktif | Petugas Tiket & Check-in |
| `budi_santoso` | `budi@example.com` | `user` | Aktif | Peserta Acara |
| `siti_aminah` | `siti@example.com` | `user` | Aktif | Peserta Acara |
| `user_inactive` | `inactive@example.com` | `user` | **Nonaktif** | Untuk pengujian login akun nonaktif |

---

## Endpoint API yang Siap Dicoba

Berikut adalah daftar endpoint yang telah diimplementasikan dan dapat diuji melalui Postman, Bruno, Hoppscotch, atau cURL.

Base URL: `http://localhost:3000/api/v1`

### 1. Health Check
Memeriksa status operasional web server dan konektivitas pool PostgreSQL.
- **Method:** `GET`
- **URL:** `/api/v1/health`
- **Contoh Request:**
  ```bash
  curl -X GET http://localhost:3000/api/v1/health
  ```
- **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "server and database is OK!"
  }
  ```

---

### 2. Registrasi Pengguna Baru
Mendaftarkan akun baru dengan role default `user`.
- **Method:** `POST`
- **URL:** `/api/v1/auth/register`
- **Headers:** `Content-Type: application/json`
- **Body JSON:**
  ```json
  {
    "username": "ahmad_fauzi",
    "email": "ahmad@example.com",
    "password": "Password123!"
  }
  ```
  *(Catatan: Password harus minimal 8 karakter, mengandung huruf & angka, serta bukan kata sandi umum).*
- **Response (201 Created):**
  ```json
  {
    "success": true,
    "message": "user berhasil didaftarkan",
    "data": {
      "id": 6,
      "username": "ahmad_fauzi",
      "email": "ahmad@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2026-11-18T10:00:00Z"
    }
  }
  ```

---

### 3. Login Pengguna
Melakukan autentikasi dan menerbitkan pasangan token (`access_token` dan `refresh_token`). Terdapat rate limiter bawaan untuk keamanan brute-force.
- **Method:** `POST`
- **URL:** `/api/v1/auth/login`
- **Headers:** `Content-Type: application/json`
- **Body JSON:**
  ```json
  {
    "username": "budi_santoso",
    "password": "Password123!"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "login berhasil",
    "data": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6...",
      "refresh_token": "a1b2c3d4e5f6...",
      "token_type": "Bearer",
      "expired_in": 604800
    }
  }
  ```

---

### 4. Refresh Token (Rotasi Token)
Memperbarui access token yang telah kadaluarsa dengan menukarkan refresh token aktif. Refresh token lama akan di-*revoke* secara otomatis.
- **Method:** `POST`
- **URL:** `/api/v1/auth/refresh`
- **Headers:** `Content-Type: application/json`
- **Body JSON:**
  ```json
  {
    "refresh_token": "ISI_REFRESH_TOKEN_DARI_LOGIN"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "berhasil membuat token baru",
    "data": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6...",
      "refresh_token": "token_baru_hasil_rotasi...",
      "token_type": "Bearer",
      "expired_in": 604800
    }
  }
  ```

---

### 5. Logout Pengguna
Mencabut status aktif refresh token sehingga tidak dapat digunakan kembali.
- **Method:** `POST`
- **URL:** `/api/v1/auth/logout`
- **Headers:** `Content-Type: application/json`
- **Body JSON:**
  ```json
  {
    "refresh_token": "ISI_REFRESH_TOKEN_DARI_LOGIN"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "berhasil logout"
  }
  ```

---

### 6. Profil Pengguna Saat Ini (`Me`)
Mengambil data pengguna yang terautentikasi melalui header JWT.
- **Method:** `GET`
- **URL:** `/api/v1/auth/me`
- **Headers:**
  - `Authorization: Bearer <ACCESS_TOKEN>`
- **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "user berhasil ditemukan",
    "data": {
      "user": {
        "id": 3,
        "username": "budi_santoso",
        "email": "budi@example.com",
        "role": "user",
        "is_active": true,
        "created_at": "2026-11-18T09:00:00Z"
      }
    }
  }
  ```

---

## Struktur Direktori Saat Ini

```text
uas-event-seat/
├── app/
│   ├── model/         # Entitas domain, struct DTO request/response, dan definisi error
│   ├── repository/    # Akses data PostgreSQL (UserRepository, TokenRepository, sentinel error)
│   └── service/       # Layanan logika bisnis (AuthService)
├── config/            # Inisialisasi Fiber app, error handler, dan parsing variabel lingkungan (.env)
├── database/          # Inisialisasi koneksi pgxpool PostgreSQL
├── files/             # File pendukung (common_password.txt)
├── helper/            # Utilitas helper (JWT manager, RBAC, hash password, validator, errors)
├── middleware/        # Middleware Fiber (RequireAuth, RateLimiter, CORS, JSON guard)
├── migrations/        # Skrip DDL database PostgreSQL & dummy seed
├── routes/            # Pendaftaran rute API v1
└── main.go            # Entrypoint utama aplikasi
```
