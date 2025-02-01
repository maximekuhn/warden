CREATE TABLE auth (
    user_id TEXT PRIMARY KEY,
    email TEXT NOT NULL,
    hashed_password BLOB NOT NULL,
    created_at DATE NOT NULL,
    session_id TEXT,
    session_expire_date DATE NOT NULL,
    UNIQUE(email)
);
