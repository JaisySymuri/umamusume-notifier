package storage

const createPointSystemsTable = `
CREATE TABLE IF NOT EXISTS point_systems (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,

    max_points INTEGER NOT NULL,
    current_points INTEGER NOT NULL,

    regen_minutes INTEGER NOT NULL,

    elapsed_seconds INTEGER NOT NULL DEFAULT 0,

    last_tick DATETIME NOT NULL
);
`

const createReminderStatesTable = `
CREATE TABLE IF NOT EXISTS reminder_states (
    system_id TEXT PRIMARY KEY
        REFERENCES point_systems(id),

    alert_sent BOOLEAN NOT NULL DEFAULT FALSE,

    full_sent BOOLEAN NOT NULL DEFAULT FALSE,

    last_message_id INTEGER
);
`

const loadPointSystemsQuery = `
SELECT
    id,
    name,
    max_points,
    current_points,
    regen_minutes,
    elapsed_seconds,
    last_tick
FROM point_systems
ORDER BY id;
`

const savePointSystemQuery = `
UPDATE point_systems
SET
    current_points = ?,
    elapsed_seconds = ?,
    last_tick = ?
WHERE id = ?;
`

const loadReminderStatesQuery = `
SELECT
    system_id,
    alert_sent,
    full_sent,
    last_message_id
FROM reminder_states
ORDER BY system_id;
`

const saveReminderStateQuery = `
INSERT INTO reminder_states (
    system_id,
    alert_sent,
    full_sent,
    last_message_id
)
VALUES (?, ?, ?, ?)
ON CONFLICT(system_id) DO UPDATE SET
    alert_sent = excluded.alert_sent,
    full_sent = excluded.full_sent,
    last_message_id = excluded.last_message_id;
`