package loot

import (
	"context"
	"errors"
	"prospect/internal/player"

	"github.com/google/uuid"
)

var (
	ErrInvalidItemID   = errors.New("item id is required")
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
	ErrInvalidLootID   = errors.New("loot id is required")
	ErrInvalidPlayerID = errors.New("player id is required")
)

type Service struct {
	repository       *Repository
	playerRepository *player.Repository
}

func NewService(
	repository *Repository,
	playerRepository *player.Repository,
) *Service {
	return &Service{
		repository:       repository,
		playerRepository: playerRepository,
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

func (s *Service) Pickup(
	ctx context.Context,
	userID uuid.UUID,
	lootID uuid.UUID,
) error {
	if userID == uuid.Nil {
		return ErrInvalidPlayerID
	}

	if lootID == uuid.Nil {
		return ErrInvalidLootID
	}

	player, err := s.playerRepository.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.repository.PickupTx(
		ctx,
		player.ID,
		lootID,
	)
}
