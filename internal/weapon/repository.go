package weapon

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
	weapon Weapon,
) (Weapon, error) {
	var createdWeapon Weapon

	err := r.db.QueryRow(ctx, `
		INSERT INTO weapons (
			item_id,
			damage,
			fire_rate,
			magazine_size,
			reload_time,
			range
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			item_id,
			damage,
			fire_rate,
			magazine_size,
			reload_time,
			range,
			created_at
	`,
		weapon.ItemID,
		weapon.Damage,
		weapon.FireRate,
		weapon.MagazineSize,
		weapon.ReloadTime,
		weapon.Range,
	).Scan(
		&createdWeapon.ID,
		&createdWeapon.ItemID,
		&createdWeapon.Damage,
		&createdWeapon.FireRate,
		&createdWeapon.MagazineSize,
		&createdWeapon.ReloadTime,
		&createdWeapon.Range,
		&createdWeapon.CreatedAt,
	)

	if err != nil {
		return Weapon{}, err
	}

	return createdWeapon, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (Weapon, error) {
	var weapon Weapon

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			item_id,
			damage,
			fire_rate,
			magazine_size,
			reload_time,
			range,
			created_at
		FROM weapons
		WHERE id = $1
	`, id).Scan(
		&weapon.ID,
		&weapon.ItemID,
		&weapon.Damage,
		&weapon.FireRate,
		&weapon.MagazineSize,
		&weapon.ReloadTime,
		&weapon.Range,
		&weapon.CreatedAt,
	)

	if err != nil {
		return Weapon{}, err
	}

	return weapon, nil
}

func (r *Repository) GetByItemID(
	ctx context.Context,
	itemID uuid.UUID,
) (Weapon, error) {
	var weapon Weapon

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			item_id,
			damage,
			fire_rate,
			magazine_size,
			reload_time,
			range,
			created_at
		FROM weapons
		WHERE item_id = $1
	`, itemID).Scan(
		&weapon.ID,
		&weapon.ItemID,
		&weapon.Damage,
		&weapon.FireRate,
		&weapon.MagazineSize,
		&weapon.ReloadTime,
		&weapon.Range,
		&weapon.CreatedAt,
	)

	if err != nil {
		return Weapon{}, err
	}

	return weapon, nil
}

func (r *Repository) GetAll(
	ctx context.Context,
) ([]Weapon, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			item_id,
			damage,
			fire_rate,
			magazine_size,
			reload_time,
			range,
			created_at
		FROM weapons
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	weapons := make([]Weapon, 0)

	for rows.Next() {
		var weapon Weapon

		err := rows.Scan(
			&weapon.ID,
			&weapon.ItemID,
			&weapon.Damage,
			&weapon.FireRate,
			&weapon.MagazineSize,
			&weapon.ReloadTime,
			&weapon.Range,
			&weapon.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		weapons = append(weapons, weapon)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return weapons, nil
}
