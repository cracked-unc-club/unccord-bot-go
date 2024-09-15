
-- PostgreSQL DDL for the unccord-bot-go application

-- Create the starboard table
CREATE TABLE starboard (
    id SERIAL PRIMARY KEY,               -- Auto-incrementing ID for each record
    message_id TEXT NOT NULL UNIQUE,     -- ID of the starred message (Discord message ID)
    channel_id TEXT NOT NULL,            -- ID of the channel where the message was posted
    author_id TEXT NOT NULL,             -- ID of the message author
    content TEXT NOT NULL,               -- Content of the starred message
    star_count INT NOT NULL DEFAULT 1,   -- Number of stars (reactions) the message has
    starboard_message_id TEXT,           -- ID of the message posted to the starboard
    posted_at TIMESTAMP DEFAULT NOW()    -- Timestamp when the message was starred
);

-- Index for quick lookup by message_id
CREATE INDEX idx_message_id ON starboard (message_id);

-- Index for quick lookup by starboard_message_id
CREATE INDEX idx_starboard_message_id ON starboard (starboard_message_id);
