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
	ctx, span := r.logger.GetTracer().Start(ctx, "UserRepository.GetSounds")
	defer span.End()

	res, err := r.db.QueryContext(ctx, `SELECT 
		id, author_id, sound_name, sound_album, sound_genre, duration, file_name, file_size, file_format, upload_date
		FROM sounds`)
	if err != nil {
		return nil, err
	}

	sounds := make([]*sound.Sound, 16)
	for res.Next() {
		var id, authorId, duration, fileSize int
		var soundName, soundAlbum, soundGenre, fileName, fileFormat, uploadDate string
		err = res.Scan(&id, &authorId, &soundName, &soundAlbum, &soundGenre, &duration, &fileName, &fileSize, &fileFormat, &uploadDate)
		if err != nil {
			return nil, err
		}
		sound := sound.RebuildSoundFromStorage(id, authorId, duration, soundName, soundAlbum, soundGenre, fileName, fileSize, fileFormat, uploadDate)
		sounds = append(sounds, sound)
	}

	return nil, nil
}

func (r *SoundRepository) GetSoundByName(ctx context.Context, name string) (*sound.Sound, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "UserRepository.GetUserByName")
	defer span.End()

	return nil, nil
}

func (r *SoundRepository) GetSoundByID(ctx context.Context, id int) (*sound.Sound, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "UserRepository.GetSoundByID")
	defer span.End()

	return nil, nil
}

func (r *SoundRepository) CreateSound(ctx context.Context, sound *sound.Sound) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "UserRepository.CreateSound")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	query := `INSERT INTO sounds (id, author_id, sound_name, sound_album, sound_genre, duration, file_name, file_size, file_format)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = tx.ExecContext(ctx, query, sound.ID(), sound.AuthorID(), sound.Name(), sound.Ablum(), sound.Genre(), sound.Duration(), sound.FileName(), sound.FileSize(), sound.FileFormat())
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *SoundRepository) DeleteSound(ctx context.Context, sound *sound.Sound) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "UserRepository.DeleteSound")
	defer span.End()

	return nil
}
