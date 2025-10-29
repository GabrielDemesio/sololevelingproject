package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/core/battle"
	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/httpx/middleware"
	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GateHandler struct {
	db *pgxpool.Pool
}

func NewGateHandler(db *pgxpool.Pool) *GateHandler {
	return &GateHandler{db: db}
}
func (h *GateHandler) Open(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var in struct {
		QuestID *uuid.UUID `json:"questid"`
		Rank    string     `json:"rank"`
		Minutes int        `json:"minutes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Minutes <= 0 {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if in.Rank == "" {
		in.Rank = "E"
	}
	run := store.FocusRun{
		ID:            uuid.New(),
		UserID:        uuid.MustParse(uid),
		QuestID:       in.QuestID,
		DungeonRank:   in.Rank,
		StartAt:       time.Now(),
		TargetMinutes: in.Minutes,
		XPEarned:      0,
		GoldEarned:    0,
	}
	if err := store.CreateFocusRun(r.Context(), h.db, &run); err != nil {
		http.Error(w, "failed to create run", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(run)
}
func (h *GateHandler) Close(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	runID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var in struct {
		Result  string  `json:"result"`
		Quality float64 `json:"quality"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	run, err := store.GetFocusRunByID(r.Context(), h.db, runID)
	if err != nil || run.UserID.String() != uid {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if run.EndAt != nil {
		http.Error(w, "already closed", http.StatusBadRequest)
		return
	}
	if in.Quality <= 0 {
		in.Quality = 1.0
	}
	mult := battle.RankMultiplier(run.DungeonRank)
	questWeight := 1
	if run.QuestID != nil {
		if w, err := store.GetQuestWeight(r.Context(), h.db, *run.QuestID); err == nil {
			questWeight = w
		}
	}

	xp, gold := battle.ComputeRewards(run.TargetMinutes, mult, questWeight, in.Quality)
	now := time.Now()
	run.EndAt = &now
	run.Result = &in.Result
	run.XPEarned = xp
	run.GoldEarned = gold
	if err := store.FinishFocusRun(r.Context(), h.db, &run); err != nil {
		http.Error(w, "failed to update run", http.StatusBadRequest)
		return
	}
	if in.Result == "success" {
		_ = store.AddXPAndGold(r.Context(), h.db, run.UserID, xp, gold, true)
	} else {
		_ = store.UpdateStreak(r.Context(), h.db, run.UserID, false)
	}
	json.NewEncoder(w).Encode(run)
}
