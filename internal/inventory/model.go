package inventory

import (
	"time"

	"github.com/google/uuid"
)

type InventoryItem struct {
	ID        uuid.UUID `json:"id"`
	PlayerID  uuid.UUID `json:"player_id"`
	ItemID    uuid.UUID `json:"item_id"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
