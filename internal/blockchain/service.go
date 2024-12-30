package blockchain

import (
	"context"
	"fmt"
	"sync"
)

type service struct {
	providers map[Network]Provider
	mu        sync.RWMutex
}

// NewService creates a new blockchain service
func NewService() Service {
	return &service{
		providers: make(map[Network]Provider),
	}
}

// RegisterProvider registers a new blockchain provider
func (s *service) RegisterProvider(provider Provider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	network := provider.Network()
	if _, exists := s.providers[network]; exists {
		return fmt.Errorf("provider for network %s already registered", network)
	}

	s.providers[network] = provider
	return nil
}

// GetProvider returns the provider for the specified network
func (s *service) GetProvider(network Network) (Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	provider, exists := s.providers[network]
	if !exists {
		return nil, fmt.Errorf("no provider registered for network %s", network)
	}

	return provider, nil
}

// CreateWallet creates a new wallet for the specified network
func (s *service) CreateWallet(ctx context.Context, network Network) (*Wallet, error) {
	provider, err := s.GetProvider(network)
	if err != nil {
		return nil, err
	}

	return provider.CreateWallet(ctx)
}

// GetWallet retrieves a wallet by its address
func (s *service) GetWallet(ctx context.Context, network Network, address string) (*Wallet, error) {
	provider, err := s.GetProvider(network)
	if err != nil {
		return nil, err
	}

	return provider.GetWallet(ctx, address)
}

// GetBalance retrieves the balance for a wallet
func (s *service) GetBalance(ctx context.Context, network Network, address string) (Amount, error) {
	provider, err := s.GetProvider(network)
	if err != nil {
		return Amount{}, err
	}

	return provider.GetBalance(ctx, address)
}

// Buy executes a buy transaction
func (s *service) Buy(ctx context.Context, network Network, req BuyRequest) (*Transaction, error) {
	provider, err := s.GetProvider(network)
	if err != nil {
		return nil, err
	}

	return provider.Buy(ctx, req)
}

// Sell executes a sell transaction
func (s *service) Sell(ctx context.Context, network Network, req SellRequest) (*Transaction, error) {
	provider, err := s.GetProvider(network)
	if err != nil {
		return nil, err
	}

	return provider.Sell(ctx, req)
}

// GetTransaction retrieves a transaction by its ID
func (s *service) GetTransaction(ctx context.Context, network Network, txID string) (*Transaction, error) {
	provider, err := s.GetProvider(network)
	if err != nil {
		return nil, err
	}

	return provider.GetTransaction(ctx, txID)
}

// GetTransactions retrieves transactions for a wallet
func (s *service) GetTransactions(ctx context.Context, network Network, address string, limit int) ([]Transaction, error) {
	provider, err := s.GetProvider(network)
	if err != nil {
		return nil, err
	}

	return provider.GetTransactions(ctx, address, limit)
}
