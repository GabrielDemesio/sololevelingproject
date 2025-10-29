package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db        *pgxpool.Pool
	jwtSecret []byte
}

func NewAuthHandler(db *pgxpool.Pool, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{db: db, jwtSecret: jwtSecret}
}
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || in.Password == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	id := uuid.New()
	now := time.Now()
	u := store.User{
		ID:        id,
		Email:     in.Email,
		PassHash:  string(hash),
		Level:     1,
		XP:        0,
		Gold:      0,
		Stats:     map[string]int{"focus": 1, "discipline": 1, "energy": 1, "creativity": 1},
		Streak:    0,
		CreatedAt: now,
	}
	if err := store.CreateUser(r.Context(), h.db, &u); err != nil {
		http.Error(w, "failed to create user", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": u.ID})
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	u, err := store.GetUserByEmail(r.Context(), h.db, in.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PassHash), []byte(in.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": "solo-leveling",
		"sub": 	u.ID.String(),
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
		"rnd": randomNonce(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(h.jwtSecret)
	json.NewEncoder(w).Encode(map[string]any{
		"accesToken"	: ss,
	})
}
func randomNonce() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "0"
	}
	return hex.EncodeToString(b)
}