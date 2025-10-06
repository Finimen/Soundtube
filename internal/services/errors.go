package services

import "errors"

var (
	UserAlreadyExits = errors.New("user already exists")
	DBError          = errors.New("db error")
)
