package repositories

import (
	"context"
	"database/sql"
	"soundtube/internal/domain/sound"
	"soundtube/pkg"

	_ "embed"
)

type SoundRepository struct {
	db     *sql.DB
	logger *pkg.CustomLogger
}

//go:embed migrations/sound/001_create_sound_table_up.sql
var createSoundTable string

func NewSoundRepository(db *sql.DB, logger *pkg.CustomLogger) (*SoundRepository, error) {
	soundRepository := SoundRepository{
		db:     db,
		logger: logger,
	}

	_, err := soundRepository.db.Exec(createSoundTable)
	if err != nil {
		return nil, err
	}

	return &soundRepository, nil
}

func (r *SoundRepository) GetSounds(ctx context.Context) ([]*sound.Sound, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundRepository.GetSounds")
	defer span.End()

	rows, err := r.db.QueryContext(ctx, `SELECT 
		id, author_id, sound_name, sound_album, sound_genre, duration, file_name, file_size, file_format, upload_date, file_path
		FROM sounds`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sounds []*sound.Sound
	for rows.Next() {
		var id, authorId, duration, fileSize int
		var soundName, soundAlbum, soundGenre, fileName, fileFormat, uploadDate, filePath string
		err = rows.Scan(&id, &authorId, &soundName, &soundAlbum, &soundGenre, &duration, &fileName, &fileSize, &fileFormat, &uploadDate, &filePath)
		if err != nil {
			return nil, err
		}
		sound := sound.RebuildSoundFromStorage(id, authorId, duration, soundName, soundAlbum, soundGenre, fileName, filePath, fileSize, fileFormat, uploadDate)
		sounds = append(sounds, sound)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sounds, nil
}

func (r *SoundRepository) GetSoundByName(ctx context.Context, name string) (*sound.Sound, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundRepository.GetSoundByName")
	defer span.End()

	query := `SELECT id, author_id, sound_name, sound_album, sound_genre, duration, file_name, file_size, file_format, upload_date, file_path
        FROM sounds WHERE sound_name = $1`

	var id, authorID, duration, fileSize int
	var soundName, soundAlbum, soundGenre, fileName, fileFormat, uploadDate, filePath string

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&id, &authorID, &soundName, &soundAlbum, &soundGenre,
		&duration, &fileName, &fileSize, &fileFormat, &uploadDate, &filePath,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	sound := sound.RebuildSoundFromStorage(id, authorID, duration, soundName, soundAlbum, soundGenre, fileName, filePath, fileSize, fileFormat, uploadDate)
	return sound, nil
}

func (r *SoundRepository) GetSoundByID(ctx context.Context, id int) (*sound.Sound, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundRepository.GetSoundByID")
	defer span.End()

	query := `SELECT id, author_id, sound_name, sound_album, sound_genre, duration, file_name, file_size, file_format, upload_date, file_path
        FROM sounds WHERE id = $1`

	var authorID, duration, fileSize int
	var soundName, soundAlbum, soundGenre, fileName, fileFormat, uploadDate, filePath string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&id, &authorID, &soundName, &soundAlbum, &soundGenre,
		&duration, &fileName, &fileSize, &fileFormat, &uploadDate, &filePath,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	sound := sound.RebuildSoundFromStorage(id, authorID, duration, soundName, soundAlbum, soundGenre, fileName, filePath, fileSize, fileFormat, uploadDate)
	return sound, nil
}

func (r *SoundRepository) CreateSound(ctx context.Context, sound *sound.Sound) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundRepository.CreateSound")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	query := `INSERT INTO sounds (author_id, sound_name, sound_album, sound_genre, duration, file_name, file_size, file_format)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.ExecContext(ctx, query, sound.AuthorID(), sound.Name(), sound.Ablum(), sound.Genre(), sound.Duration(), sound.FileName(), sound.FileSize(), sound.FileFormat())
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *SoundRepository) DeleteSound(ctx context.Context, sound *sound.Sound) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundRepository.DeleteSound")
	defer span.End()

	return nil
}

func (r *SoundRepository) UpdateSoundFile(ctx context.Context, name, fileName, filePath string, fileSize int64) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundRepository.UploadSoundFile")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE sounds SET file_name = $1, file_path = $2, file_size = $3 WHERE sound_name = $4`

	_, err = tx.ExecContext(ctx, query, fileName, filePath, fileSize, name)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
