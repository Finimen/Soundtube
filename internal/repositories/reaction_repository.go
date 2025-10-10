package repositories

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"soundtube/internal/domain/reactions"
	"soundtube/pkg"

	"github.com/lib/pq"
)

type SoundReactionRepository struct {
	logger *pkg.CustomLogger
	db     *sql.DB
}

type ReactionStatus struct {
	Likes    int
	Dislikes int
}

//go:embed migrations/reactions/001_create_sound_reaction_table_up.sql
var createReactionTable string

func NewReactionRepository(db *sql.DB, logger *pkg.CustomLogger) (*SoundReactionRepository, error) {
	repository := SoundReactionRepository{db: db, logger: logger}

	_, err := db.Exec(createReactionTable)

	if err != nil {
		return nil, err
	}

	return &repository, nil
}

func (r *SoundReactionRepository) Delete(ctx context.Context, soundID int, reactType string) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundReactionRepository.Delete")
	defer span.End()

	_, err := r.db.ExecContext(ctx, `
        UPDATE sound_reactions 
        SET total_likes = total_likes - CASE WHEN $2 = 'like' THEN 1 ELSE 0 END,
            total_dislikes = total_dislikes - CASE WHEN $2 = 'dislike' THEN 1 ELSE 0 END
        WHERE sound_id = $1 AND (
            (total_likes > 0 AND $2 = 'like') OR 
            (total_dislikes > 0 AND $2 = 'dislike')
        )`,
		soundID, reactType)

	return err
}

func (r *SoundReactionRepository) Create(ctx context.Context, react *reactions.Reaction) (int, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundReactionRepository.Create")
	defer span.End()

	var id int
	err := r.db.QueryRowContext(ctx, `
        INSERT INTO sound_reactions (sound_id, total_likes, total_dislikes) 
        VALUES ($1, 
            CASE WHEN $2 = 'like' THEN 1 ELSE 0 END,
            CASE WHEN $2 = 'dislike' THEN 1 ELSE 0 END
        )
        ON CONFLICT (sound_id) DO UPDATE SET
            total_likes = CASE 
                WHEN $2 = 'like' THEN sound_reactions.total_likes + 1 
                ELSE sound_reactions.total_likes 
            END,
            total_dislikes = CASE 
                WHEN $2 = 'dislike' THEN sound_reactions.total_dislikes + 1 
                ELSE sound_reactions.total_dislikes 
            END
        RETURNING id`,
		react.GetTargetID(), react.GetType()).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to increment reaction: %w", err)
	}

	return id, nil
}

func (r *SoundReactionRepository) GetReactionBatch(ctx context.Context, soundIDs []int) (map[int]*ReactionStatus, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundReactionRepository.GetReactionBatch")
	defer span.End()

	if len(soundIDs) == 0 {
		return map[int]*ReactionStatus{}, nil
	}

	query := `SELECT sound_id, total_likes, total_dislikes
		FROM sound_reactions
		WHERE sound_id = ANY($1)`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(soundIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[int]*ReactionStatus)
	for rows.Next() {
		var soundID, likes, dislikes int
		if err := rows.Scan(&soundID, &likes, &dislikes); err != nil {
			return nil, err
		}
		stats[soundID] = &ReactionStatus{
			Likes:    likes,
			Dislikes: dislikes,
		}
	}

	for _, id := range soundIDs {
		if _, exists := stats[id]; !exists {
			stats[id] = &ReactionStatus{Likes: 0, Dislikes: 0}
		}
	}

	return stats, nil
}

func (r *SoundReactionRepository) GetReactionStats(ctx context.Context, soundID int) (*ReactionStatus, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundReactionRepository.GetReactionStats")
	defer span.End()

	var likes, dislikes int

	err := r.db.QueryRowContext(ctx, `
		SELECT total_likes, total_dislikes 
		FROM sound_reactions 
		WHERE sound_id = $1`,
		soundID).Scan(&likes, &dislikes)

	if err == sql.ErrNoRows {
		return &ReactionStatus{Likes: 0, Dislikes: 0}, nil
	}
	if err != nil {
		return nil, err
	}

	return &ReactionStatus{
		Likes:    likes,
		Dislikes: dislikes,
	}, nil
}

func (r *SoundReactionRepository) GetCountByType(ctx context.Context, soundID int, reactionType string) (int, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "SoundReactionRepository.GetReactionStats")
	defer span.End()

	var count int
	var column string

	if reactionType == "like" {
		column = "total_likes"
	} else if reactionType == "dislike" {
		column = "total_dislikes"
	} else {
		return 0, fmt.Errorf("invalid reaction type: %s", reactionType)
	}

	query := fmt.Sprintf("SELECT %s FROM sound_reactions WHERE sound_id = $1", column)
	err := r.db.QueryRowContext(ctx, query, soundID).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}
