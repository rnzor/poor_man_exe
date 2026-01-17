-- Users
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- SSH Public Keys  
CREATE TABLE IF NOT EXISTS public_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    fingerprint TEXT UNIQUE NOT NULL,
    key_data TEXT NOT NULL,
    comment TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Apps
CREATE TABLE IF NOT EXISTS apps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id),
    container_name TEXT,
    image TEXT DEFAULT 'alpine:latest',
    http_port INTEGER DEFAULT 80,
    is_public BOOLEAN DEFAULT FALSE,
    status TEXT DEFAULT 'running',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- App Shares (Allowed identities)
CREATE TABLE IF NOT EXISTS app_shares (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_id INTEGER REFERENCES apps(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(app_id, email)
);

-- Audit Log
CREATE TABLE IF NOT EXISTS audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    event TEXT NOT NULL,
    user_id INTEGER,
    app_name TEXT,
    source_ip TEXT,
    details TEXT
);
