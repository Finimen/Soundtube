CREATE TABLE IN NO EXISTS comments(
    id SERIAL PRIMARY KEY,
    author_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    is_response BOOLEAN DEFAULT FALSE,
    response_target INTEGER DEFAULT -1,
);

CREATE INDEX IF NOT EXISTS idx_comments_id ON comments(id);
CREATE INDEX IF NOT EXISTS idx_comments_author_id ON comments(author_id); 