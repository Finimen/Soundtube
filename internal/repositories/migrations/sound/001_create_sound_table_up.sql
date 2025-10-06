CREATE TABLE IF NOT EXISTS sounds(
    id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    sound_name VARCHAR(255) NOT NULL,
    sound_album VARCHAR(255),
    sound_genre VARCHAR(100),
    duration INTEGER,
    file_name VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    file_format VARCHAR(10) NOT NULL,
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active'
)

CREATE INDEX idx_sounds_author_id ON sounds(author_id)
CREATE INDEX idx_sounds_genre ON sounds(sound_genre) 