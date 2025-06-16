CREATE TABLE IF NOT EXISTS participants (
    conversation_id TEXT CHECK (
        conversation_id LIKE '________-____-____-____-____________'
    ),
    user_id TEXT CHECK (
        user_id LIKE '________-____-____-____-____________'
    ),
    PRIMARY KEY (conversation_id, user_id),
    FOREIGN KEY (conversation_id) REFERENCES private_conversations (conversation_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_participants_conversation_id ON participants (conversation_id);
