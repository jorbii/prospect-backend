package loot

import (
	"time"

	"github.com/google/uuid"
)

type Loot struct {
	ID        uuid.UUID `json:"id"`
	ItemID    uuid.UUID `json:"item_id"`
	Quantity  int64     `json:"quantity"`
	PositionX float32   `json:"position_x"`
	PositionY float32   `json:"position_y"`
	CreatedAt time.Time `json:"created_at"`
}
