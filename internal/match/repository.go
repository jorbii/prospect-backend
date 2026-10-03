package match

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateMatch(ctx context.Context, matches matches) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO matches (
			id,
			status,
			max_players,
			started_at,
			created_at,
			ended_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		matches.id,
		matches.status,
		matches.MaxPlayers,
		matches.StartedAt,
		matches.CreatedAt,
		matches.EndedAt,
	)

	return err
}
