package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) AddFunds(ctx context.Context, amount float64, paymentMethod string, transactionID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create payment record
		payment := &models.Payment{
			ID:            uuid.New(),
			Amount:        amount,
			PaymentMethod: paymentMethod,
			TransactionID: transactionID,
			Status:        models.PaymentStatusCompleted,
			CreatedAt:     time.Now(),
		}

		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// Update user's balance
		if err := tx.Model(&models.User{}).
			Where("id = ?", payment.UserID).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount)).
			Error; err != nil {
			return err
		}

		return nil
	})
}
