package weapon

import (
	"encoding/json"
	"net/http"
	"strings"

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

func (h *Handler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	weapons, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get weapons",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(weapons); err != nil {
		return
	}
}

func (h *Handler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		"/api/weapons/",
	)

	if strings.HasPrefix(path, "item/") {
		h.GetByItemID(w, r)
		return
	}

	id, err := uuid.Parse(path)
	if err != nil {
		http.Error(
			w,
			"invalid weapon id",
			http.StatusBadRequest,
		)
		return
	}

	weapon, err := h.service.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			"weapon not found",
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(weapon); err != nil {
		return
	}
}

func (h *Handler) GetByItemID(
	w http.ResponseWriter,
	r *http.Request,
) {
	itemIDString := strings.TrimPrefix(
		r.URL.Path,
		"/api/weapons/item/",
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

	weapon, err := h.service.GetByItemID(
		r.Context(),
		itemID,
	)
	if err != nil {
		http.Error(
			w,
			"weapon not found",
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(weapon); err != nil {
		return
	}
}
