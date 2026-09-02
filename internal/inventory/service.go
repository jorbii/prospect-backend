package inventory

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidQuantity   = errors.New("quantity must be greater than zero")
	ErrInsufficientItems = errors.New("insufficient item quantity")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetInventory(
	ctx context.Context,
	playerID uuid.UUID,
) ([]InventoryItem, error) {
	return s.repository.GetByPlayerID(ctx, playerID)
}

func (s *Service) AddItem(
	ctx context.Context,
	playerID uuid.UUID,
	itemID uuid.UUID,
	quantity int64,
) (InventoryItem, error) {
	if quantity <= 0 {
		return InventoryItem{}, ErrInvalidQuantity
	}

	return s.repository.AddItem(
		ctx,
		playerID,
		itemID,
		quantity,
	)
}

func (s *Service) RemoveItem(
	ctx context.Context,
	playerID uuid.UUID,
	itemID uuid.UUID,
	quantity int64,
) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	err := s.repository.RemoveItem(
		ctx,
		playerID,
		itemID,
		quantity,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteItem(
	ctx context.Context,
	playerID uuid.UUID,
	itemID uuid.UUID,
) error {
	return s.repository.DeleteItem(
		ctx,
		playerID,
		itemID,
	)
}
