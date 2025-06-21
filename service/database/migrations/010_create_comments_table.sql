CREATE TABLE IF NOT EXISTS comments (
    comment_id TEXT PRIMARY KEY CHECK (
        comment_id LIKE '________-____-____-____-____________'
    ),
    emoji TEXT NOT NULL CHECK (
        LENGTH (emoji) >= 1
        AND LENGTH (emoji) <= 10
    ),
    commented_at TEXT NOT NULL CHECK (
        commented_at LIKE "____-__-__T__:__:__Z" OR
        commented_at LIKE "____-__-__T__:__:__+__:__" OR
        commented_at LIKE "____-__-__T__:__:__-__:__"
    ),
    message_id TEXT NOT NULL CHECK (
        message_id LIKE '________-____-____-____-____________'
    ),
    user_id TEXT NOT NULL CHECK (
        message_id LIKE '________-____-____-____-____________'
    ),
    FOREIGN KEY (message_id) REFERENCES messages (message_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);
