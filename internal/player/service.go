package player

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"prospect/internal/auth"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetMe(w http.ResponseWriter, r *http.Request) {
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

	player, err := s.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, contextCanceledError()) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}

		http.Error(w, "player not found", http.StatusNotFound)
		return
	}

	response := MeResponse{
		ID:       player.ID,
		Username: player.Username,
		Email:    player.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func contextCanceledError() error {
	return context.Canceled
}
