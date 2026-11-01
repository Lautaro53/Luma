
-- Luma v1.0
-- Migración 002: Organización mediante carpetas

CREATE TABLE folders (
    id INTEGER PRIMARY KEY,

    name TEXT NOT NULL COLLATE NOCASE UNIQUE
        CHECK (length(trim(name)) BETWEEN 1 AND 80),

    created_at TEXT NOT NULL DEFAULT
        (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),

    updated_at TEXT NOT NULL DEFAULT
        (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- La asociación con una carpeta es opcional.
ALTER TABLE tasks
ADD COLUMN folder_id INTEGER
    REFERENCES folders(id)
    ON DELETE SET NULL;

-- Facilita consultar las tareas de una carpeta.
CREATE INDEX idx_tasks_folder
ON tasks(folder_id);
