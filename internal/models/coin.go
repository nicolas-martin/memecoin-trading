package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Coin struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Symbol          string         `json:"symbol" gorm:"type:varchar(50);not null"`
	Name            string         `json:"name" gorm:"type:varchar(255);not null"`
	ContractAddress string         `json:"contract_address" gorm:"type:varchar(255);unique;not null"`
	LogoURL         string         `json:"logo_url" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	Prices          []CoinPrice    `json:"prices,omitempty" gorm:"foreignKey:CoinID"`
}

type CoinPrice struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CoinID         uuid.UUID      `json:"coin_id" gorm:"type:uuid;not null"`
	Price          float64        `json:"price" gorm:"type:decimal(20,8);not null"`
	MarketCap      float64        `json:"market_cap" gorm:"type:decimal(30,2)"`
	Volume24h      float64        `json:"volume_24h" gorm:"type:decimal(30,2)"`
	PriceChange24h float64        `json:"price_change_24h" gorm:"type:decimal(20,8)"`
	Timestamp      time.Time      `json:"timestamp" gorm:"index"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	Coin           *Coin          `json:"-" gorm:"foreignKey:CoinID"`
}
