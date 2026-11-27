-- 1. TABEL REGISTRATIONS
CREATE TABLE IF NOT EXISTS registrations (
    id            SERIAL       PRIMARY KEY,
    user_id       INTEGER      NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    session_id    INTEGER      NOT NULL REFERENCES sessions(id) ON DELETE RESTRICT,
    ticket_code   VARCHAR(32)  NOT NULL UNIQUE,
    status        VARCHAR(20)  NOT NULL DEFAULT 'active', -- 'active', 'cancelled', 'checked_in'
    checked_in_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- PARTIAL UNIQUE INDEX (Kunci pencegah duplikasi sekaligus izin daftar ulang pasca-batal)
CREATE UNIQUE INDEX IF NOT EXISTS uq_user_session_active 
    ON registrations (user_id, session_id) 
    WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_registrations_user ON registrations (user_id);
CREATE INDEX IF NOT EXISTS idx_registrations_session ON registrations (session_id);
CREATE INDEX IF NOT EXISTS idx_registrations_ticket_code ON registrations (ticket_code);

-- 2. TABEL TRANSACTIONAL OUTBOX
CREATE TABLE IF NOT EXISTS outbox (
    id             BIGSERIAL    PRIMARY KEY,
    aggregate_type VARCHAR(50)  NOT NULL, -- 'registration', 'session'
    aggregate_id   INTEGER      NOT NULL,
    event_type     VARCHAR(100) NOT NULL, -- 'ticket.send', 'reminder.h1', 'registration.cancelled'
    payload        JSONB        NOT NULL,
    status         VARCHAR(20)  NOT NULL DEFAULT 'pending', -- 'pending', 'processed', 'failed'
    attempts       INTEGER      NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    processed_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_outbox_pending 
    ON outbox (status, created_at ASC) 
    WHERE status = 'pending';

-- 3. TABEL AUDIT NOTIFIKASI
CREATE TABLE IF NOT EXISTS notifications (
    id         BIGSERIAL    PRIMARY KEY,
    outbox_id  BIGINT       NOT NULL UNIQUE,
    user_id    INTEGER      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel    VARCHAR(20)  NOT NULL DEFAULT 'email',
    recipient  VARCHAR(255) NOT NULL,
    subject    VARCHAR(255) NOT NULL,
    body       TEXT         NOT NULL,
    status     VARCHAR(20)  NOT NULL DEFAULT 'sent',
    sent_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications (user_id, sent_at DESC);
