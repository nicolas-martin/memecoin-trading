package postgres

import (
	"context"
	"time"

	"github.com/nicolas-martin/memecoin-trading/internal/models"

	"gorm.io/gorm"
)

type LeaderboardRepository struct {
	db *gorm.DB
}

func NewLeaderboardRepository(db *gorm.DB) *LeaderboardRepository {
	return &LeaderboardRepository{db: db}
}

func (r *LeaderboardRepository) GetTopTraders(ctx context.Context, timeframe time.Duration, limit int) ([]models.LeaderboardEntry, error) {
	var entries []models.LeaderboardEntry

	query := `
		WITH trader_profits AS (
			SELECT 
				u.id as user_id,
				u.username,
				SUM(CASE 
					WHEN t.type = 'SELL' THEN t.amount * t.price
					WHEN t.type = 'BUY' THEN -t.amount * t.price
					ELSE 0
				END) as profit,
				ROW_NUMBER() OVER (ORDER BY SUM(CASE 
					WHEN t.type = 'SELL' THEN t.amount * t.price
					WHEN t.type = 'BUY' THEN -t.amount * t.price
					ELSE 0
				END) DESC) as rank
			FROM users u
			JOIN transactions t ON t.user_id = u.id
			WHERE t.created_at >= NOW() - $1::interval
				AND t.status = 'COMPLETED'
			GROUP BY u.id, u.username
		)
		SELECT 
			gen_random_uuid() as id,
			user_id,
			username,
			profit,
			rank,
			$2 as timeframe,
			NOW() as created_at
		FROM trader_profits
		LIMIT $3
	`

	err := r.db.WithContext(ctx).Raw(
		query,
		timeframe,
		timeframe.String(),
		limit,
	).Scan(&entries).Error

	if err != nil {
		return nil, err
	}

	return entries, nil
}
