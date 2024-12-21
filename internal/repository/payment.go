package repository

import (
	"context"

	"github.com/nicolas-martin/memecoin-trading/internal/models"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPayment(ctx context.Context, id string) (*models.Payment, error)
	UpdatePayment(ctx context.Context, payment *models.Payment) error
	AddFunds(ctx context.Context, amount float64, paymentMethod string, transactionID string) error
}
