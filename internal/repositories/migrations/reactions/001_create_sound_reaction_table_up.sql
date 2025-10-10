CREATE TABLE IF NOT EXISTS sound_reactions(
    id SERIAL PRIMARY KEY,
    sound_id INTEGER UNIQUE REFERENCES sounds(id) ON DELETE CASCADE,
    total_likes INTEGER DEFAULT 0,
    total_dislikes INTEGER DEFAULT 0
);