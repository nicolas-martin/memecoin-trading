package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PriceChange struct {
	H1  float64 `json:"h1"`
	H24 float64 `json:"h24"`
	D7  float64 `json:"d7"`
}

type Volume struct {
	H24 float64 `json:"h24"`
	H6  float64 `json:"h6"`
	H1  float64 `json:"h1"`
	M5  float64 `json:"m5"`
}

type Liquidity struct {
	USD   float64 `json:"usd"`
	Base  float64 `json:"base"`
	Quote float64 `json:"quote"`
}

type Coin struct {
	ID          string      `json:"id" gorm:"primaryKey"`
	Name        string      `json:"name" gorm:"type:varchar(255)"`
	Symbol      string      `json:"symbol" gorm:"type:varchar(50)"`
	PairAddress string      `json:"pairAddress" gorm:"type:varchar(255)"`
	ChainID     string      `json:"chainId" gorm:"type:varchar(50)"`
	Price       string      `json:"price" gorm:"type:varchar(50)"`
	PriceChange PriceChange `json:"priceChange" gorm:"embedded"`
	Volume      Volume      `json:"volume" gorm:"embedded"`
	Liquidity   Liquidity   `json:"liquidity" gorm:"embedded"`
	MarketCap   float64     `json:"marketCap" gorm:"type:decimal(20,2)"`
	FDV         float64     `json:"fdv" gorm:"type:decimal(20,2)"`
	Logo        string      `json:"logo,omitempty" gorm:"type:text"`
	Description string      `json:"description,omitempty" gorm:"type:text"`
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
