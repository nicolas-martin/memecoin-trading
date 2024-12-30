package blockchain

import (
	"context"
	"math/big"
)

// Network represents a blockchain network
type Network string

const (
	NetworkSolana Network = "solana"
	// Add more networks as needed:
	// NetworkEthereum Network = "ethereum"
	// NetworkPolygon Network = "polygon"
)

// Amount represents a blockchain amount with high precision
type Amount struct {
	Value    *big.Int
	Decimals uint8
}

// Wallet represents a blockchain wallet
type Wallet struct {
	ID            string
	Network       Network
	Address       string
	PublicKey     string
	PrivateKey    string // Should be encrypted in production
	CreatedAt     int64
	LastUpdatedAt int64
}

// Transaction represents a blockchain transaction
type Transaction struct {
	ID            string
	Network       Network
	Type          TransactionType
	Status        TransactionStatus
	FromAddress   string
	ToAddress     string
	Amount        Amount
	TokenAddress  string
	Signature     string
	BlockHash     string
	BlockNumber   uint64
	Timestamp     int64
	GasFee        Amount
	ErrorMessage  string
	CreatedAt     int64
	LastUpdatedAt int64
}

type TransactionType string

const (
	TransactionTypeBuy  TransactionType = "buy"
	TransactionTypeSell TransactionType = "sell"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// Provider defines the interface for blockchain-specific implementations
type Provider interface {
	// Network information
	Network() Network
	IsValidAddress(address string) bool

	// Wallet operations
	CreateWallet(ctx context.Context) (*Wallet, error)
	GetWallet(ctx context.Context, address string) (*Wallet, error)
	GetBalance(ctx context.Context, address string) (Amount, error)

	// Transaction operations
	Buy(ctx context.Context, req BuyRequest) (*Transaction, error)
	Sell(ctx context.Context, req SellRequest) (*Transaction, error)
	GetTransaction(ctx context.Context, txID string) (*Transaction, error)
	GetTransactions(ctx context.Context, address string, limit int) ([]Transaction, error)
}

// BuyRequest represents a request to buy tokens
type BuyRequest struct {
	WalletAddress string
	TokenAddress  string
	Amount        Amount
	MaxPrice      Amount // Maximum price willing to pay (slippage protection)
}

// SellRequest represents a request to sell tokens
type SellRequest struct {
	WalletAddress string
	TokenAddress  string
	Amount        Amount
	MinPrice      Amount // Minimum price willing to accept (slippage protection)
}

// Service provides a high-level interface for blockchain operations
type Service interface {
	// Provider management
	RegisterProvider(provider Provider) error
	RegisterProviderWithConfig(provider Provider, config ProviderConfig) error

	// Wallet operations
	CreateWallet(ctx context.Context, network Network) (*Wallet, error)
	GetWallet(ctx context.Context, network Network, address string) (*Wallet, error)
	GetBalance(ctx context.Context, network Network, address string) (Amount, error)

	// Transaction operations
	Buy(ctx context.Context, network Network, req BuyRequest) (*Transaction, error)
	Sell(ctx context.Context, network Network, req SellRequest) (*Transaction, error)
	GetTransaction(ctx context.Context, network Network, txID string) (*Transaction, error)
	GetTransactions(ctx context.Context, network Network, address string, limit int) ([]Transaction, error)
}
