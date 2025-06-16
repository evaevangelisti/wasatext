CREATE TABLE IF NOT EXISTS members (
    conversation_id TEXT CHECK (
        conversation_id LIKE '________-____-____-____-____________'
    ),
    user_id TEXT CHECK (
        user_id LIKE '________-____-____-____-____________'
    ),
    joined_at TEXT NOT NULL CHECK (joined_at LIKE "____-__-__T__:__:__Z"),
    PRIMARY KEY (conversation_id, user_id),
    FOREIGN KEY (conversation_id) REFERENCES group_conversations (conversation_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_members_conversation_id ON members (conversation_id);
