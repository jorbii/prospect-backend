package match

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MatchCreator interface {
	CreateMatch(ctx context.Context) (uuid.UUID, error)
}

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateMatch(ctx context.Context) (uuid.UUID, error) {
	match := matches{
		id:         uuid.New(),
		status:     "waiting",
		MaxPlayers: 10,
		CreatedAt:  time.Now(),
	}

	err := s.repository.CreateMatch(ctx, match)
	if err != nil {
		return uuid.Nil, err
	}

	return match.id, nil
}

/*
func (s *Service) JoinMatch(
	ctx context.Context,
	matchID uuid.UUID,
	userID uuid.UUID,
) error {

	match, err := s.repository.GetMatch(ctx, matchID)
	if err != nil {
		return err
	}

	if match.Status != "waiting" {
		return ErrMatchNotWaiting
	}

	playersCount, err := s.repository.CountPlayers(ctx, matchID)
	if err != nil {
		return err
	}

	if playersCount >= match.MaxPlayers {
		return ErrMatchFull
	}

	exists, err := s.repository.IsPlayerInMatch(ctx, matchID, userID)
	if err != nil {
		return err
	}

	if exists {
		return ErrAlreadyInMatch
	}

	return s.repository.AddPlayer(ctx, matchID, userID)
}
*/
