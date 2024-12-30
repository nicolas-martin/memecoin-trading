package solana

import (
	"context"
	"fmt"
	"math/big"
	"meme-trader/internal/blockchain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestRaydiumClient is a mock implementation of the Raydium client for testing
type TestRaydiumClient struct {
	mock.Mock
}

func (m *TestRaydiumClient) SwapTokens(ctx context.Context, req SwapRequest) (*blockchain.Transaction, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*blockchain.Transaction), args.Error(1)
}

func TestSwapTokens(t *testing.T) {
	// Create mock Raydium client
	mockClient := new(TestRaydiumClient)

	// Set up test data
	fromAddress := "11111111111111111111111111111111"
	toAddress := "22222222222222222222222222222222"
	tokenAddress := "33333333333333333333333333333333"
	amount := blockchain.Amount{
		Value:    big.NewInt(1000000),
		Decimals: 9,
	}
	minAmountOut := blockchain.Amount{
		Value:    big.NewInt(900000),
		Decimals: 9,
	}

	// Create swap request
	req := SwapRequest{
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		TokenAddress:     tokenAddress,
		Amount:           amount,
		MinimumAmountOut: minAmountOut,
		Type:             blockchain.TransactionTypeBuy,
		Timestamp:        time.Now().Unix(),
	}

	// Set up mock expectations
	expectedTx := &blockchain.Transaction{
		ID:            "test-tx",
		Network:       blockchain.NetworkSolana,
		Type:          blockchain.TransactionTypeBuy,
		Status:        blockchain.TransactionStatusPending,
		FromAddress:   fromAddress,
		ToAddress:     toAddress,
		TokenAddress:  tokenAddress,
		Amount:        amount,
		CreatedAt:     req.Timestamp,
		LastUpdatedAt: req.Timestamp,
	}

	mockClient.On("SwapTokens", mock.Anything, req).Return(expectedTx, nil)

	// Execute swap
	tx, err := mockClient.SwapTokens(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, blockchain.TransactionStatusPending, tx.Status)
	assert.Equal(t, expectedTx, tx)
}

func TestSwapTokensError(t *testing.T) {
	// Create mock Raydium client
	mockClient := new(TestRaydiumClient)

	// Set up test data
	fromAddress := "11111111111111111111111111111111"
	toAddress := "22222222222222222222222222222222"
	tokenAddress := "33333333333333333333333333333333"
	amount := blockchain.Amount{
		Value:    big.NewInt(1000000),
		Decimals: 9,
	}
	minAmountOut := blockchain.Amount{
		Value:    big.NewInt(900000),
		Decimals: 9,
	}

	// Create swap request
	req := SwapRequest{
		FromAddress:      fromAddress,
		ToAddress:        toAddress,
		TokenAddress:     tokenAddress,
		Amount:           amount,
		MinimumAmountOut: minAmountOut,
		Type:             blockchain.TransactionTypeBuy,
		Timestamp:        time.Now().Unix(),
	}

	// Set up mock expectations with error
	mockClient.On("SwapTokens", mock.Anything, req).Return(nil, fmt.Errorf("failed to execute swap"))

	// Execute swap
	tx, err := mockClient.SwapTokens(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, tx)
	assert.Contains(t, err.Error(), "failed to execute swap")
}
