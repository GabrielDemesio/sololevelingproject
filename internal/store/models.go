package store

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID      `json:"id"`
	Email          string         `json:"email"`
	PassHash       string         `json:"-"`
	Level          int            `json:"level"`
	XP             int64          `json:"xp"`
	Gold           int64          `json:"gold"`
	Stats          map[string]int `json:"stats"`
	Streak         int            `json:"streak"`
	LastActiveDate *time.Time     `json:"lastActiveDate,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
}

type Quest struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Weight      int        `json:"weight"`
	Status      string     `json:"status"`
	DueAt       *time.Time `json:"dueAt,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type FocusRun struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"userId"`
	QuestID       *uuid.UUID `json:"questId,omitempty"`
	DungeonRank   string     `json:"rank"`
	StartAt       time.Time  `json:"startAt"`
	EndAt         *time.Time `json:"endAt,omitempty"`
	TargetMinutes int        `json:"targetMinutes"`
	Result        *string    `json:"result,omitempty"`
	XPEarned      int64      `json:"xpEarned"`
	GoldEarned    int64      `json:"goldEarned"`
}
