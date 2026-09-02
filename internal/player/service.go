package player

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreatePlayerTx(
	ctx context.Context,
	tx pgx.Tx,
	userID uuid.UUID,
) error {
	now := time.Now()

	player := Player{
		ID:              uuid.New(),
		UserID:          userID,
		XP:              0,
		Level:           1,
		Cash:            0,
		PremiumCurrency: 0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	return s.repository.CreateTx(ctx, tx, player)
}

func (s *Service) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (Player, error) {
	return s.repository.GetByUserID(ctx, userID)
}

func (s *Service) AddXP(
	ctx context.Context,
	playerID uuid.UUID,
	amount int64,
) (Player, error) {
	if amount <= 0 {
		return Player{}, errors.New("xp amount must be positive")
	}

	return s.repository.AddXP(
		ctx,
		playerID,
		amount,
	)
}

func calculateLevel(xp int64) int {
	switch {
	case xp >= 10000:
		return 10
	case xp >= 8000:
		return 9
	case xp >= 6500:
		return 8
	case xp >= 5000:
		return 7
	case xp >= 4000:
		return 6
	case xp >= 3000:
		return 5
	case xp >= 2000:
		return 4
	case xp >= 1000:
		return 3
	case xp >= 500:
		return 2
	default:
		return 1
	}
}

func (s *Service) AddCash(
	ctx context.Context,
	playerID uuid.UUID,
	amount int64,
) error {
	if amount <= 0 {
		return errors.New("cash amount must be positive")
	}

	return s.repository.AddCash(ctx, playerID, amount)
}

func (s *Service) RemoveCash(
	ctx context.Context,
	playerID uuid.UUID,
	amount int64,
) error {
	if amount <= 0 {
		return errors.New("cash amount must be positive")
	}

	removed, err := s.repository.RemoveCash(
		ctx,
		playerID,
		amount,
	)
	if err != nil {
		return err
	}

	if !removed {
		return errors.New("insufficient cash")
	}

	return nil
}

func (s *Service) AddPremiumCurrency(
	ctx context.Context,
	playerID uuid.UUID,
	amount int64,
) error {
	if amount <= 0 {
		return errors.New("premium currency amount must be positive")
	}

	return s.repository.AddPremiumCurrency(
		ctx,
		playerID,
		amount,
	)
}

func (s *Service) RemovePremiumCurrency(
	ctx context.Context,
	playerID uuid.UUID,
	amount int64,
) error {
	if amount <= 0 {
		return errors.New("premium currency amount must be positive")
	}

	removed, err := s.repository.RemovePremiumCurrency(
		ctx,
		playerID,
		amount,
	)
	if err != nil {
		return err
	}

	if !removed {
		return errors.New("insufficient premium currency")
	}

	return nil
}
