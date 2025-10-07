package auth

import (
	"context"
	"time"
)

type IUserRepository interface {
	IUserRepositoryReader
	IUserRepositoryWriter
	MarkUserAsVerified(ctx context.Context, id int) error
}

type IUserRepositoryReader interface {
	GetUserByName(ctx context.Context, name string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByToken(ctx context.Context, token string) (*User, error)
	UserExists(ctx context.Context, userID int) (bool, error)
}

type IUserRepositoryWriter interface {
	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) error
}

type ITokenBlacklist interface {
	Add(ctx context.Context, token string, duration time.Duration) error
	Exist(ctx context.Context, token string) (bool, error)
}

type IEmailSener interface {
	SendVerificationEmail(ctx context.Context, email, verifyToken string) error
}
