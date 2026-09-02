package player

import (
	"encoding/json"
	"errors"
	"net/http"

	"prospect/internal/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r)
	if !ok || userID == "" {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	player, err := h.service.GetByUserID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "player not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := PlayerResponse{
		ID:              player.ID,
		UserID:          player.UserID,
		XP:              player.XP,
		Level:           player.Level,
		Cash:            player.Cash,
		PremiumCurrency: player.PremiumCurrency,
		CreatedAt:       player.CreatedAt,
		UpdatedAt:       player.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
