package item

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidName      = errors.New("item name is required")
	ErrInvalidType      = errors.New("item type is required")
	ErrInvalidMaxStack  = errors.New("max stack must be greater than zero")
	ErrInvalidStackData = errors.New("non-stackable item must have max stack equal to 1")
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
	item Item,
) (Item, error) {
	item.Name = strings.TrimSpace(item.Name)
	item.Type = strings.TrimSpace(item.Type)
	item.Rarity = strings.TrimSpace(item.Rarity)

	if item.Name == "" {
		return Item{}, ErrInvalidName
	}

	if item.Type == "" {
		return Item{}, ErrInvalidType
	}

	if item.MaxStack <= 0 {
		return Item{}, ErrInvalidMaxStack
	}

	if !item.Stackable && item.MaxStack != 1 {
		return Item{}, ErrInvalidStackData
	}

	if item.Rarity == "" {
		item.Rarity = "common"
	}

	return s.repository.Create(ctx, item)
}

func (s *Service) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (Item, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetAll(
	ctx context.Context,
) ([]Item, error) {
	return s.repository.GetAll(ctx)
}
