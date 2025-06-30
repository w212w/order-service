package handlers

import (
	"encoding/json"
	"net/http"
	e "order-service/internal/entity"
	"order-service/internal/repository"
	"order-service/internal/storage/cache"

	"github.com/gorilla/mux"
)

type Handler struct {
	cache *cache.Cache
	repo  repository.OrderRepository
}

func NewHandler(c *cache.Cache, r repository.OrderRepository) *Handler {
	return &Handler{cache: c, repo: r}
}

func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(e.ErrorResponse{Message: message})
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	orderUID := vars["order_uid"]

	if orderUID == "" {
		WriteJSONError(w, http.StatusBadRequest, "[error] empty order_uid")
		return
	}

	// Поиск значения в кэше
	if order, ok := h.cache.Get(orderUID); ok {
		resp := ConvertToResponse(order)
		json.NewEncoder(w).Encode(resp)
		return
	}

	order, err := h.repo.GetByUID(orderUID)
	if err != nil {
		WriteJSONError(w, http.StatusNotFound, "[error] order not found")
		return
	}

	h.cache.Set(order)
	resp := ConvertToResponse(order)
	json.NewEncoder(w).Encode(resp)
}
