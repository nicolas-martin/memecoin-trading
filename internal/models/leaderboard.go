package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaderboardEntry struct {
	ID               uuid.UUID `json:"id" db:"id"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	Username         string    `json:"username" db:"username"`
	Profit           float64   `json:"profit" db:"profit"`
	ProfitPercentage float64   `json:"profit_percentage" db:"profit_percentage"`
	Rank             int       `json:"rank" db:"rank"`
	TimeFrame        string    `json:"timeframe" db:"timeframe"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}
