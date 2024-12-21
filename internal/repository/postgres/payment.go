package postgres

import (
	"context"
	"fmt"

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

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	result := r.db.WithContext(ctx).Create(payment)
	if result.Error != nil {
		return fmt.Errorf("failed to create payment: %w", result.Error)
	}
	return nil
}

func (r *PaymentRepository) GetPayment(ctx context.Context, id string) (*models.Payment, error) {
	var payment models.Payment
	result := r.db.WithContext(ctx).First(&payment, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get payment: %w", result.Error)
	}
	return &payment, nil
}

func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	result := r.db.WithContext(ctx).Save(payment)
	if result.Error != nil {
		return fmt.Errorf("failed to update payment: %w", result.Error)
	}
	return nil
}

func (r *PaymentRepository) AddFunds(ctx context.Context, amount float64, paymentMethod string, transactionID string) error {
	payment := &models.Payment{
		ID:        uuid.New(),
		Amount:    amount,
		Method:    models.PaymentMethod(paymentMethod),
		Status:    models.PaymentStatusCompleted,
		Reference: transactionID,
	}

	return r.CreatePayment(ctx, payment)
}
