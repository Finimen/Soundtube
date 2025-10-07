package services

import (
	"context"
	"errors"
	"soundtube/internal/domain/sound"
	"soundtube/pkg"
)

type SoundService struct {
	repository sound.ISoundRepository
	logger     *pkg.CustomLogger
}

func NewSoundService(repository sound.ISoundRepository, logger *pkg.CustomLogger) *SoundService {
	return &SoundService{repository: repository, logger: logger}
}

func (s *SoundService) CreateSound(ctx context.Context, name, album, genre string, authorID int) error {
	ctx, span := s.logger.GetTracer().Start(ctx, "SoundService.CreateSound")
	defer span.End()

	sound, err := sound.NewSound(name, album, genre, authorID)
	if err != nil {
		s.logger.Error("invalid sound params", err).WithTrace(ctx)
		return err
	}

	existsSound, err := s.repository.GetSoundByName(ctx, name)
	if err != nil {
		s.logger.Error("db error", err)
		return err
	}

	if existsSound != nil {
		err = errors.New("user already exits")
		s.logger.Error("invalid sound params", err).WithTrace(ctx)
		return err
	}

	err = s.repository.CreateSound(ctx, sound)
	if err != nil {
		s.logger.Error("db error", err).WithTrace(ctx)
		return err
	}

	return nil
}

func (s *SoundService) DeleteSound(ctx context.Context, name string) error {
	_, span := s.logger.GetTracer().Start(ctx, "SoundService.CreateSound")
	defer span.End()
	return nil
}

func (s *SoundService) GetSounds(ctx context.Context) ([]*sound.Sound, error) {
	ctx, span := s.logger.GetTracer().Start(ctx, "SoundService.CreateSound")
	defer span.End()

	sounds, err := s.repository.GetSounds(ctx)
	if err != nil {
		s.logger.Error("db error", err)
		return nil, err
	}

	return sounds, nil
}
