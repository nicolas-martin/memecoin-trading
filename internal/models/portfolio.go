package models

import (
	"time"

	"github.com/google/uuid"
)

type PortfolioHolding struct {
	ID                uuid.UUID `json:"id" db:"id"`
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	CoinID            uuid.UUID `json:"coin_id" db:"coin_id"`
	Amount            float64   `json:"amount" db:"amount"`
	AveragePrice      float64   `json:"average_price" db:"average_price"`
	CurrentPrice      float64   `json:"current_price" db:"current_price"`
	Value             float64   `json:"value" db:"value"`
	ProfitLoss        float64   `json:"profit_loss" db:"profit_loss"`
	ProfitLossPercent float64   `json:"profit_loss_percent" db:"profit_loss_percent"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type PortfolioValue struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Value     float64   `json:"value" db:"value"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}
