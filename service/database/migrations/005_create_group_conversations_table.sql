CREATE TABLE IF NOT EXISTS group_conversations (
    conversation_id TEXT PRIMARY KEY CHECK (
        conversation_id LIKE '________-____-____-____-____________'
    ),
    name TEXT NOT NULL CHECK (
        LENGTH (name) >= 1
        AND LENGTH (name) <= 50
    ),
    photo TEXT CHECK (
        LENGTH (photo) >= 11
        AND LENGTH (photo) <= 255
    ),
    FOREIGN KEY (conversation_id) REFERENCES conversations (conversation_id) ON DELETE CASCADE
);
