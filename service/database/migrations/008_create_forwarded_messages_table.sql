CREATE TABLE IF NOT EXISTS forwarded_messages (
    forwarded_message_id TEXT PRIMARY KEY CHECK (
        forwarded_message_id LIKE '________-____-____-____-____________'
    ),
    forwarded_at TEXT NOT NULL CHECK (
        forwarded_at LIKE "____-__-__T__:__:__Z" OR
        forwarded_at LIKE "____-__-__T__:__:__+__:__" OR
        forwarded_at LIKE "____-__-__T__:__:__-__:__"
    ),
    original_message_id TEXT NOT NULL CHECK (
        original_message_id LIKE '________-____-____-____-____________'
    ),
    conversation_id TEXT NOT NULL CHECK (
        conversation_id LIKE '________-____-____-____-____________'
    ),
    sender_id TEXT NOT NULL CHECK (
        sender_id LIKE '________-____-____-____-____________'
    ),
    FOREIGN KEY (original_message_id) REFERENCES messages (message_id) ON DELETE CASCADE,
    FOREIGN KEY (conversation_id) REFERENCES conversations (conversation_id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users (user_id) ON DELETE SET NULL
);
