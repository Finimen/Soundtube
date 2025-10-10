package repositories

import (
	"context"
	"database/sql"
	_ "embed"
	"soundtube/pkg"

	"github.com/lib/pq"
)

type SoundPartisipantsRepository struct {
	db     *sql.DB
	logger *pkg.CustomLogger
}

type SoundPartisipantsResponse struct {
	SoundID   int
	UserID    int
	ReactID   int
	ReactType string
}

//go:embed migrations/reactions/001_create_soupd_partisipants_table_up.sql
var createPartisipantsTable string

func NewSoundPartisipantsRepository(db *sql.DB, logger *pkg.CustomLogger) (*SoundPartisipantsRepository, error) {
	repository := SoundPartisipantsRepository{db: db, logger: logger}
	_, err := repository.db.Exec(createPartisipantsTable)
	if err != nil {
		return nil, err
	}

	return &repository, nil
}

func (r *SoundPartisipantsRepository) Get(ctx context.Context, userID, soundID int) (*SoundPartisipantsResponse, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundPartisipantsRepository.Get")
	defer span.End()

	query := "SELECT react_type FROM sound_participants WHERE user_id = $1 AND sound_id = $2"
	var reactType string
	err := r.db.QueryRowContext(ctx, query, userID, soundID).Scan(&reactType)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	response := SoundPartisipantsResponse{
		SoundID:   soundID,
		UserID:    userID,
		ReactType: reactType,
	}
	return &response, nil
}

func (r *SoundPartisipantsRepository) AddOrUpdate(ctx context.Context, userID, soundID int, reactType string) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundPartisipantsRepository.AddOrUpdate")
	defer span.End()

	query := `INSERT INTO sound_participants (sound_id, user_id, react_type) 
              VALUES ($1, $2, $3)
              ON CONFLICT (sound_id, user_id) 
              DO UPDATE SET react_type = $3`
	_, err := r.db.ExecContext(ctx, query, soundID, userID, reactType)
	return err
}

func (r *SoundPartisipantsRepository) Remove(ctx context.Context, userID, soundID int) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundPartisipantsRepository.SoundPartisipantsRepository")
	defer span.End()

	query := "DELETE FROM sound_participants WHERE user_id = $1 AND sound_id = $2"
	_, err := r.db.ExecContext(ctx, query, userID, soundID)
	return err
}

func (r *SoundPartisipantsRepository) GetUserReactonBatch(ctx context.Context, userID int, soundIDs []int) (map[int]string, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundPartisipantsRepository.GetUserReactonBatch")
	defer span.End()

	if len(soundIDs) == 0 {
		return map[int]string{}, nil
	}

	query := `SELECT sound_id, react_type
	FROM sound_participants
	WHERE user_id = $1 AND sound_id = ANY($2)`

	rows, err := r.db.QueryContext(ctx, query, userID, pq.Array(soundIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userReactions := make(map[int]string)
	for rows.Next() {
		var soundID int
		var reactType string
		if err := rows.Scan(&soundID, &reactType); err != nil {
			return nil, err
		}
		userReactions[soundID] = reactType
	}

	return userReactions, nil
}
