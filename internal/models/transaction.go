package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeBuy  TransactionType = "BUY"
	TransactionTypeSell TransactionType = "SELL"

	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

type Transaction struct {
	ID        uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID         `json:"user_id" gorm:"type:uuid;not null"`
	WalletID  uuid.UUID         `json:"wallet_id" gorm:"type:uuid;not null"`
	CoinID    uuid.UUID         `json:"coin_id" gorm:"type:uuid;not null"`
	Type      TransactionType   `json:"type" gorm:"type:varchar(20);not null"`
	Amount    float64           `json:"amount" gorm:"type:decimal(20,8);not null"`
	Price     float64           `json:"price" gorm:"type:decimal(20,8);not null"`
	Status    TransactionStatus `json:"status" gorm:"type:varchar(20);not null"`
	TxHash    *string           `json:"tx_hash,omitempty" gorm:"type:varchar(255)"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `json:"-" gorm:"index"`
	User      *User             `json:"-" gorm:"foreignKey:UserID"`
	Wallet    *Wallet           `json:"-" gorm:"foreignKey:WalletID"`
	Coin      *Coin             `json:"-" gorm:"foreignKey:CoinID"`
}
