package player

import "github.com/google/uuid"

type Player struct {
	ID       uuid.UUID
	Username string
	Email    string
}

type MeResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}
