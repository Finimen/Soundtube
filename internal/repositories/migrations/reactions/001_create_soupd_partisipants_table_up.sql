CREATE TABLE IF NOT EXISTS sound_participants( 
    sound_id INTEGER REFERENCES sounds(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    react_type TEXT NOT NULL CHECK (react_type IN ('like', 'dislike')),
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (sound_id, user_id) 
);

CREATE INDEX IF NOT EXISTS idx_sound_participants_user_id ON sound_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_sound_participants_sound_id ON sound_participants(sound_id);