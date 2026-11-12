# EventSeat

Backend REST API berbasis **Go** dan **Fiber v2** untuk manajemen pendaftaran seminar, workshop, dan kompetisi kampus berkapasitas terbatas, beserta pengelolaan tiket event, notifikasi asinkron, dan feedback peserta.

## Anggota Tim

| Nama | NIM |
| :--- | :--- |
| LUTHFI ALAN PERDANA (Git Master) | 434241052 |
| MUHAMMAD RAKA RAZZANI | 434241056 |
| SYIFA AULIA ALIMATUZZAHRA | 434241050 |

## Teknologi Utama

- **Bahasa & Framework:** Go, Fiber v2
- **Database Relasional:** PostgreSQL (`pgx/v5`) untuk data transaksional, kuota sesi, RBAC, dan tabel `outbox`
- **Database Dokumen:** MongoDB (`mongo-driver/v2`) untuk feedback peserta, balasan panitia, dan analitik rating/tag
- **Antrean & Worker:** Redis + Asynq untuk pengiriman e-tiket dan pengingat H-1 dengan pola *Transactional Outbox*
- **Dokumentasi & Deployment:** OpenAPI/Swagger, Docker, dan Docker Compose

## Fitur Utama

- **Autentikasi & RBAC:** Registrasi, login dengan rate limiter, JWT + refresh token rotation, serta kontrol akses berbasis peran (`user`, `staff`, `admin`) dan kepemilikan resource.
- **Manajemen Event & Sesi:** CRUD event dan sesi, validasi jadwal, *cursor pagination*, dan *optimistic locking* untuk pembaruan kapasitas sesi.
- **Pendaftaran & Tiket Atomik:** Pendaftaran sesi aman konkurensi (anti-overbooking), pembatalan, pindah sesi bebas deadlock, check-in kode tiket acak, dan ekspor daftar peserta (JSON/CSV).
- **Notifikasi Asinkron:** Pengiriman e-tiket dan pengingat H-1 melalui background worker (`-role=api|worker|all`) yang terpisah dari siklus request HTTP.
- **Feedback & Statistik:** Pengiriman ulasan/rating, balasan panitia, vote *helpful*, serta agregasi statistik rating dan tag populer di MongoDB.

## Peran Pengguna

| Role | Deskripsi |
| :--- | :--- |
| `user` | Peserta yang dapat melihat katalog event, mendaftar/pindah/membatalkan sesi, melihat tiket & notifikasi, serta mengirim feedback. |
| `staff` | Panitia acara yang mengelola event/sesi, melihat & mengekspor daftar peserta, melakukan check-in tiket, dan membalas feedback. |
| `admin` | Administrator dengan akses penuh, termasuk menghapus event, mengatur role pengguna, moderasi feedback, dan memantau antrean job. |

## Struktur Rencana Direktori

```text
uas-event-seat/
├── app/
│   ├── model/         # Entitas domain & DTO request/response
│   ├── repository/    # Akses data PostgreSQL & MongoDB
│   └── service/       # Business rules, service HTTP, relay & worker
├── config/            # Konfigurasi aplikasi, environment, logger, & queue
├── database/          # Koneksi pool PostgreSQL & MongoDB
├── docs/              # Dokumentasi Swagger/OpenAPI
├── helper/            # Utilitas stateless (response, validator, JWT, RBAC, cursor, CSV)
├── middleware/        # Middleware Fiber (auth, RBAC, CORS, rate limiter, JSON guard)
├── migrations/        # Skema migrasi SQL PostgreSQL
├── route/             # Definisi endpoint API & Swagger UI
├── tests/             # Integration test
└── tools/             # Alat uji beban (load test) konkurensi
```

## Status Proyek

Inisialisasi awal repositori (persiapan struktur dan dokumentasi proyek).
