CREATE TABLE IF NOT EXISTS users (
    phone_number TEXT PRIMARY KEY,
    verified BOOLEAN NOT NULL DEFAULT 0,
    display_name TEXT,
    email TEXT,
    notify_by_sms BOOLEAN NOT NULL DEFAULT 1,
    notify_by_email BOOLEAN NOT NULL DEFAULT 1,
    payment_methods TEXT
);

CREATE TABLE IF NOT EXISTS groups (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_by TEXT NOT NULL REFERENCES users(phone_number),
    default_currency TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS group_members (
    group_id TEXT NOT NULL REFERENCES groups(id),
    phone_number TEXT NOT NULL REFERENCES users(phone_number),
    display_name_override TEXT,
    PRIMARY KEY (group_id, phone_number)
);

CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_phone_number ON group_members(phone_number);
