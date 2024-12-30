package providers

import (
	"context"
	"testing"
	"time"

	"meme-trader/internal/blockchain"
	"meme-trader/internal/blockchain/solana"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderIntegration(t *testing.T) {
	// Skip in short mode as these are integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create blockchain service
	service := blockchain.NewService()

	// Create and register Solana provider
	provider, err := solana.NewProvider(true) // Use devnet for testing
	require.NoError(t, err, "Failed to create Solana provider")
	err = service.RegisterProvider(provider)
	require.NoError(t, err, "Failed to register Solana provider")

	// Test cases with real token addresses
	testCases := []struct {
		name          string
		chainID       blockchain.Network
		tokenAddress  string
		expectedError bool
	}{
		{
			name:         "Fetch BONK token data",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263",
		},
		{
			name:         "Fetch WIF token data",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm",
		},
		{
			name:         "Fetch MYRO token data",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "HhJpBhRRn4g56VsyLuT8DL5Bv31HkXqsrahTTUCZeZg4",
		},
		{
			name:         "Fetch PEPE token data",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "E4vX7kEegE3qVHUdk9fRAXwK1ZNWkL6YVxJVsHPNcvVM",
		},
		{
			name:          "Invalid token address",
			chainID:       blockchain.NetworkSolana,
			tokenAddress:  "invalid",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Get wallet data
			wallet, err := service.GetWallet(ctx, tc.chainID, tc.tokenAddress)

			if tc.expectedError {
				assert.Error(t, err, "Expected error for invalid token")
				return
			}

			require.NoError(t, err, "Failed to get wallet data")
			assert.NotNil(t, wallet, "Wallet data should not be nil")
			assert.Equal(t, tc.chainID, wallet.Network, "Network mismatch")
			assert.Equal(t, tc.tokenAddress, wallet.Address, "Address mismatch")
			assert.NotEmpty(t, wallet.PublicKey, "Public key should not be empty")
		})
	}
}

func TestProviderBalanceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping balance integration test in short mode")
	}

	// Create blockchain service
	service := blockchain.NewService()

	// Create and register Solana provider
	provider, err := solana.NewProvider(true) // Use devnet for testing
	require.NoError(t, err, "Failed to create Solana provider")
	err = service.RegisterProvider(provider)
	require.NoError(t, err, "Failed to register Solana provider")

	testCases := []struct {
		name          string
		chainID       blockchain.Network
		tokenAddress  string
		expectedError bool
	}{
		{
			name:         "Fetch BONK balance",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263",
		},
		{
			name:         "Fetch WIF balance",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm",
		},
		{
			name:         "Fetch MYRO balance",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "HhJpBhRRn4g56VsyLuT8DL5Bv31HkXqsrahTTUCZeZg4",
		},
		{
			name:         "Fetch PEPE balance",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "E4vX7kEegE3qVHUdk9fRAXwK1ZNWkL6YVxJVsHPNcvVM",
		},
		{
			name:          "Invalid token address",
			chainID:       blockchain.NetworkSolana,
			tokenAddress:  "invalid",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Get balance
			balance, err := service.GetBalance(ctx, tc.chainID, tc.tokenAddress)

			if tc.expectedError {
				assert.Error(t, err, "Expected error for invalid token")
				return
			}

			require.NoError(t, err, "Failed to get balance")
			assert.NotNil(t, balance, "Balance should not be nil")
			assert.Greater(t, balance.Decimals, uint8(0), "Decimals should be greater than 0")
		})
	}
}

func TestProviderTransactionsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping transactions integration test in short mode")
	}

	// Create blockchain service
	service := blockchain.NewService()

	// Create and register Solana provider
	provider, err := solana.NewProvider(true) // Use devnet for testing
	require.NoError(t, err, "Failed to create Solana provider")
	err = service.RegisterProvider(provider)
	require.NoError(t, err, "Failed to register Solana provider")

	testCases := []struct {
		name          string
		chainID       blockchain.Network
		tokenAddress  string
		expectedError bool
	}{
		{
			name:         "Fetch BONK transactions",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263",
		},
		{
			name:         "Fetch WIF transactions",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm",
		},
		{
			name:         "Fetch MYRO transactions",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "HhJpBhRRn4g56VsyLuT8DL5Bv31HkXqsrahTTUCZeZg4",
		},
		{
			name:         "Fetch PEPE transactions",
			chainID:      blockchain.NetworkSolana,
			tokenAddress: "E4vX7kEegE3qVHUdk9fRAXwK1ZNWkL6YVxJVsHPNcvVM",
		},
		{
			name:          "Invalid token address",
			chainID:       blockchain.NetworkSolana,
			tokenAddress:  "invalid",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Get transactions
			transactions, err := service.GetTransactions(ctx, tc.chainID, tc.tokenAddress, 10)

			if tc.expectedError {
				assert.Error(t, err, "Expected error for invalid token")
				return
			}

			require.NoError(t, err, "Failed to get transactions")
			if transactions != nil && len(transactions) > 0 {
				tx := transactions[0]
				assert.Equal(t, tc.chainID, tx.Network, "Network mismatch")
				assert.NotEmpty(t, tx.ID, "Transaction ID should not be empty")
				assert.NotEmpty(t, tx.Signature, "Signature should not be empty")
			} else {
				t.Logf("No transactions found for %s, this is acceptable for new or inactive tokens", tc.name)
			}
		})
	}
}

func TestProviderFallbackIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping fallback integration test in short mode")
	}

	// Create blockchain service
	service := blockchain.NewService()

	// Create and register two Solana providers with different priorities
	provider1, err := solana.NewProvider(true) // Use devnet for testing
	require.NoError(t, err, "Failed to create first Solana provider")
	provider2, err := solana.NewProvider(true) // Use devnet for testing
	require.NoError(t, err, "Failed to create second Solana provider")

	// Register providers with different priorities
	err = service.RegisterProviderWithConfig(provider1, blockchain.ProviderConfig{
		Priority:           1,
		RequestsPerWindow:  1, // Set low limit to trigger fallback
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 3,
	})
	require.NoError(t, err, "Failed to register first provider")

	err = service.RegisterProviderWithConfig(provider2, blockchain.ProviderConfig{
		Priority:           2,
		RequestsPerWindow:  100,
		WindowDuration:     time.Minute,
		HealthCheckPeriod:  time.Minute,
		MaxConsecutiveErrs: 3,
	})
	require.NoError(t, err, "Failed to register second provider")

	t.Run("Solana fallback test", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Make multiple requests to trigger fallback
		for i := 0; i < 3; i++ {
			// Get transactions
			transactions, err := service.GetTransactions(ctx, blockchain.NetworkSolana, "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263", 10)
			require.NoError(t, err, "Failed to get transactions after fallback")

			// Some tokens might not have transactions yet, so we don't assert NotNil
			if transactions != nil && len(transactions) > 0 {
				tx := transactions[0]
				assert.Equal(t, blockchain.NetworkSolana, tx.Network, "Network mismatch")
				assert.NotEmpty(t, tx.ID, "Transaction ID should not be empty")
				assert.NotEmpty(t, tx.Signature, "Signature should not be empty")
			} else {
				t.Log("No transactions found in fallback test, this is acceptable for new or inactive tokens")
			}

			// Add a small delay between requests
			time.Sleep(100 * time.Millisecond)
		}
	})
}
