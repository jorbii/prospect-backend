package item

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
	items, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get items",
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

func (h *Handler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/items/",
	)

	id, err := uuid.Parse(idString)
	if err != nil {
		http.Error(
			w,
			"invalid item id",
			http.StatusBadRequest,
		)
		return
	}

	item, err := h.service.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			"item not found",
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(item); err != nil {
		return
	}
}
