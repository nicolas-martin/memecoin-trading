package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// Create creates a new wallet
func (r *WalletRepository) Create(ctx context.Context, wallet *models.Wallet) error {
	result := r.db.WithContext(ctx).Create(wallet)
	return result.Error
}

// GetByID retrieves a wallet by its ID
func (r *WalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	result := r.db.WithContext(ctx).First(&wallet, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &wallet, nil
}

// GetByUserID retrieves all wallets for a user
func (r *WalletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Wallet, error) {
	var wallets []models.Wallet
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&wallets)
	if result.Error != nil {
		return nil, result.Error
	}
	return wallets, nil
}

// GetByAddress retrieves a wallet by its address
func (r *WalletRepository) GetByAddress(ctx context.Context, address string) (*models.Wallet, error) {
	var wallet models.Wallet
	result := r.db.WithContext(ctx).First(&wallet, "address = ?", address)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &wallet, nil
}

// UpdateBalance updates the wallet balance within a transaction
func (r *WalletRepository) UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, "id = ?", id).Error; err != nil {
			return err
		}

		wallet.Balance += amount
		if wallet.Balance < 0 {
			return errors.New("insufficient balance")
		}

		return tx.Save(&wallet).Error
	})
}

// Update updates a wallet
func (r *WalletRepository) Update(ctx context.Context, wallet *models.Wallet) error {
	result := r.db.WithContext(ctx).Save(wallet)
	return result.Error
}

// Delete soft-deletes a wallet
func (r *WalletRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Wallet{}, "id = ?", id)
	return result.Error
}

// GetTotalBalance gets the total balance across all wallets for a user
func (r *WalletRepository) GetTotalBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	var totalBalance float64
	result := r.db.WithContext(ctx).Model(&models.Wallet{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&totalBalance)
	if result.Error != nil {
		return 0, result.Error
	}
	return totalBalance, nil
}

// TransferBalance transfers balance between two wallets within a transaction
func (r *WalletRepository) TransferBalance(ctx context.Context, fromID, toID uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Deduct from source wallet
		if err := r.UpdateBalance(ctx, fromID, -amount); err != nil {
			return err
		}

		// Add to destination wallet
		if err := r.UpdateBalance(ctx, toID, amount); err != nil {
			return err
		}

		return nil
	})
}

// GetWalletWithUser retrieves a wallet with its associated user
func (r *WalletRepository) GetWalletWithUser(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	result := r.db.WithContext(ctx).
		Preload("User").
		First(&wallet, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &wallet, nil
}
