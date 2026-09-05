package weapon

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidItemID       = errors.New("item id is required")
	ErrInvalidDamage       = errors.New("damage must be greater than zero")
	ErrInvalidFireRate     = errors.New("fire rate must be greater than zero")
	ErrInvalidMagazineSize = errors.New("magazine size must be greater than zero")
	ErrInvalidReloadTime   = errors.New("reload time must be greater than zero")
	ErrInvalidRange        = errors.New("range must be greater than zero")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	weapon Weapon,
) (Weapon, error) {
	if weapon.ItemID == uuid.Nil {
		return Weapon{}, ErrInvalidItemID
	}

	if weapon.Damage <= 0 {
		return Weapon{}, ErrInvalidDamage
	}

	if weapon.FireRate <= 0 {
		return Weapon{}, ErrInvalidFireRate
	}

	if weapon.MagazineSize <= 0 {
		return Weapon{}, ErrInvalidMagazineSize
	}

	if weapon.ReloadTime <= 0 {
		return Weapon{}, ErrInvalidReloadTime
	}

	if weapon.Range <= 0 {
		return Weapon{}, ErrInvalidRange
	}

	return s.repository.Create(ctx, weapon)
}

func (s *Service) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (Weapon, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetByItemID(
	ctx context.Context,
	itemID uuid.UUID,
) (Weapon, error) {
	return s.repository.GetByItemID(ctx, itemID)
}

func (s *Service) GetAll(
	ctx context.Context,
) ([]Weapon, error) {
	return s.repository.GetAll(ctx)
}
