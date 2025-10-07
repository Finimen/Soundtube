package repositories

import (
	"context"
	"database/sql"
	_ "embed"
	"soundtube/internal/domain/auth"
	"soundtube/pkg"
)

type UserRepository struct {
	db     *sql.DB
	logger *pkg.CustomLogger
}

//go:embed migrations/user/001_create_user_table_up.sql
var createUserTable string

func NewUserRepository(db *sql.DB, logger *pkg.CustomLogger) (*UserRepository, error) {
	var userRepository = UserRepository{
		db:     db,
		logger: logger,
	}

	_, err := db.Exec(createUserTable)
	if err != nil {
		return nil, err
	}

	return &userRepository, nil
}

func (r *UserRepository) GetUserByName(ctx context.Context, name string) (*auth.User, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "UserRepository.GetUserByName")
	defer span.End()

	query := `SELECT id, user_password, user_email, is_verified, is_banned, verify_token
				FROM users WHERE user_name = $1`
	row := r.db.QueryRowContext(ctx, query, name)

	var id int
	var password, email, verifyToken string
	var isVerified, isBanned bool
	err := row.Scan(&id, &password, &email, &isVerified, &isBanned, &verifyToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	user := auth.RebuildUserFromStorage(id, name, email, password, isVerified, isBanned)
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*auth.User, error) {
	ctx, span := r.logger.GetTracer().Start(ctx, "UserRepository.GetUserByName")
	defer span.End()

	query := `SELECT user_name, user_password, user_email, is_verified, is_banned, verify_token
				FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var name, password, email, verifyToken string
	var isVerified, isBanned bool
	err := row.Scan(&id, &password, &email, &isVerified, &isBanned, &verifyToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	user := auth.RebuildUserFromStorage(id, name, email, password, isVerified, isBanned)
	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *auth.User) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "UserRepository.GetUserByName")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	query := `INSERT INTO users (user_name, user_email, user_password, is_verified, is_banned, verify_token)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, query, user.Username(), user.Email(), user.Password(), user.IsVerified(), user.IsBanned(), user.VerifyToken())
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "UserRepository.DeleteUser")
	defer span.End()

	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("delete user failed", err).WithTrace(ctx)
		return err
	}

	return nil
}

func (r *UserRepository) MarkUserAsVerified(ctx context.Context, id int) error {
	ctx, span := r.logger.GetTracer().Start(ctx, "UserRepository.MarkUserAsVerified")
	defer span.End()

	query := "UPDATE users SET is_verified = true WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Verify update failed", err).WithTrace(ctx)
		return err
	}

	return nil
}
