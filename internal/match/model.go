package match

import (
	"time"

	"github.com/google/uuid"
)

type matches struct {
	id         uuid.UUID
	status     string
	MaxPlayers int64
	StartedAt  *time.Time
	CreatedAt  time.Time
	EndedAt    *time.Time
}

type matchesPlayers struct {
	MatchId  uuid.UUID `json:"match_id"`
	PlayerId uuid.UUID `json:"player_id"`
	JoinedAt time.Time `json:"joined_at"`
}
