package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"gorm.io/gorm"
)

type PortfolioRepository struct {
	db *gorm.DB
}

func NewPortfolioRepository(db *gorm.DB) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

func (r *PortfolioRepository) GetHoldings(ctx context.Context, userID uuid.UUID) ([]models.PortfolioHolding, error) {
	var holdings []models.PortfolioHolding

	query := `
		WITH current_prices AS (
			SELECT DISTINCT ON (coin_id) 
				coin_id,
				price as current_price
			FROM coin_prices
			ORDER BY coin_id, timestamp DESC
		)
		SELECT 
			h.id,
			h.user_id,
			h.coin_id,
			h.amount,
			h.average_price,
			cp.current_price,
			h.amount * cp.current_price as value,
			(h.amount * cp.current_price) - (h.amount * h.average_price) as profit_loss,
			CASE 
				WHEN h.average_price > 0 THEN 
					((cp.current_price - h.average_price) / h.average_price) * 100
				ELSE 0 
			END as profit_loss_percent,
			h.updated_at
		FROM holdings h
		JOIN current_prices cp ON cp.coin_id = h.coin_id
		WHERE h.user_id = $1
		AND h.amount > 0
	`

	err := r.db.WithContext(ctx).Raw(query, userID).Scan(&holdings).Error
	if err != nil {
		return nil, err
	}

	return holdings, nil
}

func (r *PortfolioRepository) GetHistory(ctx context.Context, userID uuid.UUID, timeframe string) ([]models.PortfolioValue, error) {
	var history []models.PortfolioValue

	query := `
		SELECT 
			id,
			user_id,
			value,
			timestamp
		FROM portfolio_values
		WHERE user_id = $1
		AND timestamp >= NOW() - $2::interval
		ORDER BY timestamp ASC
	`

	err := r.db.WithContext(ctx).Raw(query, userID, timeframe).Scan(&history).Error
	if err != nil {
		return nil, err
	}

	return history, nil
}
