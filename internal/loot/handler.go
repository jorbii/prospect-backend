package loot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"prospect/internal/auth"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type createLootRequest struct {
	ItemID    uuid.UUID `json:"item_id"`
	Quantity  int64     `json:"quantity"`
	PositionX float32   `json:"position_x"`
	PositionY float32   `json:"position_y"`
}

func (h *Handler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var req createLootRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	loot, err := h.service.Create(
		r.Context(),
		Loot{
			ItemID:    req.ItemID,
			Quantity:  req.Quantity,
			PositionX: req.PositionX,
			PositionY: req.PositionY,
		},
	)
	if err != nil {
		if err == ErrInvalidItemID || err == ErrInvalidQuantity {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		http.Error(
			w,
			"failed to create loot",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(loot)
}

func (h *Handler) HandleLoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)

	case http.MethodPost:
		h.Create(w, r)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *Handler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	loots, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get loot",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(loots); err != nil {
		return
	}
}

func (h *Handler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/loot/",
	)

	id, err := uuid.Parse(idString)
	if err != nil {
		http.Error(
			w,
			"invalid loot id",
			http.StatusBadRequest,
		)
		return
	}

	loot, err := h.service.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			"loot not found",
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(loot); err != nil {
		return
	}
}

func (h *Handler) Pickup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lootIDString := strings.TrimPrefix(
		r.URL.Path,
		"/api/loot/pickup/",
	)

	lootID, err := uuid.Parse(lootIDString)
	if err != nil {
		http.Error(w, "invalid loot id", http.StatusBadRequest)
		return
	}

	userIDString, ok := auth.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		http.Error(w, "invalid player id", http.StatusUnauthorized)
		return
	}

	err = h.service.Pickup(
		r.Context(),
		userID,
		lootID,
	)
	if err != nil {
		fmt.Println("pickup error:", err)

		http.Error(
			w,
			"failed to pickup loot",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "loot picked up successfully",
	})
}
