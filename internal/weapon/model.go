package weapon

import (
	"time"

	"github.com/google/uuid"
)

type Weapon struct {
	ID           uuid.UUID `json:"id"`
	ItemID       uuid.UUID `json:"item_id"`
	Damage       int       `json:"damage"`
	FireRate     int       `json:"fire_rate"`
	MagazineSize int       `json:"magazine_size"`
	ReloadTime   float64   `json:"reload_time"`
	Range        int       `json:"range"`
	CreatedAt    time.Time `json:"created_at"`
}
