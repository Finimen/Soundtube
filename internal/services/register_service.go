package services

import (
	"context"
	"errors"
	"soundtube/internal/domain/auth"
	"soundtube/pkg"
)

type RegisterService struct {
	repository auth.IUserRepository
	logger     *pkg.CustomLogger
}

var (
	UserAlreadyExits = errors.New("user already exists")
)

func NewRegisterService(repository auth.IUserRepository, logger *pkg.CustomLogger) *RegisterService {
	return &RegisterService{repository: repository, logger: logger}
}

func (s *RegisterService) Register(с context.Context, username, email, password string) error {
	ctx, span := s.logger.GetTracer().Start(с, "RegisterService.Register")
	defer span.End()

	user, err := auth.NewUser(username, email, password)
	if err != nil {
		s.logger.Error("invalid user params", err).WithTrace(ctx)
		return err
	}

	existenceUser, err := s.repository.GetUserByName(ctx, user.Username())
	if err != nil {
		s.logger.Error("failed to check existence user", err).WithTrace(ctx)
		return err
	}

	if existenceUser != nil {
		err = UserAlreadyExits
		s.logger.Error("user already exists", err).WithTrace(ctx)
		return err
	}

	err = s.repository.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("db error", err).WithTrace(ctx)
		return err
	}

	s.logger.Info("user succesful registrated", user.Username())
	return nil
}
