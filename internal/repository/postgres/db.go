package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

type MemeCoin struct {
	ID                       string
	Symbol                   string
	Name                     string
	Price                    float64
	MarketCap                float64
	Volume24h                float64
	PriceChange24h           float64
	PriceChangePercentage24h float64
	ContractAddress          string
	DataProvider             string
	LastUpdated              time.Time
	LogoURL                  string
	Description              string
}

type PriceHistory struct {
	CoinID    string
	Price     float64
	Volume    float64
	Timestamp int64
}

func NewDatabase(connStr string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	if err := createBlockchainTables(db); err != nil {
		return nil, fmt.Errorf("failed to create blockchain tables: %w", err)
	}

	return &Database{db: db}, nil
}

func createTables(db *sql.DB) error {
	// Create memecoins table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS memecoins (
			id TEXT PRIMARY KEY,
			symbol TEXT NOT NULL,
			name TEXT NOT NULL,
			price DOUBLE PRECISION NOT NULL,
			market_cap DOUBLE PRECISION,
			volume_24h DOUBLE PRECISION,
			price_change_24h DOUBLE PRECISION,
			price_change_percentage_24h DOUBLE PRECISION,
			contract_address TEXT NOT NULL,
			data_provider TEXT NOT NULL,
			last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			logo_url TEXT,
			description TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create memecoins table: %w", err)
	}

	// Create price_history table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS price_history (
			coin_id TEXT REFERENCES memecoins(id),
			price DOUBLE PRECISION NOT NULL,
			volume DOUBLE PRECISION,
			timestamp BIGINT NOT NULL,
			PRIMARY KEY (coin_id, timestamp)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create price_history table: %w", err)
	}

	return nil
}

func (db *Database) UpdateMemeCoin(coin *MemeCoin) error {
	log.Printf("Updating memecoin %s (%s) with logo URL: %s", coin.Name, coin.Symbol, coin.LogoURL)
	_, err := db.db.Exec(`
		INSERT INTO memecoins (
			id, symbol, name, price, market_cap, volume_24h,
			price_change_24h, price_change_percentage_24h, contract_address,
			data_provider, last_updated, logo_url, description
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
			symbol = EXCLUDED.symbol,
			name = EXCLUDED.name,
			price = EXCLUDED.price,
			market_cap = EXCLUDED.market_cap,
			volume_24h = EXCLUDED.volume_24h,
			price_change_24h = EXCLUDED.price_change_24h,
			price_change_percentage_24h = EXCLUDED.price_change_percentage_24h,
			contract_address = EXCLUDED.contract_address,
			data_provider = EXCLUDED.data_provider,
			last_updated = CURRENT_TIMESTAMP,
			logo_url = CASE 
				WHEN EXCLUDED.logo_url IS NOT NULL AND EXCLUDED.logo_url != '' THEN EXCLUDED.logo_url 
				ELSE memecoins.logo_url 
			END,
			description = EXCLUDED.description
	`, coin.ID, coin.Symbol, coin.Name, coin.Price, coin.MarketCap,
		coin.Volume24h, coin.PriceChange24h, coin.PriceChangePercentage24h,
		coin.ContractAddress, coin.DataProvider, coin.LogoURL, coin.Description)

	if err != nil {
		log.Printf("Error updating memecoin %s: %v", coin.Symbol, err)
		return fmt.Errorf("failed to update memecoin: %w", err)
	}

	// Verify the update
	var storedLogoURL string
	err = db.db.QueryRow("SELECT logo_url FROM memecoins WHERE id = $1", coin.ID).Scan(&storedLogoURL)
	if err != nil {
		log.Printf("Error verifying logo URL for %s: %v", coin.Symbol, err)
	} else {
		log.Printf("Stored logo URL for %s: %s", coin.Symbol, storedLogoURL)
	}

	return nil
}

func (db *Database) GetTopMemeCoins(limit int) ([]MemeCoin, error) {
	log.Printf("Fetching top %d meme coins", limit)
	rows, err := db.db.Query(`
		SELECT id, symbol, name, price, market_cap, volume_24h,
			price_change_24h, price_change_percentage_24h, contract_address,
			data_provider, last_updated, logo_url, description
		FROM memecoins
		ORDER BY market_cap DESC
		LIMIT $1
	`, limit)
	if err != nil {
		log.Printf("Error fetching top meme coins: %v", err)
		return nil, fmt.Errorf("failed to get top memecoins: %w", err)
	}
	defer rows.Close()

	var coins []MemeCoin
	for rows.Next() {
		var coin MemeCoin
		err := rows.Scan(
			&coin.ID, &coin.Symbol, &coin.Name, &coin.Price,
			&coin.MarketCap, &coin.Volume24h, &coin.PriceChange24h,
			&coin.PriceChangePercentage24h, &coin.ContractAddress,
			&coin.DataProvider, &coin.LastUpdated, &coin.LogoURL,
			&coin.Description,
		)
		if err != nil {
			log.Printf("Error scanning memecoin: %v", err)
			return nil, fmt.Errorf("failed to scan memecoin: %w", err)
		}
		log.Printf("Found coin %s (%s) with logo URL: %s", coin.Name, coin.Symbol, coin.LogoURL)
		coins = append(coins, coin)
	}

	log.Printf("Returning %d meme coins", len(coins))
	return coins, nil
}

func (db *Database) GetMemeCoinByID(id string) (*MemeCoin, error) {
	var coin MemeCoin
	err := db.db.QueryRow(`
		SELECT id, symbol, name, price, market_cap, volume_24h,
			price_change_24h, price_change_percentage_24h, contract_address,
			data_provider, last_updated, logo_url, description
		FROM memecoins
		WHERE id = $1
	`, id).Scan(
		&coin.ID, &coin.Symbol, &coin.Name, &coin.Price,
		&coin.MarketCap, &coin.Volume24h, &coin.PriceChange24h,
		&coin.PriceChangePercentage24h, &coin.ContractAddress,
		&coin.DataProvider, &coin.LastUpdated, &coin.LogoURL,
		&coin.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get memecoin: %w", err)
	}

	return &coin, nil
}

func (db *Database) AddPriceHistory(history *PriceHistory) error {
	_, err := db.db.Exec(`
		INSERT INTO price_history (coin_id, price, volume, timestamp)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (coin_id, timestamp) DO UPDATE SET
			price = EXCLUDED.price,
			volume = EXCLUDED.volume
	`, history.CoinID, history.Price, history.Volume, history.Timestamp)

	if err != nil {
		return fmt.Errorf("failed to add price history: %w", err)
	}

	return nil
}

func (db *Database) GetPriceHistory(coinID string) ([]PriceHistory, error) {
	rows, err := db.db.Query(`
		SELECT coin_id, price, volume, timestamp
		FROM price_history
		WHERE coin_id = $1
		ORDER BY timestamp DESC
		LIMIT 100
	`, coinID)
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %w", err)
	}
	defer rows.Close()

	var history []PriceHistory
	for rows.Next() {
		var h PriceHistory
		err := rows.Scan(&h.CoinID, &h.Price, &h.Volume, &h.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan price history: %w", err)
		}
		history = append(history, h)
	}

	return history, nil
}
