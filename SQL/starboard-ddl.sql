CEATE TABLE starboard (
    id SERIAL PRIMARY KEY,
    message_id TEXT NOT NULL,
    channel_id TEXT NOT NULL.
    author_id TEXT NOT NULL,
    content TEXT NOT NULL,
    star_count INT NOT NULL DEFAULT 1,
    posted_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (message_id)
)