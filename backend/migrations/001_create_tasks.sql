
-- Luma v1.0
-- Migración 001: Gestión de tareas

CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY,

    title TEXT NOT NULL
        CHECK (length(trim(title)) BETWEEN 1 AND 200),

    description TEXT,

    priority TEXT NOT NULL DEFAULT 'media'
        CHECK (priority IN ('baja', 'media', 'alta')),

    status TEXT NOT NULL DEFAULT 'pendiente'
        CHECK (status IN ('pendiente', 'completada')),

    created_at TEXT NOT NULL DEFAULT
        (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),

    updated_at TEXT NOT NULL DEFAULT
        (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),

    completed_at TEXT,

    CHECK (
        (status = 'pendiente' AND completed_at IS NULL)
        OR
        (status = 'completada' AND completed_at IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_tasks_status_created
ON tasks(status, created_at DESC);
