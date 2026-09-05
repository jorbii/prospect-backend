package item

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Stackable bool      `json:"stackable"`
	MaxStack  int       `json:"max_stack"`
	Rarity    string    `json:"rarity"`
	CreatedAt time.Time `json:"created_at"`
}
