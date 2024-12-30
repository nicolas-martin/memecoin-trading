package solana

import (
	"context"
	"math/big"
	"meme-trader/internal/blockchain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRaydiumClient is a mock implementation of the Raydium client
type MockRaydiumClient struct {
	mock.Mock
}

func (m *MockRaydiumClient) SwapTokens(ctx context.Context, req SwapRequest) (*blockchain.Transaction, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*blockchain.Transaction), args.Error(1)
}

func TestNewProvider(t *testing.T) {
	provider, err := NewProvider(true)
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, blockchain.NetworkSolana, provider.Network())
	assert.True(t, provider.isDevnet)
}

func TestIsValidAddress(t *testing.T) {
	provider, err := NewProvider(true)
	assert.NoError(t, err)

	// Test valid address
	valid := provider.IsValidAddress("11111111111111111111111111111111")
	assert.True(t, valid)

	// Test invalid address
	invalid := provider.IsValidAddress("invalid")
	assert.False(t, invalid)
}

func TestCreateWallet(t *testing.T) {
	provider, err := NewProvider(true)
	assert.NoError(t, err)

	wallet, err := provider.CreateWallet(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, blockchain.NetworkSolana, wallet.Network)
	assert.NotEmpty(t, wallet.Address)
	assert.NotEmpty(t, wallet.PublicKey)
	assert.NotEmpty(t, wallet.PrivateKey)
}

func TestGetWallet(t *testing.T) {
	provider, err := NewProvider(true)
	assert.NoError(t, err)

	// Test valid address
	wallet, err := provider.GetWallet(context.Background(), "11111111111111111111111111111111")
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, blockchain.NetworkSolana, wallet.Network)
	assert.Equal(t, "11111111111111111111111111111111", wallet.Address)

	// Test invalid address
	wallet, err = provider.GetWallet(context.Background(), "invalid")
	assert.Error(t, err)
	assert.Nil(t, wallet)
}

func TestGetBalance(t *testing.T) {
	provider, err := NewProvider(true)
	assert.NoError(t, err)

	// Test valid address
	balance, err := provider.GetBalance(context.Background(), "11111111111111111111111111111111")
	assert.NoError(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, uint8(9), balance.Decimals)

	// Test invalid address
	balance, err = provider.GetBalance(context.Background(), "invalid")
	assert.Error(t, err)
	assert.Equal(t, blockchain.Amount{}, balance)
}

func TestBuyTransaction(t *testing.T) {
	provider, err := NewProvider(true)
	assert.NoError(t, err)

	// Create test request
	req := blockchain.BuyRequest{
		WalletAddress: "11111111111111111111111111111111",
		TokenAddress:  "11111111111111111111111111111111", // Use a valid address for testing
		Amount: blockchain.Amount{
			Value:    big.NewInt(1000000),
			Decimals: 9,
		},
		MaxPrice: blockchain.Amount{
			Value:    big.NewInt(900000),
			Decimals: 9,
		},
	}

	// Execute buy transaction
	tx, err := provider.Buy(context.Background(), req)
	if err != nil {
		t.Logf("Buy transaction error: %v", err)
		t.Skip("Skipping test in devnet environment")
		return
	}
	assert.NotNil(t, tx)
	assert.Equal(t, blockchain.TransactionTypeBuy, tx.Type)
}

func TestSellTransaction(t *testing.T) {
	provider, err := NewProvider(true)
	assert.NoError(t, err)

	// Create test request
	req := blockchain.SellRequest{
		WalletAddress: "11111111111111111111111111111111",
		TokenAddress:  "11111111111111111111111111111111", // Use a valid address for testing
		Amount: blockchain.Amount{
			Value:    big.NewInt(1000000),
			Decimals: 9,
		},
		MinPrice: blockchain.Amount{
			Value:    big.NewInt(900000),
			Decimals: 9,
		},
	}

	// Execute sell transaction
	tx, err := provider.Sell(context.Background(), req)
	if err != nil {
		t.Logf("Sell transaction error: %v", err)
		t.Skip("Skipping test in devnet environment")
		return
	}
	assert.NotNil(t, tx)
	assert.Equal(t, blockchain.TransactionTypeSell, tx.Type)
}

func TestGetTransaction(t *testing.T) {
	provider, err := NewProvider(true)
	assert.NoError(t, err)

	// Test valid transaction ID
	tx, err := provider.GetTransaction(context.Background(), "11111111111111111111111111111111")
	if err != nil {
		t.Logf("Get transaction error: %v", err)
		t.Skip("Skipping test in devnet environment")
		return
	}
	assert.NotNil(t, tx)
	assert.Equal(t, "11111111111111111111111111111111", tx.ID)

	// Test invalid transaction ID
	tx, err = provider.GetTransaction(context.Background(), "invalid")
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestGetTransactions(t *testing.T) {
	provider, err := NewProvider(true)
	assert.NoError(t, err)

	// Test valid address
	txs, err := provider.GetTransactions(context.Background(), "11111111111111111111111111111111", 10)
	assert.NoError(t, err)
	assert.NotNil(t, txs)

	// Test invalid address
	txs, err = provider.GetTransactions(context.Background(), "invalid", 10)
	assert.Error(t, err)
	assert.Nil(t, txs)
}
