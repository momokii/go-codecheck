-- sqlite table schema

-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- User table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    is_setup_completed INTEGER NOT NULL DEFAULT 0, -- sqlite not have boolean as default datatype,
    session_token TEXT DEFAULT NULL,
    session_expired TEXT DEFAULT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

-- Repository table
CREATE TABLE IF NOT EXISTS repositories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    -- FOREIGN KEY (user_id) REFERENCES users(id), -- will implement later
    name TEXT NOT NULL,
    description TEXT,
    path TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Scan results table
CREATE TABLE IF NOT EXISTS scans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repository_id INTEGER NOT NULL,
    result TEXT,
    vulnerabilities INTEGER DEFAULT 0,
    status TEXT NOT NULL CHECK(status IN ('pending', 'in_progress', 'completed', 'error')),
    scan_time TEXT DEFAULT (datetime('now')),
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE
);

-- Create index for faster queries
CREATE INDEX IF NOT EXISTS idx_repository_id ON scans(repository_id);

-- Create default user data
INSERT INTO users (username, password, is_setup_completed) VALUES ('admin', '$2a$16$U2KzLZbEQnZcLFb2Oi.sAemuWWKZ4X3V8.gHw/ClEB9lKzvv3K896', 0);