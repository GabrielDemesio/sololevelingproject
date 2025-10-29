package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateFocusRun(ctx context.Context, db *pgxpool.Pool, r *FocusRun) error {
	_, err := db.Exec(ctx, `INSERT INTO focus_runs(id, user_id, quest_id, dungeon_rank, start_at, target_minutes, xp_earned, gold_earned)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8)`,
		r.ID, r.UserID, r.QuestID, r.DungeonRank, r.StartAt, r.TargetMinutes, r.XPEarned, r.GoldEarned)
	return err
}

func GetFocusRunByID(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) (FocusRun, error) {
	row := db.QueryRow(ctx, `SELECT id, user_id, quest_id, dungeon_rank, start_at, end_at, target_minutes, result, xp_earned, gold_earned
	FROM focus_runs WHERE id=$1`, id)

	var r FocusRun
	if err := row.Scan(&r.ID, &r.UserID, &r.QuestID, &r.DungeonRank, &r.StartAt, &r.EndAt,
		&r.TargetMinutes, &r.Result, &r.XPEarned, &r.GoldEarned); err != nil {
		return FocusRun{}, err
	}
	return r, nil
}

func FinishFocusRun(ctx context.Context, db *pgxpool.Pool, r *FocusRun) error {
	now := time.Now()
	if r.EndAt == nil {
		r.EndAt = &now
	}
	_, err := db.Exec(ctx, `UPDATE focus_runs
	SET end_at=$2, result=$3, xp_earned=$4, gold_earned=$5
	WHERE id=$1`,
		r.ID, r.EndAt, r.Result, r.XPEarned, r.GoldEarned)
	return err
}
