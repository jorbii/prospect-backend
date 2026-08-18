package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(user User) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO users (
			id,
			username,
			email,
			password_hash,
			created_at
		) VALUES ($1, $2, $3, $4, $5)`,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				switch pgErr.ConstraintName {
				case "users_username_key":
					return ErrUsernameExists

				case "users_email_key":
					return ErrEmailExists
				}
			}
		}

		return err
	}

	return nil
}

func (r *Repository) GetUserByEmail(email string) (User, error) {
	var user User

	err := r.db.QueryRow(
		context.Background(),
		`SELECT
			id,
			username,
			email,
			password_hash,
			created_at,
			last_login_at
		FROM users
		WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.LastLoginAt,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
