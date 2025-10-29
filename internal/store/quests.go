package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListQuestsByUser(ctx context.Context, db *pgxpool.Pool, userID uuid.UUID) ([]Quest, error) {
	rows, err := db.Query(ctx, `SELECT id, user_id, title, description, weight, status, due_at, tags, created_at
	FROM quests WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Quest
	for rows.Next() {
		var q Quest
		if err := rows.Scan(&q.ID, &q.UserID, &q.Title, &q.Description, &q.Weight, &q.Status,
			&q.DueAt, &q.Tags, &q.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, nil
}

func CreateQuest(ctx context.Context, db *pgxpool.Pool, q *Quest) error {
	_, err := db.Exec(ctx, `INSERT INTO quests(id, user_id, title, description, weight, status, due_at, tags, created_at)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		q.ID, q.UserID, q.Title, q.Description, q.Weight, q.Status, q.DueAt, q.Tags, q.CreatedAt)
	return err
}

func GetQuestByID(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) (*Quest, error) {
	row := db.QueryRow(ctx, `SELECT id, user_id, title, description, weight, status, due_at, tags, created_at
	FROM quests WHERE id=$1`, id)

	var q Quest
	if err := row.Scan(&q.ID, &q.UserID, &q.Title, &q.Description, &q.Weight, &q.Status,
		&q.DueAt, &q.Tags, &q.CreatedAt); err != nil {
		return nil, err
	}
	return &q, nil
}

func UpdateQuest(ctx context.Context, db *pgxpool.Pool, q *Quest) error {
	_, err := db.Exec(ctx, `UPDATE quests
	SET title=$2, description=$3, weight=$4, status=$5, tags=$6
	WHERE id=$1`,
		q.ID, q.Title, q.Description, q.Weight, q.Status, q.Tags)
	return err
}

func DeleteQuest(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) error {
	_, err := db.Exec(ctx, `DELETE FROM quests WHERE id=$1`, id)
	return err
}

// usado pelo cálculo de recompensa
func GetQuestWeight(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) (int, error) {
	row := db.QueryRow(ctx, `SELECT weight FROM quests WHERE id=$1`, id)
	var w int
	if err := row.Scan(&w); err != nil {
		return 1, err
	}
	if w <= 0 {
		w = 1
	}
	return w, nil
}
