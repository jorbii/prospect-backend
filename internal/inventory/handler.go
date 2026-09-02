package inventory

import (
	"encoding/json"
	"net/http"
	"strings"

	"prospect/internal/auth"
	"prospect/internal/player"

	"github.com/google/uuid"
)

type Handler struct {
	service       *Service
	playerService *player.Service
}

func NewHandler(
	service *Service,
	playerService *player.Service,
) *Handler {
	return &Handler{
		service:       service,
		playerService: playerService,
	}
}

type addItemRequest struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity int64     `json:"quantity"`
}

func (h *Handler) getPlayerID(
	r *http.Request,
) (uuid.UUID, bool) {
	userID, ok := auth.UserIDFromContext(r)
	if !ok {
		return uuid.Nil, false
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, false
	}

	player, err := h.playerService.GetByUserID(
		r.Context(),
		userUUID,
	)
	if err != nil {
		return uuid.Nil, false
	}

	return player.ID, true
}

func (h *Handler) GetInventory(
	w http.ResponseWriter,
	r *http.Request,
) {
	playerID, ok := h.getPlayerID(r)
	if !ok {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	items, err := h.service.GetInventory(
		r.Context(),
		playerID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to get inventory",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(items); err != nil {
		return
	}
}

func (h *Handler) AddItem(
	w http.ResponseWriter,
	r *http.Request,
) {
	playerID, ok := h.getPlayerID(r)
	if !ok {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	var req addItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	item, err := h.service.AddItem(
		r.Context(),
		playerID,
		req.ItemID,
		req.Quantity,
	)
	if err != nil {
		if err == ErrInvalidQuantity {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		http.Error(
			w,
			"failed to add item",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(item); err != nil {
		return
	}
}

func (h *Handler) DeleteItem(
	w http.ResponseWriter,
	r *http.Request,
) {
	playerID, ok := h.getPlayerID(r)
	if !ok {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	itemIDString := strings.TrimPrefix(
		r.URL.Path,
		"/api/inventory/items/",
	)

	itemID, err := uuid.Parse(itemIDString)
	if err != nil {
		http.Error(
			w,
			"invalid item id",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.DeleteItem(
		r.Context(),
		playerID,
		itemID,
	); err != nil {
		http.Error(
			w,
			"failed to delete item",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
