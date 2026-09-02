package player

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, player Player) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO players (
			id,
			user_id,
			xp,
			level,
			cash,
			premium_currency,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		player.ID,
		player.UserID,
		player.XP,
		player.Level,
		player.Cash,
		player.PremiumCurrency,
		player.CreatedAt,
		player.UpdatedAt,
	)

	return err
}

func (r *Repository) CreateTx(
	ctx context.Context,
	tx pgx.Tx,
	player Player,
) error {
	_, err := tx.Exec(
		ctx,
		`INSERT INTO players (
			id,
			user_id,
			xp,
			level,
			cash,
			premium_currency,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		player.ID,
		player.UserID,
		player.XP,
		player.Level,
		player.Cash,
		player.PremiumCurrency,
		player.CreatedAt,
		player.UpdatedAt,
	)

	return err
}

func (r *Repository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (Player, error) {
	var player Player

	err := r.db.QueryRow(
		ctx,
		`SELECT
			id,
			user_id,
			xp,
			level,
			cash,
			premium_currency,
			created_at,
			updated_at
		FROM players
		WHERE user_id = $1`,
		userID,
	).Scan(
		&player.ID,
		&player.UserID,
		&player.XP,
		&player.Level,
		&player.Cash,
		&player.PremiumCurrency,
		&player.CreatedAt,
		&player.UpdatedAt,
	)

	if err != nil {
		return Player{}, err
	}

	return player, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	playerID uuid.UUID,
) (Player, error) {
	var player Player

	err := r.db.QueryRow(
		ctx,
		`SELECT
			id,
			user_id,
			xp,
			level,
			cash,
			premium_currency,
			created_at,
			updated_at
		FROM players
		WHERE id = $1`,
		playerID,
	).Scan(
		&player.ID,
		&player.UserID,
		&player.XP,
		&player.Level,
		&player.Cash,
		&player.PremiumCurrency,
		&player.CreatedAt,
		&player.UpdatedAt,
	)

	if err != nil {
		return Player{}, err
	}

	return player, nil
}

func (r *Repository) UpdateProgression(
	ctx context.Context,
	player Player,
) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE players
		SET
			xp = $1,
			level = $2,
			updated_at = $3
		WHERE id = $4`,
		player.XP,
		player.Level,
		player.UpdatedAt,
		player.ID,
	)

	return err
}

func (r *Repository) AddCash(
	ctx context.Context,
	playerID uuid.UUID,
	amount int64,
) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE players
		SET
			cash = cash + $1,
			updated_at = NOW()
		WHERE id = $2`,
		amount,
		playerID,
	)

	return err
}

func (r *Repository) RemoveCash(
	ctx context.Context,
	playerID uuid.UUID,
	amount int64,
) (bool, error) {
	result, err := r.db.Exec(
		ctx,
		`UPDATE players
		SET
			cash = cash - $1,
			updated_at = NOW()
		WHERE id = $2
		  AND cash >= $1`,
		amount,
		playerID,
	)

	if err != nil {
		return false, err
	}

	return result.RowsAffected() == 1, nil
}

func (r *Repository) AddPremiumCurrency(
	ctx context.Context,
	playerID uuid.UUID,
	amount int64,
) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE players
		SET
			premium_currency = premium_currency + $1,
			updated_at = NOW()
		WHERE id = $2`,
		amount,
		playerID,
	)

	return err
}

func (r *Repository) RemovePremiumCurrency(
	ctx context.Context,
	playerID uuid.UUID,
	amount int64,
) (bool, error) {
	result, err := r.db.Exec(
		ctx,
		`UPDATE players
		SET
			premium_currency = premium_currency - $1,
			updated_at = NOW()
		WHERE id = $2
		  AND premium_currency >= $1`,
		amount,
		playerID,
	)

	if err != nil {
		return false, err
	}

	return result.RowsAffected() == 1, nil
}

func (r *Repository) AddXP(
	ctx context.Context,
	playerID uuid.UUID,
	amount int64,
) (Player, error) {
	var player Player

	err := r.db.QueryRow(
		ctx,
		`UPDATE players
		SET
			xp = xp + $1,
			level = CASE
				WHEN xp + $1 >= 10000 THEN 10
				WHEN xp + $1 >= 8000 THEN 9
				WHEN xp + $1 >= 6500 THEN 8
				WHEN xp + $1 >= 5000 THEN 7
				WHEN xp + $1 >= 4000 THEN 6
				WHEN xp + $1 >= 3000 THEN 5
				WHEN xp + $1 >= 2000 THEN 4
				WHEN xp + $1 >= 1000 THEN 3
				WHEN xp + $1 >= 500 THEN 2
				ELSE 1
			END,
			updated_at = NOW()
		WHERE id = $2
		RETURNING
			id,
			user_id,
			xp,
			level,
			cash,
			premium_currency,
			created_at,
			updated_at`,
		amount,
		playerID,
	).Scan(
		&player.ID,
		&player.UserID,
		&player.XP,
		&player.Level,
		&player.Cash,
		&player.PremiumCurrency,
		&player.CreatedAt,
		&player.UpdatedAt,
	)

	if err != nil {
		return Player{}, err
	}

	return player, nil
}
