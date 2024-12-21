package payment

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/nicolas-martin/memecoin-trading/internal/config"
	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"github.com/nicolas-martin/memecoin-trading/internal/repository/postgres"
	"github.com/nicolas-martin/memecoin-trading/internal/repository/redis"
)

type Service struct {
	repo           *postgres.PaymentRepository
	cache          redis.Cache
	applePayConfig *config.ApplePayConfig
}

func NewService(repo *postgres.PaymentRepository, cache redis.Cache, applePayConfig *config.ApplePayConfig) *Service {
	return &Service{
		repo:           repo,
		cache:          cache,
		applePayConfig: applePayConfig,
	}
}

func (s *Service) ValidateApplePayMerchant(ctx context.Context, validationURL string) (map[string]interface{}, error) {
	// Load your Apple Pay merchant certificate and private key
	cert, err := tls.LoadX509KeyPair(s.applePayConfig.CertificatePath, s.applePayConfig.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load merchant certificate: %w", err)
	}

	// Create HTTP client with certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	// Create merchant validation payload
	payload := map[string]interface{}{
		"merchantIdentifier": s.applePayConfig.MerchantID,
		"displayName":        "MemeCoin Trading",
		"initiative":         "web",
		"initiativeContext":  s.applePayConfig.DomainName,
	}

	// Send validation request to Apple
	response, err := s.sendValidationRequest(ctx, client, validationURL, payload)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) ProcessApplePayPayment(ctx context.Context, payment map[string]interface{}, amount float64) (*models.PaymentResult, error) {
	// Verify the payment data
	if err := s.verifyPayment(payment, amount); err != nil {
		return nil, err
	}

	// Process the payment with your payment processor
	result, err := s.processPayment(ctx, payment, amount)
	if err != nil {
		return nil, err
	}

	// Store the payment record
	paymentRecord := &models.Payment{
		UserID:    getUserIDFromContext(ctx),
		Amount:    amount,
		Method:    "apple_pay",
		Status:    "completed",
		Reference: result.TransactionID,
	}

	if err := s.repo.CreatePayment(ctx, paymentRecord); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) AddFunds(ctx context.Context, amount float64, paymentMethod string, transactionID string) error {
	// Implementation for adding funds to user's account
	return s.repo.AddFunds(ctx, amount, paymentMethod, transactionID)
}
