package repositories

import "database/sql"

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	var userRepository = UserRepository{}
	userRepository.db = db
	return &userRepository, nil
}
