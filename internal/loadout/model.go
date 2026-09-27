package loadout

import (
	"time"

	"github.com/google/uuid"
)

type Loadout struct {
	ID                uuid.UUID  `json:"id"`
	PlayerID          uuid.UUID  `json:"player_id"`
	PrimaryWeaponID   *uuid.UUID `json:"primary_weapon_id,omitempty"`
	SecondaryWeaponID *uuid.UUID `json:"secondary_weapon_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
