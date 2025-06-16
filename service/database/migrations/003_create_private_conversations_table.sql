CREATE TABLE IF NOT EXISTS private_conversations (
    conversation_id TEXT PRIMARY KEY CHECK (
        conversation_id LIKE '________-____-____-____-____________'
    ),
    FOREIGN KEY (conversation_id) REFERENCES conversations (conversation_id) ON DELETE CASCADE
);
