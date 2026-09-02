package player

import (
	"time"

	"github.com/google/uuid"
)

type Player struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	XP              int64
	Level           int
	Cash            int64
	PremiumCurrency int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type PlayerResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	XP              int64     `json:"xp"`
	Level           int       `json:"level"`
	Cash            int64     `json:"cash"`
	PremiumCurrency int64     `json:"premium_currency"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
