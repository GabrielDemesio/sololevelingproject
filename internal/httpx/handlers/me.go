package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/httpx/middleware"
	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MeHandler struct {
	db *pgxpool.Pool
}

func NewMeHandler(db *pgxpool.Pool) *MeHandler {
	return &MeHandler{db: db}
}

func (h *MeHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	u, err := store.GetUserByID(r.Context(), h.db, uuid.MustParse(uid))
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(u)
}
