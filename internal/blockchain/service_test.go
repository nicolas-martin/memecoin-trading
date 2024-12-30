package blockchain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProvider is a mock implementation of the Provider interface
type MockProvider struct {
	mock.Mock
}

func (m *MockProvider) Network() Network {
	args := m.Called()
	return args.Get(0).(Network)
}

func (m *MockProvider) IsValidAddress(address string) bool {
	args := m.Called(address)
	return args.Bool(0)
}

func (m *MockProvider) CreateWallet(ctx context.Context) (*Wallet, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Wallet), args.Error(1)
}

func (m *MockProvider) GetWallet(ctx context.Context, address string) (*Wallet, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Wallet), args.Error(1)
}

func (m *MockProvider) GetBalance(ctx context.Context, address string) (Amount, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(Amount), args.Error(1)
}

func (m *MockProvider) Buy(ctx context.Context, req BuyRequest) (*Transaction, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Transaction), args.Error(1)
}

func (m *MockProvider) Sell(ctx context.Context, req SellRequest) (*Transaction, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Transaction), args.Error(1)
}

func (m *MockProvider) GetTransaction(ctx context.Context, txID string) (*Transaction, error) {
	args := m.Called(ctx, txID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Transaction), args.Error(1)
}

func (m *MockProvider) GetTransactions(ctx context.Context, address string, limit int) ([]Transaction, error) {
	args := m.Called(ctx, address, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Transaction), args.Error(1)
}

func (m *MockProvider) GetTopMemeCoins(ctx context.Context, req TopMemeCoinsRequest) ([]MemeCoin, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]MemeCoin), args.Error(1)
}

func TestNewService(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service, "Service should not be nil")
}

func TestRegisterProvider(t *testing.T) {
	service := NewService()
	mockProvider := new(MockProvider)
	mockProvider.On("Network").Return(NetworkSolana)

	err := service.RegisterProvider(mockProvider)
	assert.NoError(t, err, "Should register provider without error")

	// Register another provider with different priority
	mockProvider2 := new(MockProvider)
	mockProvider2.On("Network").Return(NetworkSolana)
	config := ProviderConfig{
		Priority:           1,
		RequestsPerWindow:  200,
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 5,
	}
	err = service.RegisterProviderWithConfig(mockProvider2, config)
	assert.NoError(t, err, "Should register second provider without error")
}

func TestCreateWallet(t *testing.T) {
	service := NewService()
	mockProvider := new(MockProvider)
	mockProvider.On("Network").Return(NetworkSolana)

	now := time.Now().Unix()
	expectedWallet := &Wallet{
		ID:            "test-wallet",
		Network:       NetworkSolana,
		Address:       "test-address",
		PublicKey:     "test-pubkey",
		PrivateKey:    "test-privkey",
		CreatedAt:     now,
		LastUpdatedAt: now,
	}

	mockProvider.On("CreateWallet", mock.Anything).Return(expectedWallet, nil)

	// Register provider
	err := service.RegisterProvider(mockProvider)
	assert.NoError(t, err, "Should register provider without error")

	// Create wallet
	wallet, err := service.CreateWallet(context.Background(), NetworkSolana)
	assert.NoError(t, err, "Should create wallet without error")
	assert.Equal(t, expectedWallet, wallet, "Should return expected wallet")
}

func TestBuyTransaction(t *testing.T) {
	service := NewService()
	mockProvider := new(MockProvider)
	mockProvider.On("Network").Return(NetworkSolana)

	expectedTx := &Transaction{
		ID:            "test-tx",
		Network:       NetworkSolana,
		Type:          TransactionTypeBuy,
		Status:        TransactionStatusPending,
		FromAddress:   "test-from",
		ToAddress:     "test-to",
		TokenAddress:  "test-token",
		CreatedAt:     time.Now().Unix(),
		LastUpdatedAt: time.Now().Unix(),
	}

	buyReq := BuyRequest{
		WalletAddress: "test-from",
		TokenAddress:  "test-token",
		Amount:        Amount{},
		MaxPrice:      Amount{},
	}

	mockProvider.On("Buy", mock.Anything, buyReq).Return(expectedTx, nil)

	// Register provider
	err := service.RegisterProvider(mockProvider)
	assert.NoError(t, err, "Should register provider without error")

	// Execute buy transaction
	tx, err := service.Buy(context.Background(), NetworkSolana, buyReq)
	assert.NoError(t, err, "Should execute buy without error")
	assert.Equal(t, expectedTx, tx, "Should return expected transaction")
}

func TestSellTransaction(t *testing.T) {
	service := NewService()
	mockProvider := new(MockProvider)
	mockProvider.On("Network").Return(NetworkSolana)

	expectedTx := &Transaction{
		ID:            "test-tx",
		Network:       NetworkSolana,
		Type:          TransactionTypeSell,
		Status:        TransactionStatusPending,
		FromAddress:   "test-from",
		ToAddress:     "test-to",
		TokenAddress:  "test-token",
		CreatedAt:     time.Now().Unix(),
		LastUpdatedAt: time.Now().Unix(),
	}

	sellReq := SellRequest{
		WalletAddress: "test-from",
		TokenAddress:  "test-token",
		Amount:        Amount{},
		MinPrice:      Amount{},
	}

	mockProvider.On("Sell", mock.Anything, sellReq).Return(expectedTx, nil)

	// Register provider
	err := service.RegisterProvider(mockProvider)
	assert.NoError(t, err, "Should register provider without error")

	// Execute sell transaction
	tx, err := service.Sell(context.Background(), NetworkSolana, sellReq)
	assert.NoError(t, err, "Should execute sell without error")
	assert.Equal(t, expectedTx, tx, "Should return expected transaction")
}

func TestProviderFallback(t *testing.T) {
	service := NewService()

	// Create two providers with different priorities
	mockProvider1 := new(MockProvider)
	mockProvider1.On("Network").Return(NetworkSolana)
	mockProvider2 := new(MockProvider)
	mockProvider2.On("Network").Return(NetworkSolana)

	// Register providers with different priorities
	err := service.RegisterProviderWithConfig(mockProvider1, ProviderConfig{
		Priority:           1,
		RequestsPerWindow:  100,
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 3,
	})
	assert.NoError(t, err)

	err = service.RegisterProviderWithConfig(mockProvider2, ProviderConfig{
		Priority:           2,
		RequestsPerWindow:  100,
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 3,
	})
	assert.NoError(t, err)

	// Set up mock responses
	expectedWallet := &Wallet{
		ID:      "test-wallet",
		Network: NetworkSolana,
		Address: "test-address",
	}

	// First provider fails
	mockProvider1.On("CreateWallet", mock.Anything).Return(nil, assert.AnError)
	// Second provider succeeds
	mockProvider2.On("CreateWallet", mock.Anything).Return(expectedWallet, nil)

	// Test fallback
	wallet, err := service.CreateWallet(context.Background(), NetworkSolana)
	assert.NoError(t, err)
	assert.Equal(t, expectedWallet, wallet)

	mockProvider1.AssertExpectations(t)
	mockProvider2.AssertExpectations(t)
}
