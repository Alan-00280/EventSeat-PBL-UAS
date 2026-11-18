-- Keterangan: Data dummy pengguna untuk berbagai role sesuai helper/authz.go
-- Hash dibuat menggunakan bcrypt cost 12 sesuai helper/security.go

-- password untuk akun dummy di bawah ini: Password123!

INSERT INTO users (username, email, password, role, is_active, created_at, updated_at)
VALUES
    -- 1. Role: STAFF (Petugas Operasional Event & Check-in)
    (
        'staff_event',
        'staff@eventseat.com',
        '$2a$12$r5beRCkvZKslIfHMMco3Qe75elv6.L826UtL7AS4mZrVe6WnK1D5O',
        'staff',
        TRUE,
        NOW(),
        NOW()
    ),
    (
        'staff_ticketing',
        'ticketing@eventseat.com',
        '$2a$12$r5beRCkvZKslIfHMMco3Qe75elv6.L826UtL7AS4mZrVe6WnK1D5O',
        'staff',
        TRUE,
        NOW(),
        NOW()
    ),

    -- 2. Role: USER (Peserta / Pelanggan Aktif)
    (
        'budi_santoso',
        'budi@example.com',
        '$2a$12$r5beRCkvZKslIfHMMco3Qe75elv6.L826UtL7AS4mZrVe6WnK1D5O',
        'user',
        TRUE,
        NOW(),
        NOW()
    ),
    (
        'siti_aminah',
        'siti@example.com',
        '$2a$12$r5beRCkvZKslIfHMMco3Qe75elv6.L826UtL7AS4mZrVe6WnK1D5O',
        'user',
        TRUE,
        NOW(),
        NOW()
    ),

    -- 3. Role: USER (Contoh Akun Nonaktif untuk pengujian status login)
    (
        'user_inactive',
        'inactive@example.com',
        '$2a$12$r5beRCkvZKslIfHMMco3Qe75elv6.L826UtL7AS4mZrVe6WnK1D5O',
        'user',
        FALSE,
        NOW(),
        NOW()
    )
ON CONFLICT ((LOWER(username))) DO NOTHING;