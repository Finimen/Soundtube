package services

import (
	"context"
	"encoding/json"
	"fmt"
	"soundtube/internal/domain"
	"soundtube/internal/domain/reactions"
	"soundtube/internal/repositories"
	"soundtube/pkg"
	"sync"
	"time"
)

type ReactionService struct {
	repository   *repositories.SoundReactionRepository
	participants *repositories.SoundPartisipantsRepository
	logger       *pkg.CustomLogger
	cache        domain.ICache
}

type SoundReactionsResponse struct {
	SoundID      int     `json:"sound_id"`
	Likes        int     `json:"likes"`
	Dislikes     int     `json:"dislikes"`
	UserReaction *string `json:"user_reaction,omitempty"`
}

func NewRactionService(repository *repositories.SoundReactionRepository, participants *repositories.SoundPartisipantsRepository, cache domain.ICache, logger *pkg.CustomLogger) *ReactionService {
	return &ReactionService{
		repository:   repository,
		participants: participants,
		cache:        cache,
		logger:       logger,
	}
}

func (s *ReactionService) SetSoundReaction(ctx context.Context, userID, soundID int, reactionType string) error {
	ctx, span := s.logger.GetTracer().Start(ctx, "ReactionService.SetSoundReaction")
	defer span.End()

	existingReaction, err := s.participants.Get(ctx, userID, soundID)
	if err != nil {
		return err
	}

	if existingReaction != nil && existingReaction.ReactType == reactionType {
		err = s.repository.Delete(ctx, soundID, reactionType)
		if err != nil {
			return err
		}
		return s.participants.Remove(ctx, userID, soundID)
	}

	if existingReaction != nil {
		err = s.repository.Delete(ctx, soundID, existingReaction.ReactType)
		if err != nil {
			return err
		}
	}

	newReaction := reactions.NewReaction(soundID, reactionType)
	_, err = s.repository.Create(ctx, newReaction)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("sound_reactions:stats:%d", soundID)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Warn("failed to invalidate reaction cache", err).WithTrace(ctx)
	} else {
		s.logger.Info("reaction cache invalidated", "sound_id", soundID).WithTrace(ctx)
	}

	return s.participants.AddOrUpdate(ctx, userID, soundID, reactionType)
}

func (s *ReactionService) GetSoundReactions(ctx context.Context, userID, soundID int) (*SoundReactionsResponse, error) {
	ctx, span := s.logger.GetTracer().Start(ctx, "ReactionService.GetSoundReactions")
	defer span.End()

	cashedKey := fmt.Sprintf("sound_reactions:stats%d", soundID)

	var reactionStats *repositories.ReactionStatus
	if cashed, err := s.cache.Get(ctx, cashedKey); err == nil {
		if err := json.Unmarshal([]byte(cashed), &reactionStats); err == nil {
			s.logger.Info("reaction stats loaded from cache", "sound_id", soundID).WithTrace(ctx)
		}
	}

	if reactionStats == nil {
		var err error
		reactionStats, err = s.repository.GetReactionStats(ctx, soundID)
		if err != nil {
			return nil, err
		}

		if data, err := json.Marshal(reactionStats); err != nil {
			if err := s.cache.Set(ctx, cashedKey, data, 15*time.Minute); err != nil {
				s.logger.Warn("failed to cache reaction stats", err).WithTrace(ctx)
			}
		}
	}

	var userReaction *string
	if userID > 0 {
		userReactionObj, err := s.participants.Get(ctx, userID, soundID)
		if err != nil {
			return nil, err
		}
		if userReactionObj != nil {
			userReaction = &userReactionObj.ReactType
		}
	}

	return &SoundReactionsResponse{
		SoundID:      soundID,
		Likes:        reactionStats.Likes,
		Dislikes:     reactionStats.Dislikes,
		UserReaction: userReaction,
	}, nil
}

func (s *ReactionService) GetSoundsReactions(ctx context.Context, userID int, soundIDs []int) ([]SoundReactionsResponse, error) {
	ctx, span := s.logger.GetTracer().Start(ctx, "ReactionService.GetSoundsReactions")
	defer span.End()

	if len(soundIDs) == 0 {
		return []SoundReactionsResponse{}, nil
	}

	var reactionStats map[int]*repositories.ReactionStatus
	var userReactions map[int]string
	var err1, err2 error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		reactionStats, err1 = s.repository.GetReactionBatch(ctx, soundIDs)
	}()

	go func() {
		defer wg.Done()
		if userID > 0 {
			userReactions, err2 = s.participants.GetUserReactonBatch(ctx, userID, soundIDs)
		} else {
			userReactions = make(map[int]string)
		}
	}()

	wg.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	responses := make([]SoundReactionsResponse, 0, len(soundIDs))

	for _, soundID := range soundIDs {
		var likes, dislikes int
		var userReaction *string

		if stats, exists := reactionStats[soundID]; exists {
			likes = stats.Likes
			dislikes = stats.Dislikes
		}

		if reactType, exists := userReactions[soundID]; exists {
			userReaction = &reactType
		}

		response := SoundReactionsResponse{
			SoundID:      soundID,
			Likes:        likes,
			Dislikes:     dislikes,
			UserReaction: userReaction,
		}

		responses = append(responses, response)
	}

	return responses, nil
}
