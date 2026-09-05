package loot

import (
	"context"

	"github.com/google/uuid"
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

func (r *Repository) Create(
	ctx context.Context,
	loot Loot,
) (Loot, error) {
	var createdLoot Loot

	err := r.db.QueryRow(ctx, `
		INSERT INTO loot (
			item_id,
			quantity,
			position_x,
			position_y
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			item_id,
			quantity,
			position_x,
			position_y,
			created_at
	`,
		loot.ItemID,
		loot.Quantity,
		loot.PositionX,
		loot.PositionY,
	).Scan(
		&createdLoot.ID,
		&createdLoot.ItemID,
		&createdLoot.Quantity,
		&createdLoot.PositionX,
		&createdLoot.PositionY,
		&createdLoot.CreatedAt,
	)

	if err != nil {
		return Loot{}, err
	}

	return createdLoot, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (Loot, error) {
	var loot Loot

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			item_id,
			quantity,
			position_x,
			position_y,
			created_at
		FROM loot
		WHERE id = $1
	`, id).Scan(
		&loot.ID,
		&loot.ItemID,
		&loot.Quantity,
		&loot.PositionX,
		&loot.PositionY,
		&loot.CreatedAt,
	)

	if err != nil {
		return Loot{}, err
	}

	return loot, nil
}

func (r *Repository) GetAll(
	ctx context.Context,
) ([]Loot, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			item_id,
			quantity,
			position_x,
			position_y,
			created_at
		FROM loot
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loots := make([]Loot, 0)

	for rows.Next() {
		var loot Loot

		err := rows.Scan(
			&loot.ID,
			&loot.ItemID,
			&loot.Quantity,
			&loot.PositionX,
			&loot.PositionY,
			&loot.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		loots = append(loots, loot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return loots, nil
}

func (r *Repository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM loot
		WHERE id = $1
	`, id)

	return err
}
