-- 1. TABEL EVENTS
CREATE TABLE IF NOT EXISTS events (
    id            SERIAL       PRIMARY KEY,
    organizer_id  INTEGER      NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    title         VARCHAR(200) NOT NULL,
    description   TEXT,
    venue         VARCHAR(200) NOT NULL,
    category      VARCHAR(50)  NOT NULL, -- 'seminar', 'workshop', 'kompetisi'
    status        VARCHAR(20)  NOT NULL DEFAULT 'draft', -- 'draft', 'open', 'closed', 'cancelled'
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_events_status_created ON events (status, created_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_events_organizer ON events (organizer_id);

-- 2. TABEL SESSIONS
CREATE TABLE IF NOT EXISTS sessions (
    id                 SERIAL       PRIMARY KEY,
    event_id           INTEGER      NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name               VARCHAR(100) NOT NULL,
    start_time         TIMESTAMPTZ  NOT NULL,
    end_time           TIMESTAMPTZ  NOT NULL,
    total_capacity     INTEGER      NOT NULL CHECK (total_capacity > 0),
    available_capacity INTEGER      NOT NULL CHECK (available_capacity >= 0),
    version            INTEGER      NOT NULL DEFAULT 1,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_session_times CHECK (end_time > start_time),
    CONSTRAINT chk_capacity_logical CHECK (available_capacity <= total_capacity)
);

CREATE INDEX IF NOT EXISTS idx_sessions_event_id ON sessions (event_id);
CREATE INDEX IF NOT EXISTS idx_sessions_times ON sessions (start_time, end_time);
