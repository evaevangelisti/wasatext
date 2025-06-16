CREATE TABLE IF NOT EXISTS conversations (
    conversation_id TEXT PRIMARY KEY CHECK (
        conversation_id LIKE '________-____-____-____-____________'
    ),
    type TEXT NOT NULL CHECK (type IN ("private", "group")),
    created_at TEXT NOT NULL CHECK (created_at LIKE "____-__-__T__:__:__Z")
);
