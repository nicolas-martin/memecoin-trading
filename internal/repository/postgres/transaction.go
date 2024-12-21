package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"gorm.io/gorm"
)

const (
	defaultBatchSize = 1000
	maxQueryTimeout  = 30 * time.Second
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Transaction(func(dtx *gorm.DB) error {
		// Create the transaction record
		if err := dtx.Create(tx).Error; err != nil {
			return err
		}

		// Update wallet balance based on transaction type
		var balanceChange float64
		if tx.Type == models.TransactionTypeBuy {
			balanceChange = -tx.Amount * tx.Price // Deduct for buy
		} else {
			balanceChange = tx.Amount * tx.Price // Add for sell
		}

		// Update wallet balance
		if err := dtx.Model(&models.Wallet{}).
			Where("id = ?", tx.WalletID).
			UpdateColumn("balance", gorm.Expr("balance + ?", balanceChange)).
			Error; err != nil {
			return err
		}

		return nil
	})
}

// GetByID retrieves a transaction by its ID
func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var tx models.Transaction
	result := r.db.WithContext(ctx).
		Preload("User").
		Preload("Wallet").
		Preload("Coin").
		First(&tx, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &tx, nil
}

// GetByUserID retrieves all transactions for a user
func (r *TransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := r.db.WithContext(ctx).
		Preload("Coin").
		Preload("Wallet").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

// GetByWalletID retrieves all transactions for a wallet
func (r *TransactionRepository) GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := r.db.WithContext(ctx).
		Preload("Coin").
		Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

// UpdateStatus updates the transaction status
func (r *TransactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TransactionStatus, txHash *string) error {
	return r.db.WithContext(ctx).Transaction(func(dtx *gorm.DB) error {
		result := dtx.Model(&models.Transaction{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"status":  status,
				"tx_hash": txHash,
			})
		return result.Error
	})
}

// GetUserStats retrieves user transaction statistics
func (r *TransactionRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (map[string]float64, error) {
	var stats struct {
		TotalBuyAmount  float64 `gorm:"column:total_buy_amount"`
		TotalSellAmount float64 `gorm:"column:total_sell_amount"`
		TotalBuyValue   float64 `gorm:"column:total_buy_value"`
		TotalSellValue  float64 `gorm:"column:total_sell_value"`
	}

	result := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Select(`
			SUM(CASE WHEN type = 'BUY' THEN amount ELSE 0 END) as total_buy_amount,
			SUM(CASE WHEN type = 'SELL' THEN amount ELSE 0 END) as total_sell_amount,
			SUM(CASE WHEN type = 'BUY' THEN amount * price ELSE 0 END) as total_buy_value,
			SUM(CASE WHEN type = 'SELL' THEN amount * price ELSE 0 END) as total_sell_value
		`).
		Where("user_id = ? AND status = ?", userID, models.TransactionStatusCompleted).
		Scan(&stats)

	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]float64{
		"total_buy_amount":  stats.TotalBuyAmount,
		"total_sell_amount": stats.TotalSellAmount,
		"total_buy_value":   stats.TotalBuyValue,
		"total_sell_value":  stats.TotalSellValue,
		"net_profit":        stats.TotalSellValue - stats.TotalBuyValue,
	}, nil
}

// Add this struct for the query result
type ProfitResult struct {
	UserID uuid.UUID `gorm:"column:user_id"`
	Profit float64   `gorm:"column:profit"`
}

// Update the method to use the new struct
func (r *TransactionRepository) GetTopTraders(ctx context.Context, duration time.Duration, limit int) ([]ProfitResult, error) {
	var results []ProfitResult

	query := `
		SELECT 
			user_id,
			SUM(CASE 
				WHEN type = 'SELL' THEN amount * price
				WHEN type = 'BUY' THEN -amount * price
					ELSE 0
			END) as profit
		FROM transactions
		WHERE created_at >= NOW() - ?::interval
			AND status = 'COMPLETED'
		GROUP BY user_id
		ORDER BY profit DESC
		LIMIT ?
	`

	err := r.db.WithContext(ctx).Raw(query, duration.String(), limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetTransactionHistory retrieves transaction history with pagination and filters
func (r *TransactionRepository) GetTransactionHistory(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Transaction{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	result := query.
		Preload("Coin").
		Preload("Wallet").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return transactions, total, nil
}

// BatchCreate creates multiple transactions in batches
func (r *TransactionRepository) BatchCreate(ctx context.Context, transactions []*models.Transaction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Process in batches to avoid memory issues
		for i := 0; i < len(transactions); i += defaultBatchSize {
			end := i + defaultBatchSize
			if end > len(transactions) {
				end = len(transactions)
			}

			if err := tx.CreateInBatches(transactions[i:end], defaultBatchSize).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByIDWithSelect retrieves a transaction by ID with specific fields
func (r *TransactionRepository) GetByIDWithSelect(ctx context.Context, id uuid.UUID, fields []string) (*models.Transaction, error) {
	var tx models.Transaction
	result := r.db.WithContext(ctx).
		Select(fields).
		First(&tx, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &tx, nil
}

// Optimize GetTransactionHistory with cursor-based pagination
type TransactionCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

func (r *TransactionRepository) GetTransactionHistoryWithCursor(
	ctx context.Context,
	filters map[string]interface{},
	cursor *TransactionCursor,
	limit int,
) ([]models.Transaction, *TransactionCursor, error) {
	query := r.db.WithContext(ctx).Model(&models.Transaction{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	// Apply cursor conditions if provided
	if cursor != nil {
		query = query.Where(
			"(created_at, id) < (?, ?)",
			cursor.CreatedAt,
			cursor.ID,
		)
	}

	var transactions []models.Transaction
	result := query.
		Preload("Coin").
		Preload("Wallet").
		Order("created_at DESC, id DESC").
		Limit(limit + 1). // Get one extra to determine if there are more results
		Find(&transactions)

	if result.Error != nil {
		return nil, nil, result.Error
	}

	var nextCursor *TransactionCursor
	if len(transactions) > limit {
		lastTx := transactions[limit-1]
		nextCursor = &TransactionCursor{
			CreatedAt: lastTx.CreatedAt,
			ID:        lastTx.ID,
		}
		transactions = transactions[:limit]
	}

	return transactions, nextCursor, nil
}

// Optimize GetUserStats with materialized view
func (r *TransactionRepository) GetUserStatsOptimized(ctx context.Context, userID uuid.UUID) (map[string]float64, error) {
	// Using a CTE for better performance
	query := `
		WITH user_stats AS (
			SELECT 
				SUM(CASE WHEN type = 'BUY' THEN amount ELSE 0 END) as total_buy_amount,
				SUM(CASE WHEN type = 'SELL' THEN amount ELSE 0 END) as total_sell_amount,
				SUM(CASE WHEN type = 'BUY' THEN amount * price ELSE 0 END) as total_buy_value,
				SUM(CASE WHEN type = 'SELL' THEN amount * price ELSE 0 END) as total_sell_value
			FROM transactions 
			WHERE user_id = ? 
				AND status = 'COMPLETED'
				AND created_at >= NOW() - INTERVAL '30 days'
		)
		SELECT 
			COALESCE(total_buy_amount, 0) as total_buy_amount,
			COALESCE(total_sell_amount, 0) as total_sell_amount,
			COALESCE(total_buy_value, 0) as total_buy_value,
			COALESCE(total_sell_value, 0) as total_sell_value
		FROM user_stats
	`

	var stats struct {
		TotalBuyAmount  float64 `gorm:"column:total_buy_amount"`
		TotalSellAmount float64 `gorm:"column:total_sell_amount"`
		TotalBuyValue   float64 `gorm:"column:total_buy_value"`
		TotalSellValue  float64 `gorm:"column:total_sell_value"`
	}

	ctx, cancel := context.WithTimeout(ctx, maxQueryTimeout)
	defer cancel()

	if err := r.db.WithContext(ctx).Raw(query, userID).Scan(&stats).Error; err != nil {
		return nil, err
	}

	return map[string]float64{
		"total_buy_amount":  stats.TotalBuyAmount,
		"total_sell_amount": stats.TotalSellAmount,
		"total_buy_value":   stats.TotalBuyValue,
		"total_sell_value":  stats.TotalSellValue,
		"net_profit":        stats.TotalSellValue - stats.TotalBuyValue,
	}, nil
}

// Add query timeout wrapper
func (r *TransactionRepository) withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); !ok {
		return context.WithTimeout(ctx, timeout)
	}
	return ctx, func() {}
}
