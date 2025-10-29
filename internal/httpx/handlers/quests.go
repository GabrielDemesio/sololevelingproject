package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/httpx/middleware"
	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestsHandler struct {
	db *pgxpool.Pool
}

func NewQuestsHandler(db *pgxpool.Pool) *QuestsHandler {
	return &QuestsHandler{db: db}
}
func (h *QuestsHandler) List(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	items, err := store.ListQuestsByUser(r.Context(), h.db, uuid.MustParse(uid))
	if err != nil {
		http.Error(w, "failed to list quests", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func (h *QuestsHandler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var in struct {
		Title       string   `json:"title"`
		Description *string  `json:"description"`
		Weight      int      `json:"weight"`
		DueAt       *string  `json:"dueAt"`
		Tags        []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Title == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	var due *time.Time
	if in.DueAt != nil && *in.DueAt != "" {
		t, err := time.Parse(time.RFC3339, *in.DueAt)
		if err == nil {
			due = &t
		}
	}
	weight := in.Weight
	if weight <= 0 {
		weight = 1
	}
	q := store.Quest{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(uid),
		Title:       in.Title,
		Description: in.Description,
		Weight:      weight,
		Status:      "open",
		DueAt:       due,
		Tags:        in.Tags,
		CreatedAt:   time.Now(),
	}
	if err := store.CreateQuest(r.Context(), h.db, &q); err != nil {
		http.Error(w, "failed to create quest", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}
func (h *QuestsHandler) Patch(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	qid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var in struct {
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		Weight      *int     `json:"weight"`
		Status      *string  `json:"status"`
		Tags        []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	q, err := store.GetQuestByID(r.Context(), h.db, qid)
	if err != nil || q.UserID.String() != uid {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if in.Title != nil && *in.Title != "" {
		q.Title = *in.Title
	}
	if in.Description != nil {
		q.Description = in.Description
	}
	if in.Weight != nil && *in.Weight > 0 {
		q.Weight = *in.Weight
	}
	if in.Status != nil && *in.Status != "" {
		q.Status = *in.Status
	}
	if in.Tags != nil && len(in.Tags) > 0 {
		q.Tags = in.Tags
	}
	if err := store.UpdateQuest(r.Context(), h.db, q); err != nil {
		http.Error(w, "failed to update quest", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(q)
}

func (h *QuestsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	qid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	q, err := store.GetQuestByID(r.Context(), h.db, qid)
	if err != nil || q.UserID.String() != uid {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err := store.DeleteQuest(r.Context(), h.db, qid); err != nil {
		http.Error(w, "failed to delete quest", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
