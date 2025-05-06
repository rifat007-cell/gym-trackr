CREATE TABLE chat_messages (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    role TEXT NOT NULL, 
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
