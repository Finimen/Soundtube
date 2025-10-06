package repositories

import (
	"context"
	"database/sql"
	_ "embed"
	"soundtube/internal/domain/auth"
)

type UserRepository struct {
	db *sql.DB
}

//go:embed migrations/user/001_create_user_table_up.sql
var createUserTable string

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	var userRepository = UserRepository{}
	userRepository.db = db

	_, err := db.Exec(createUserTable)
	if err != nil {
		return nil, err
	}

	return &userRepository, nil
}

func (r *UserRepository) GetUserByName(ctx context.Context, name string) (*auth.User, error) {
	query := `SELECT id, user_password, user_email, is_verified, is_banned, verify_token
				FROM users WHERE user_name = $1`
	row := r.db.QueryRowContext(ctx, query, name)

	var id int
	var password, email, verifyToken string
	var isVerified, isBanned bool
	err := row.Scan(&id, &password, &email, &isVerified, &isBanned, &verifyToken)
	if err != nil {
		return nil, err
	}
	user := auth.RebuildUserFromStorage(id, name, email, password, isVerified, isBanned)
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*auth.User, error) {
	query := `SELECT user_name, user_password, user_email, is_verified, is_banned, verify_token
				FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var name, password, email, verifyToken string
	var isVerified, isBanned bool
	err := row.Scan(&id, &password, &email, &isVerified, &isBanned, &verifyToken)
	if err != nil {
		return nil, err
	}
	user := auth.RebuildUserFromStorage(id, name, email, password, isVerified, isBanned)
	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *auth.User) error {
	query := `INSERT INTO users (user_name, user_email, user_password, is_verified, is_banned, verify_token)`
	_, err := r.db.ExecContext(ctx, query, user.Username(), user.Email(), user.Password(), user.IsVerified(), user.IsBanned(), user.VerifyToken())
	return err
}
