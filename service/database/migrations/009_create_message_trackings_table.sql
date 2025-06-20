CREATE TABLE IF NOT EXISTS message_trackings (
    message_id TEXT NOT NULL CHECK (
        message_id LIKE '________-____-____-____-____________'
    ),
    user_id TEXT NOT NULL CHECK (
        user_id LIKE '________-____-____-____-____________'
    ),
    read_at TEXT CHECK (read_at LIKE "____-__-__T__:__:__Z"),
    PRIMARY KEY (message_id, user_id),
    FOREIGN KEY (message_id) REFERENCES messages (message_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE SET NULL
);
