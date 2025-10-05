package auth

import (
	"context"
	"time"
)

type IUserRepository interface {
	IUserRepositoryReader
	IUserRepositoryWriter
}

type IUserRepositoryReader interface {
	GetUserByName(c context.Context, name string) (*User, error)
	GetUserByID(c context.Context, id int) (*User, error)
}

type IUserRepositoryWriter interface {
	CreateUser(c context.Context, user *User) error
}

type ITokenBlacklist interface {
	Add(ctx context.Context, token string, duration time.Duration) error
	Exist(ctx context.Context, token string) (bool, error)
}
