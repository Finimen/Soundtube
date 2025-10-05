package auth

import (
	"errors"
	"soundtube/scripts"
)

type User struct {
	id          int
	username    string
	email       string
	password    string
	isVerified  bool
	isBanned    bool
	verifyToken string
}

func (u *User) ID() int          { return u.id }
func (u *User) Username() string { return u.username }
func (u *User) Email() string    { return u.email }
func (u *User) IsVerified() bool { return u.isVerified }
func (u *User) IsBanned() bool   { return u.isBanned }

func (u *User) VerifyToken() string { return u.verifyToken }

func (u *User) Password() string { return u.password }

func NewUser(username, email, password string) (*User, error) {
	if username == "" || scripts.ValidateXSS(username) {
		return nil, errors.New("username cannot be empty")
	}
	if password == "" || scripts.ValidateXSS(password) {
		return nil, errors.New("password cannot be empty")
	}
	if !scripts.ValidateEmail(email) {
		return nil, errors.New("invalid email")
	}

	return &User{
		username:   username,
		email:      email,
		password:   password,
		isVerified: false,
		isBanned:   false,
	}, nil
}

func RebuildUserFromStorage(id int, username, email, password string, isVerified, isBanned bool) *User {
	return &User{
		id:         id,
		username:   username,
		email:      email,
		password:   password,
		isVerified: isVerified,
		isBanned:   isBanned,
	}
}

func (u *User) VerifyEmail() {
	u.isVerified = true
}

func (u *User) Ban() {
	u.isBanned = true
}
