package inventory

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

func (r *Repository) GetByPlayerID(ctx context.Context, playerID uuid.UUID) ([]InventoryItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, player_id, item_id, quantity, created_at, updated_at
		FROM inventory_items
		WHERE player_id = $1
		ORDER BY created_at ASC
	`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]InventoryItem, 0)

	for rows.Next() {
		var item InventoryItem

		err := rows.Scan(
			&item.ID,
			&item.PlayerID,
			&item.ItemID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
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

func (r *Repository) AddItem(
	ctx context.Context,
	playerID uuid.UUID,
	itemID uuid.UUID,
	quantity int64,
) (InventoryItem, error) {
	var item InventoryItem

	err := r.db.QueryRow(ctx, `
		INSERT INTO inventory_items (
			player_id,
			item_id,
			quantity
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (player_id, item_id)
		DO UPDATE SET
			quantity = inventory_items.quantity + EXCLUDED.quantity,
			updated_at = NOW()
		RETURNING
			id,
			player_id,
			item_id,
			quantity,
			created_at,
			updated_at
	`, playerID, itemID, quantity).Scan(
		&item.ID,
		&item.PlayerID,
		&item.ItemID,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return InventoryItem{}, err
	}

	return item, nil
}

func (r *Repository) RemoveItem(
	ctx context.Context,
	playerID uuid.UUID,
	itemID uuid.UUID,
	quantity int64,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE inventory_items
		SET
			quantity = quantity - $3,
			updated_at = NOW()
		WHERE player_id = $1
		  AND item_id = $2
		  AND quantity >= $3
	`, playerID, itemID, quantity)

	return err
}

func (r *Repository) DeleteItem(
	ctx context.Context,
	playerID uuid.UUID,
	itemID uuid.UUID,
) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM inventory_items
		WHERE player_id = $1
		  AND item_id = $2
	`, playerID, itemID)

	return err
}
