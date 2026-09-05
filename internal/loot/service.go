package loot

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidItemID   = errors.New("item id is required")
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
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
	loot Loot,
) (Loot, error) {
	if loot.ItemID == uuid.Nil {
		return Loot{}, ErrInvalidItemID
	}

	if loot.Quantity <= 0 {
		return Loot{}, ErrInvalidQuantity
	}

	return s.repository.Create(ctx, loot)
}

func (s *Service) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (Loot, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetAll(
	ctx context.Context,
) ([]Loot, error) {
	return s.repository.GetAll(ctx)
}

func (s *Service) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repository.Delete(ctx, id)
}
