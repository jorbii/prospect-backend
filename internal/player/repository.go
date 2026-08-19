package player

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(id uuid.UUID) (Player, error) {
	var player Player

	err := r.db.QueryRow(
		context.Background(),
		`SELECT id, username, email
		 FROM users
		 WHERE id = $1`,
		id,
	).Scan(
		&player.ID,
		&player.Username,
		&player.Email,
	)

	if err != nil {
		return Player{}, err
	}

	return player, nil
}
