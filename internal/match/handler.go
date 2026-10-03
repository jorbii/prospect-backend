package match

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

type createMatchResponse struct {
	ID uuid.UUID `json:"id"`
}

func (h *Handler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	matchID, err := h.service.CreateMatch(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to create match",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createMatchResponse{
		ID: matchID,
	}); err != nil {
		return
	}
}
