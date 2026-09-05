package item

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
	item Item,
) (Item, error) {
	var createdItem Item

	err := r.db.QueryRow(ctx, `
		INSERT INTO items (
			name,
			type,
			stackable,
			max_stack,
			rarity
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			name,
			type,
			stackable,
			max_stack,
			rarity,
			created_at
	`,
		item.Name,
		item.Type,
		item.Stackable,
		item.MaxStack,
		item.Rarity,
	).Scan(
		&createdItem.ID,
		&createdItem.Name,
		&createdItem.Type,
		&createdItem.Stackable,
		&createdItem.MaxStack,
		&createdItem.Rarity,
		&createdItem.CreatedAt,
	)

	if err != nil {
		return Item{}, err
	}

	return createdItem, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (Item, error) {
	var item Item

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			name,
			type,
			stackable,
			max_stack,
			rarity,
			created_at
		FROM items
		WHERE id = $1
	`, id).Scan(
		&item.ID,
		&item.Name,
		&item.Type,
		&item.Stackable,
		&item.MaxStack,
		&item.Rarity,
		&item.CreatedAt,
	)

	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func (r *Repository) GetAll(
	ctx context.Context,
) ([]Item, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			name,
			type,
			stackable,
			max_stack,
			rarity,
			created_at
		FROM items
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)

	for rows.Next() {
		var item Item

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Type,
			&item.Stackable,
			&item.MaxStack,
			&item.Rarity,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
