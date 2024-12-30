package solana

import (
	"context"
	"encoding/json"
	"math/big"
	"meme-trader/internal/blockchain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newBigInt(val uint64) *big.Int {
	return new(big.Int).SetUint64(val)
}

func TestRaydiumClient_GetTopMemeCoins(t *testing.T) {
	// Create a mock HTTP server for Raydium API
	raydiumServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pools := []RaydiumPool{
			{
				TokenAddress:   "token1",
				Symbol:         "MEME1",
				Name:           "MemeCoin1",
				LogoURL:        "https://example.com/meme1.png",
				Price:          1.0,
				MarketCap:      1000000000.0,
				Volume24h:      500000.0,
				PriceChange24h: 10.5,
				LastUpdated:    time.Now().Unix(),
				IsMemeCoin:     true,
			},
			{
				TokenAddress:   "token2",
				Symbol:         "MEME2",
				Name:           "MemeCoin2",
				LogoURL:        "https://example.com/meme2.png",
				Price:          2.0,
				MarketCap:      2000000000.0,
				Volume24h:      1000000.0,
				PriceChange24h: -5.2,
				LastUpdated:    time.Now().Unix(),
				IsMemeCoin:     true,
			},
		}
		json.NewEncoder(w).Encode(pools)
	}))
	defer raydiumServer.Close()

	// Create a mock HTTP server for Jupiter API
	jupiterServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "token.jup.ag") {
			tokens := []struct {
				Address    string                 `json:"address"`
				Symbol     string                 `json:"symbol"`
				Name       string                 `json:"name"`
				LogoURI    string                 `json:"logoURI"`
				Tags       []string               `json:"tags"`
				Extensions map[string]interface{} `json:"extensions"`
			}{
				{
					Address: "token1",
					Symbol:  "MEME1",
					Name:    "MemeCoin1",
					LogoURI: "https://example.com/meme1.png",
					Tags:    []string{"meme"},
				},
			}
			json.NewEncoder(w).Encode(tokens)
		} else {
			// Solana token list response
			response := struct {
				Tokens []struct {
					Address    string                 `json:"address"`
					Symbol     string                 `json:"symbol"`
					Name       string                 `json:"name"`
					LogoURI    string                 `json:"logoURI"`
					Tags       []string               `json:"tags"`
					Extensions map[string]interface{} `json:"extensions"`
				} `json:"tokens"`
			}{
				Tokens: []struct {
					Address    string                 `json:"address"`
					Symbol     string                 `json:"symbol"`
					Name       string                 `json:"name"`
					LogoURI    string                 `json:"logoURI"`
					Tags       []string               `json:"tags"`
					Extensions map[string]interface{} `json:"extensions"`
				}{
					{
						Address: "token1",
						Symbol:  "MEME1",
						Name:    "MemeCoin1",
						LogoURI: "https://example.com/meme1.png",
						Tags:    []string{"meme"},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer jupiterServer.Close()

	// Create a mock RPC client
	rpcClient := rpc.New(raydiumServer.URL)

	// Create Raydium client
	client := NewRaydiumClient(rpcClient, false)

	// Test GetTopMemeCoins
	memeCoins, err := client.GetTopMemeCoins(context.Background(), RaydiumTopMemeCoinsRequest{
		Limit:     10,
		TimeFrame: 24 * time.Hour,
	})

	require.NoError(t, err)
	require.Len(t, memeCoins, 2)

	// Verify first meme coin
	assert.Equal(t, "token1", memeCoins[0].Address)
	assert.Equal(t, "MEME1", memeCoins[0].Symbol)
	assert.Equal(t, "MemeCoin1", memeCoins[0].Name)
	assert.Equal(t, "https://example.com/meme1.png", memeCoins[0].LogoURL)
	assert.Equal(t, int64(1e9), memeCoins[0].Price.Value.Int64()) // 1.0 SOL = 1e9 lamports
	assert.Equal(t, int64(1000000000), memeCoins[0].MarketCap.Value.Int64())
	assert.Equal(t, int64(500000), memeCoins[0].Volume24h.Value.Int64())
	assert.Equal(t, 10.5, memeCoins[0].Change24h)

	// Verify second meme coin
	assert.Equal(t, "token2", memeCoins[1].Address)
	assert.Equal(t, "MEME2", memeCoins[1].Symbol)
	assert.Equal(t, "MemeCoin2", memeCoins[1].Name)
	assert.Equal(t, "https://example.com/meme2.png", memeCoins[1].LogoURL)
	assert.Equal(t, int64(2e9), memeCoins[1].Price.Value.Int64()) // 2.0 SOL = 2e9 lamports
	assert.Equal(t, int64(2000000000), memeCoins[1].MarketCap.Value.Int64())
	assert.Equal(t, int64(1000000), memeCoins[1].Volume24h.Value.Int64())
	assert.Equal(t, -5.2, memeCoins[1].Change24h)
}

func TestRaydiumClient_SwapTokens(t *testing.T) {
	// Create a mock RPC client
	rpcClient := rpc.New("http://localhost:8899")

	// Create Raydium client
	client := NewRaydiumClient(rpcClient, true)

	// Test buildSwapInstruction
	fromAccount := solana.NewWallet().PublicKey()
	toAccount := solana.NewWallet().PublicKey()
	amount := blockchain.Amount{Value: newBigInt(1000000)}
	minimumAmountOut := blockchain.Amount{Value: newBigInt(900000)}

	instruction := client.buildSwapInstruction(fromAccount, toAccount, amount, minimumAmountOut)

	// Verify instruction
	assert.NotNil(t, instruction)
	assert.Equal(t, "DnXyn8dAR5fJdqfBQciQ6gPSDNMQSTkQrPsR65ZF5qoW", instruction.ProgramID().String())

	accounts := instruction.Accounts()
	require.Len(t, accounts, 2)
	assert.Equal(t, fromAccount, accounts[0].PublicKey)
	assert.Equal(t, toAccount, accounts[1].PublicKey)
	assert.True(t, accounts[0].IsWritable)
	assert.True(t, accounts[1].IsWritable)
	assert.False(t, accounts[0].IsSigner)
	assert.False(t, accounts[1].IsSigner)

	data, err := instruction.Data()
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Equal(t, byte(0x0), data[0]) // Verify instruction index
}

func TestRaydiumClient_GetTokenMetadata(t *testing.T) {
	// Create a mock HTTP server for Jupiter API
	jupiterServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "token.jup.ag") {
			response := []struct {
				Address    string                 `json:"address"`
				Symbol     string                 `json:"symbol"`
				Name       string                 `json:"name"`
				LogoURI    string                 `json:"logoURI"`
				Tags       []string               `json:"tags"`
				Extensions map[string]interface{} `json:"extensions"`
			}{
				{
					Address: "token1",
					Symbol:  "MEME1",
					Name:    "MemeCoin1",
					LogoURI: "https://example.com/meme1.png",
					Tags:    []string{"meme"},
				},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			// Solana token list response
			response := struct {
				Tokens []struct {
					Address    string                 `json:"address"`
					Symbol     string                 `json:"symbol"`
					Name       string                 `json:"name"`
					LogoURI    string                 `json:"logoURI"`
					Tags       []string               `json:"tags"`
					Extensions map[string]interface{} `json:"extensions"`
				} `json:"tokens"`
			}{
				Tokens: []struct {
					Address    string                 `json:"address"`
					Symbol     string                 `json:"symbol"`
					Name       string                 `json:"name"`
					LogoURI    string                 `json:"logoURI"`
					Tags       []string               `json:"tags"`
					Extensions map[string]interface{} `json:"extensions"`
				}{
					{
						Address: "token1",
						Symbol:  "MEME1",
						Name:    "MemeCoin1",
						LogoURI: "https://example.com/meme1.png",
						Tags:    []string{"meme"},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer jupiterServer.Close()

	// Create a mock RPC client
	rpcClient := rpc.New(jupiterServer.URL)

	// Create Raydium client
	client := NewRaydiumClient(rpcClient, false)

	// Test getTokenMetadata
	metadata, err := client.getTokenMetadata(context.Background(), "token1")
	require.NoError(t, err)
	require.NotNil(t, metadata)

	assert.Equal(t, "MEME1", metadata.Symbol)
	assert.Equal(t, "MemeCoin1", metadata.Name)
	assert.Equal(t, "https://example.com/meme1.png", metadata.LogoURL)
	assert.Contains(t, metadata.Tags, "meme")
}

func TestRaydiumClient_IsMemeCoin(t *testing.T) {
	client := NewRaydiumClient(nil, false)

	testCases := []struct {
		name     string
		metadata *TokenMetadata
		expected bool
	}{
		{
			name: "Meme tag",
			metadata: &TokenMetadata{
				Symbol: "TEST",
				Name:   "Test Token",
				Tags:   []string{"meme"},
			},
			expected: true,
		},
		{
			name: "Memecoin tag",
			metadata: &TokenMetadata{
				Symbol: "TEST",
				Name:   "Test Token",
				Tags:   []string{"memecoin"},
			},
			expected: true,
		},
		{
			name: "Doge in name",
			metadata: &TokenMetadata{
				Symbol: "TEST",
				Name:   "Doge Token",
				Tags:   []string{},
			},
			expected: true,
		},
		{
			name: "Pepe in symbol",
			metadata: &TokenMetadata{
				Symbol: "PEPE",
				Name:   "Test Token",
				Tags:   []string{},
			},
			expected: true,
		},
		{
			name: "Not a meme coin",
			metadata: &TokenMetadata{
				Symbol: "TEST",
				Name:   "Test Token",
				Tags:   []string{"defi"},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := client.isMemeCoin(tc.metadata)
			assert.Equal(t, tc.expected, result)
		})
	}
}
