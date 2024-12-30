package blockchain

import (
	"context"
	"time"
)

type service struct {
	manager *ProviderManager
}

// NewService creates a new blockchain service
func NewService() Service {
	return &service{
		manager: NewProviderManager(),
	}
}

// RegisterProvider registers a new blockchain provider with default configuration
func (s *service) RegisterProvider(provider Provider) error {
	config := ProviderConfig{
		Priority:           10, // Default priority
		RequestsPerWindow:  100,
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 3,
	}
	return s.manager.RegisterProvider(provider, config)
}

// RegisterProviderWithConfig registers a new blockchain provider with custom configuration
func (s *service) RegisterProviderWithConfig(provider Provider, config ProviderConfig) error {
	return s.manager.RegisterProvider(provider, config)
}

// CreateWallet creates a new wallet for the specified network
func (s *service) CreateWallet(ctx context.Context, network Network) (*Wallet, error) {
	var wallet *Wallet
	err := s.manager.executeWithFallback(ctx, network, func(provider Provider) error {
		var err error
		wallet, err = provider.CreateWallet(ctx)
		return err
	})
	return wallet, err
}

// GetWallet retrieves a wallet by its address
func (s *service) GetWallet(ctx context.Context, network Network, address string) (*Wallet, error) {
	var wallet *Wallet
	err := s.manager.executeWithFallback(ctx, network, func(provider Provider) error {
		var err error
		wallet, err = provider.GetWallet(ctx, address)
		return err
	})
	return wallet, err
}

// GetBalance retrieves the balance for a wallet
func (s *service) GetBalance(ctx context.Context, network Network, address string) (Amount, error) {
	var balance Amount
	err := s.manager.executeWithFallback(ctx, network, func(provider Provider) error {
		var err error
		balance, err = provider.GetBalance(ctx, address)
		return err
	})
	return balance, err
}

// Buy executes a buy transaction
func (s *service) Buy(ctx context.Context, network Network, req BuyRequest) (*Transaction, error) {
	var tx *Transaction
	err := s.manager.executeWithFallback(ctx, network, func(provider Provider) error {
		var err error
		tx, err = provider.Buy(ctx, req)
		return err
	})
	return tx, err
}

// Sell executes a sell transaction
func (s *service) Sell(ctx context.Context, network Network, req SellRequest) (*Transaction, error) {
	var tx *Transaction
	err := s.manager.executeWithFallback(ctx, network, func(provider Provider) error {
		var err error
		tx, err = provider.Sell(ctx, req)
		return err
	})
	return tx, err
}

// GetTransaction retrieves a transaction by its ID
func (s *service) GetTransaction(ctx context.Context, network Network, txID string) (*Transaction, error) {
	var tx *Transaction
	err := s.manager.executeWithFallback(ctx, network, func(provider Provider) error {
		var err error
		tx, err = provider.GetTransaction(ctx, txID)
		return err
	})
	return tx, err
}

// GetTransactions retrieves transactions for a wallet
func (s *service) GetTransactions(ctx context.Context, network Network, address string, limit int) ([]Transaction, error) {
	var txs []Transaction
	err := s.manager.executeWithFallback(ctx, network, func(provider Provider) error {
		var err error
		txs, err = provider.GetTransactions(ctx, address, limit)
		return err
	})
	return txs, err
}
