package httpx

import (
	"net/http"

	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/httpx/handlers"
	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/httpx/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewServer(pool *pgxpool.Pool, jwtSecret []byte) http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	auth := handlers.NewAuthHandler(pool, jwtSecret)
	gate := handlers.NewGateHandler(pool)
	me := handlers.NewMeHandler(pool)
	quests := handlers.NewQuestsHandler(pool)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", auth.Register)
			r.Post("/login", auth.Login)
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTMiddleware(jwtSecret))

			r.Get("/me", me.Me)

			r.Route("/quests", func(r chi.Router) {
				r.Get("/", quests.List)
				r.Post("/", quests.Create)
				r.Patch("/{id}", quests.Patch)
				r.Delete("/{id}", quests.Delete)
			})
			r.Route("/gate", func(r chi.Router) {
				r.Post("/", gate.Open)
				r.Post("/{id}", gate.Close)
			})
		})
	})
	return r
}
