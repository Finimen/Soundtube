package services

import (
	"context"
	"errors"
	"soundtube/internal/domain/auth"
	"soundtube/pkg"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	repository auth.IUserRepository
	logger     *pkg.CustomLogger
	jwtkey     []byte
}

func NewLoginService(jwtKey []byte, repository auth.IUserRepository, logger *pkg.CustomLogger) *LoginService {
	return &LoginService{jwtkey: jwtKey, repository: repository, logger: logger}
}

func (s *LoginService) Login(ctx context.Context, username, password string) (string, error) {
	ctx, span := s.logger.GetTracer().Start(ctx, "LoginService.Login")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.name", username),
	)

	if username == "" || password == "" {
		err := errors.New("username & password are requered")
		s.logger.Warn(err.Error(), err).WithTrace(ctx)
		return "", err
	}

	user, err := s.repository.GetUserByName(ctx, username)
	if err != nil || user == nil {
		s.logger.Warn("user not found", err).WithTrace(ctx)
		return "", err
	}

	if !user.IsVerified() {
		err = errors.New("user not verified")
		s.logger.Warn(err.Error(), err).WithTrace(ctx)
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password())); err != nil {
		s.logger.Warn(err.Error(), err).WithTrace(ctx)
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtkey)
	if err != nil {
		s.logger.Error("token generation error", err).WithTrace(ctx)
		return "", err
	}

	s.logger.Info("login successful", "username", username).WithTrace(ctx)
	return tokenString, nil
}
