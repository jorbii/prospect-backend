package loadout

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
	playerID uuid.UUID,
) (*Loadout, error) {

	loadout := &Loadout{}

	err := r.db.QueryRow(ctx, `
		INSERT INTO loadouts (
			player_id
		)
		VALUES ($1)
		RETURNING
			id,
			player_id,
			primary_weapon_id,
			secondary_weapon_id,
			created_at,
			updated_at
	`, playerID).Scan(
		&loadout.ID,
		&loadout.PlayerID,
		&loadout.PrimaryWeaponID,
		&loadout.SecondaryWeaponID,
		&loadout.CreatedAt,
		&loadout.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return loadout, nil
}

func (r *Repository) GetByPlayerID(
	ctx context.Context,
	playerID uuid.UUID,
) (*Loadout, error) {

	loadout := &Loadout{}

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			player_id,
			primary_weapon_id,
			secondary_weapon_id,
			created_at,
			updated_at
		FROM loadouts
		WHERE player_id = $1
	`, playerID).Scan(
		&loadout.ID,
		&loadout.PlayerID,
		&loadout.PrimaryWeaponID,
		&loadout.SecondaryWeaponID,
		&loadout.CreatedAt,
		&loadout.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return loadout, nil
}
