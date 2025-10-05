package services

import (
	"context"
	"soundtube/internal/domain/auth"
	"soundtube/pkg"

	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	repository   auth.IUserRepository
	emailService auth.IEmailSener
	logger       *pkg.CustomLogger
}

func NewRegisterService(repository auth.IUserRepository, email auth.IEmailSener, logger *pkg.CustomLogger) *RegisterService {
	return &RegisterService{repository: repository, emailService: email, logger: logger}
}

func (s *RegisterService) Register(с context.Context, username, email, password string) error {
	ctx, span := s.logger.GetTracer().Start(с, "RegisterService.Register")
	defer span.End()

	existenceUser, err := s.repository.GetUserByName(ctx, username)
	if err != nil {
		s.logger.Error("failed to check existence user", err).WithTrace(ctx)
		return err
	}

	if existenceUser != nil {
		err = UserAlreadyExits
		s.logger.Warn("user already exists", err).WithTrace(ctx)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("hashing password failed", err).WithTrace(ctx)
		return err
	}

	verifyToken, err := generateVerifyToken()
	if err != nil {
		s.logger.Error("failed to generate verification token", err)
		return err
	}

	user, err := auth.NewUser(username, email, string(hashedPassword), verifyToken)
	if err != nil {
		s.logger.Error("invalid user params", err).WithTrace(ctx)
		return err
	}

	err = s.repository.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("db error", err).WithTrace(ctx)
		return err
	}

	if err := s.emailService.SendVerificationEmail(ctx, email, verifyToken); err != nil {
		s.logger.Error("verification email failed", err).WithTrace(ctx)
		return err
	}

	s.logger.Info("user succesful registrated", user.Username()).WithTrace(ctx)
	return nil
}
