package store

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrations embed.FS

func RunMigrations(pool *pgxpool.Pool) error {
	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return err
	}

	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (filename TEXT PRIMARY KEY, ran_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
		return err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".sql") {
			files = append(files, name)
		}
	}
	sort.Strings(files)

	for _, name := range files {
		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename=$1)`, name).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		b, err := migrations.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(b)); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations(filename) VALUES($1)`, name); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
