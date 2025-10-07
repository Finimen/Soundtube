package services

import (
	"context"
	"errors"
	"soundtube/internal/domain/auth"
	"soundtube/internal/domain/sound"
	"soundtube/pkg"
)

type SoundService struct {
	repository sound.ISoundRepository
	logger     *pkg.CustomLogger
	user       auth.IUserRepositoryReader
}

func NewSoundService(repository sound.ISoundRepository, user auth.IUserRepositoryReader, logger *pkg.CustomLogger) *SoundService {
	return &SoundService{repository: repository, logger: logger, user: user}
}

func (s *SoundService) CreateSound(ctx context.Context, name, album, genre string, authorID int) error {
	ctx, span := s.logger.GetTracer().Start(ctx, "SoundService.CreateSound")
	defer span.End()

	userExists, err := s.user.UserExists(ctx, authorID)
	if err != nil {
		s.logger.Error("db error checking user", err).WithTrace(ctx)
		return err
	}

	if !userExists {
		err = errors.New("user does not exist")
		s.logger.Error("invalid user id", err).WithTrace(ctx)
		return err
	}

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
	_, span := s.logger.GetTracer().Start(ctx, "SoundService.DeleteSound")
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

func (s *SoundService) UpdateSoundFile(ctx context.Context, name, filename, filepath string, fileSize int64) error {
	ctx, span := s.logger.GetTracer().Start(ctx, "SoundService.UploadSoundFile")
	defer span.End()

	err := s.repository.UpdateSoundFile(ctx, name, filename, filepath, fileSize)
	if err != nil {
		s.logger.Error("failed to update sound file info in repository", err).WithTrace(ctx)
		return err
	}

	return nil
}
