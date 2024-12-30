package postgres

import (
	"fmt"
	"meme-trader/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

type MemeCoin struct {
	ID                       string `gorm:"primaryKey"`
	Symbol                   string
	Name                     string
	LogoURL                  string
	Price                    float64
	MarketCap                float64
	Volume24h                float64
	PriceChange24h           float64
	PriceChangePercentage24h float64
	ContractAddress          string
	Description              string
}

type PriceHistory struct {
	ID        uint `gorm:"primaryKey"`
	CoinID    string
	Price     float64
	Volume    float64
	Timestamp int64
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schemas
	err = db.AutoMigrate(&MemeCoin{}, &PriceHistory{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) GetTopMemeCoins(limit int) ([]MemeCoin, error) {
	var coins []MemeCoin
	err := d.DB.Order("market_cap DESC").Limit(limit).Find(&coins).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get top meme coins: %w", err)
	}
	return coins, nil
}

func (d *Database) GetMemeCoinByID(id string) (*MemeCoin, error) {
	var coin MemeCoin
	err := d.DB.First(&coin, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get meme coin: %w", err)
	}
	return &coin, nil
}

func (d *Database) GetPriceHistory(coinID string) ([]PriceHistory, error) {
	var history []PriceHistory
	err := d.DB.Where("coin_id = ?", coinID).Order("timestamp DESC").Find(&history).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}
	return history, nil
}

func (d *Database) UpdateMemeCoin(coin *MemeCoin) error {
	return d.DB.Save(coin).Error
}

func (d *Database) AddPriceHistory(history *PriceHistory) error {
	return d.DB.Create(history).Error
}
