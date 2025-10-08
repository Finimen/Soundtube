CREATE TABLE IF NOT EXISTS sound_reactions(
    sound_id INTEGER REFERENCES sounds(id) ON DELETE CASCADE,
    total_likes INTEGER DEFAULT 0,
    total_dislikes INTEGER DEFAULT 0
)