package auth

import (
	"context"
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
