package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(ctx context.Context, db *pgxpool.Pool, u *User) error {
	_, err := db.Exec(ctx, `INSERT INTO users(id, email, pass_hash, level, xp, gold, stats, streak, last_active_date, created_at)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		u.ID, u.Email, u.PassHash, u.Level, u.XP, u.Gold, u.Stats, u.Streak, u.LastActiveDate, u.CreatedAt)
	return err
}

func GetUserByEmail(ctx context.Context, db *pgxpool.Pool, email string) (*User, error) {
	row := db.QueryRow(ctx, `SELECT id, email, pass_hash, level, xp, gold, stats, streak, last_active_date, created_at FROM users WHERE email=$1`, email)
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.PassHash, &u.Level, &u.XP, &u.Gold, &u.Stats, &u.Streak, &u.LastActiveDate, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByID(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) (*User, error) {
	row := db.QueryRow(ctx, `SELECT id, email, pass_hash, level, xp, gold, stats, streak, last_active_date, created_at FROM users WHERE id=$1`, id)
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.PassHash, &u.Level, &u.XP, &u.Gold, &u.Stats, &u.Streak, &u.LastActiveDate, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

// AddXPAndGold atualiza xp, gold e streak respeitando o dia em America/Sao_Paulo.
// success=true significa que o usuário concluiu o gate com sucesso.
func AddXPAndGold(ctx context.Context, db *pgxpool.Pool, userID uuid.UUID, xp, gold int64, success bool) error {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	today := time.Now().In(loc).Truncate(24 * time.Hour)

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var oldStreak int
	var lastDate *time.Time
	if err := tx.QueryRow(ctx, `SELECT streak, last_active_date FROM users WHERE id=$1 FOR UPDATE`, userID).Scan(&oldStreak, &lastDate); err != nil {
		return err
	}

	newStreak := oldStreak
	if success {
		if lastDate == nil {
			newStreak = 1
		} else {
			ld := lastDate.In(loc).Truncate(24 * time.Hour)
			diff := int(today.Sub(ld).Hours() / 24)
			switch diff {
			case 0:
				// mesmo dia, mantém
				newStreak = oldStreak
			case 1:
				// dia seguinte, +1
				newStreak = oldStreak + 1
			default:
				// passou mais de 1 dia, reseta pra 1
				newStreak = 1
			}
		}
	} else {
		// abandon => quebra streak
		newStreak = 0
	}

	if _, err := tx.Exec(ctx, `UPDATE users SET xp = xp + $2, gold = gold + $3, streak=$4, last_active_date=$5 WHERE id=$1`,
		userID, xp, gold, newStreak, today); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// UpdateStreak usado pra abandonar gate explícito
func UpdateStreak(ctx context.Context, db *pgxpool.Pool, userID uuid.UUID, success bool) error {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	today := time.Now().In(loc).Truncate(24 * time.Hour)

	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var oldStreak int
	if err := tx.QueryRow(ctx, `SELECT streak FROM users WHERE id=$1 FOR UPDATE`, userID).Scan(&oldStreak); err != nil {
		return err
	}

	newStreak := oldStreak
	if success {
		newStreak = oldStreak + 1
	} else {
		newStreak = 0
	}

	if _, err := tx.Exec(ctx, `UPDATE users SET streak=$2, last_active_date=$3 WHERE id=$1`,
		userID, newStreak, today); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
