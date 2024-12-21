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

type PaymentMethod string

const (
	PaymentMethodApplePay PaymentMethod = "apple_pay"
	PaymentMethodCard     PaymentMethod = "card"
)

type Payment struct {
	ID        uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID     `json:"userId" gorm:"type:uuid"`
	Amount    float64       `json:"amount" gorm:"type:decimal(20,8)"`
	Method    PaymentMethod `json:"method" gorm:"type:varchar(50)"`
	Status    PaymentStatus `json:"status" gorm:"type:varchar(50)"`
	Reference string        `json:"reference" gorm:"type:varchar(255)"`
	CreatedAt time.Time     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time     `json:"updatedAt" gorm:"autoUpdateTime"`
}

type PaymentResult struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	Amount        string `json:"amount"`
}
