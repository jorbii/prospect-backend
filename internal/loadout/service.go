package loadout

import (
	"context"
	"errors"

	"prospect/internal/player"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserID = errors.New("user id is required")
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
	userID uuid.UUID,
) (*Loadout, error) {

	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	player, err := s.playerRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repository.Create(ctx, player.ID)
}

func (s *Service) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*Loadout, error) {

	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	player, err := s.playerRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repository.GetByPlayerID(ctx, player.ID)
}
