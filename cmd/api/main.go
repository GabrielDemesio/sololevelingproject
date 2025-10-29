package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/httpx"
	"github.com/gabrieldemesio/solo-leveling-go-mvp-v2/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	addr      = ":" + getenv("PORT", "8080")
	dbDsn     = getenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/solo?sslmode=disable")
	jwtSecret = getenv("JWT_SECRET", "secret")
)

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	doMigrate := flag.Bool("migrate", false, "run migrations and exit")
	flag.Parse()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbDsn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("dn Ping: ", err)
	}
	if *doMigrate {
		if err := store.RunMigrations(pool); err != nil {
			log.Println("auto-migrate warning: ", err)
		}
	}

	router := httpx.NewServer(pool, []byte(jwtSecret))

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("server listening on", addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("Listen: ", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Shutdown: ", err)
	}
	fmt.Println("Done!")
}
