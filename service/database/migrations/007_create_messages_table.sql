CREATE TABLE IF NOT EXISTS messages (
    message_id TEXT PRIMARY KEY CHECK (
        message_id LIKE '________-____-____-____-____________'
    ),
    content TEXT CHECK (
        LENGTH (content) >= 1
        AND LENGTH (content) <= 1000
    ),
    attachment TEXT CHECK (
        LENGTH (attachment) >= 11
        AND LENGTH (attachment) <= 255
    ),
    sent_at TEXT NOT NULL CHECK (
        sent_at LIKE "____-__-__T__:__:__Z" OR
        sent_at LIKE "____-__-__T__:__:__+__:__" OR
        sent_at LIKE "____-__-__T__:__:__-__:__"
    ),
    edited_at TEXT CHECK (
        edited_at LIKE "____-__-__T__:__:__Z" OR
        edited_at LIKE "____-__-__T__:__:__+__:__" OR
        edited_at LIKE "____-__-__T__:__:__-__:__"
    ),
    conversation_id TEXT NOT NULL CHECK (
        conversation_id LIKE '________-____-____-____-____________'
    ),
    sender_id TEXT NOT NULL CHECK (
        sender_id LIKE '________-____-____-____-____________'
    ),
    FOREIGN KEY (conversation_id) REFERENCES conversations (conversation_id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users (user_id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages (conversation_id);

CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages (sender_id);
