CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY CHECK (
        user_id LIKE '________-____-____-____-____________'
    ),
    username TEXT NOT NULL CHECK (
        LENGTH (username) >= 3
        AND LENGTH (username) <= 16
    ),
    profile_picture TEXT CHECK (
        LENGTH (profile_picture) >= 11
        AND LENGTH (profile_picture) <= 255
    ),
    created_at TEXT NOT NULL CHECK (
        created_at LIKE "____-__-__T__:__:__Z" OR
        created_at LIKE "____-__-__T__:__:__+__:__" OR
        created_at LIKE "____-__-__T__:__:__-__:__"
    )
);

CREATE INDEX IF NOT EXISTS idx_users_username on users (username);
