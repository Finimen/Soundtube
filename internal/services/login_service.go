package services

import (
	"context"
	"errors"
	"fmt"
	"soundtube/internal/domain/auth"
	"soundtube/pkg"
	"soundtube/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	repository auth.IUserRepository
	blackList  auth.ITokenBlacklist
	logger     *pkg.CustomLogger
	jwtkey     []byte
	exp        int
}

func NewLoginService(cfg config.Token, repository auth.IUserRepository, blackList auth.ITokenBlacklist, logger *pkg.CustomLogger) *LoginService {
	return &LoginService{jwtkey: []byte(cfg.JwtKey), exp: cfg.Exp, repository: repository, blackList: blackList, logger: logger}
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

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password()), []byte(password)); err != nil {
		s.logger.Warn("invalid password", err).WithTrace(ctx)
		return "", err
	}

	now := time.Now()

	expiration := time.Hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID(),
		"username": username,
		"exp":      now.Add(expiration).Unix(),
		"iat":      now.Unix(),
	})

	tokenString, err := token.SignedString(s.jwtkey)
	if err != nil {
		s.logger.Error("token generation error", err).WithTrace(ctx)
		return "", err
	}

	s.logger.Info("login successful", "username", username).WithTrace(ctx)
	return tokenString, nil
}

func (s *LoginService) Logout(ctx context.Context, token string) error {
	ctx, span := s.logger.GetTracer().Start(ctx, "LoginService.Logout")
	defer span.End()

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return s.jwtkey, nil
	})
	if err != nil {
		return err
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		return errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid token expectation")
	}

	expectation := time.Until(time.Unix(int64(exp), 0))
	if expectation < 0 {
		expectation = time.Minute * 5
	}

	if err = s.blackList.Add(ctx, token, expectation); err != nil {
		s.logger.Error("falied to add token to black list", err).WithTrace(ctx)
		return err
	}

	s.logger.Info("token added to black list successfully").WithTrace(ctx)
	return nil
}

func (s *LoginService) ValidToken(ctx context.Context, token string) (string, int, error) {
	ctx, span := s.logger.GetTracer().Start(ctx, "LoginService.ValidateToken")
	defer span.End()

	inBlacklist, err := s.blackList.Exist(ctx, token)
	if err != nil {
		s.logger.Error("blacklist check failed", err)
		return "", 0, err
	}
	if inBlacklist {
		return "", 0, errors.New("token is revoked")
	}

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return s.jwtkey, nil
	})
	if err != nil || !parsed.Valid {
		return "", 0, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, errors.New("invalid token claims")
	}

	username, _ := claims["username"].(string)

	var userID int
	switch sub := claims["sub"].(type) {
	case float64:
		userID = int(sub)
	case int:
		userID = sub
	case int64:
		userID = int(sub)
	default:
		return "", 0, errors.New("invalid user id type in token")
	}

	s.logger.Info("Token validation",
		"username", username,
		"userID", userID,
		"sub_type", fmt.Sprintf("%T", claims["sub"])).WithTrace(ctx)

	return username, userID, nil
}
