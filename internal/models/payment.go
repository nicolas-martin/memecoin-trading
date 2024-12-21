package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID     `json:"user_id" gorm:"type:uuid"`
	Amount        float64       `json:"amount"`
	PaymentMethod string        `json:"payment_method"`
	TransactionID string        `json:"transaction_id"`
	Status        PaymentStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
}

type PaymentResult struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
}
