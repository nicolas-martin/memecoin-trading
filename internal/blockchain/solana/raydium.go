package solana

import (
	"context"
	"fmt"
	"log"
	"meme-trader/internal/blockchain"
	"sort"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/rpc"
)

// RaydiumTopMemeCoinsRequest represents parameters for fetching top meme coins from Raydium
type RaydiumTopMemeCoinsRequest struct {
	Limit     int
	TimeFrame time.Duration
}

// RaydiumClient handles interactions with the Raydium DEX
type RaydiumClient struct {
	rpcClient *rpc.Client
	isDevnet  bool
}

// NewRaydiumClient creates a new Raydium client
func NewRaydiumClient(rpcClient *rpc.Client, isDevnet bool) *RaydiumClient {
	return &RaydiumClient{
		rpcClient: rpcClient,
		isDevnet:  isDevnet,
	}
}

// GetTopMemeCoins fetches the top meme coins from Raydium DEX
func (c *RaydiumClient) GetTopMemeCoins(ctx context.Context, req RaydiumTopMemeCoinsRequest) ([]blockchain.MemeCoin, error) {
	// Get all pools from Raydium
	pools, err := c.getAllPools(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pools: %w", err)
	}

	// Filter and sort pools by volume
	memeCoins := c.filterAndSortMemePools(pools, req.TimeFrame)

	// Limit the results
	if len(memeCoins) > req.Limit {
		memeCoins = memeCoins[:req.Limit]
	}

	return memeCoins, nil
}

// getAllPools fetches all liquidity pools from Raydium
func (c *RaydiumClient) getAllPools(ctx context.Context) ([]RaydiumPool, error) {
	// Fetch pools from Raydium API
	pools, err := c.fetchRaydiumPools(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Raydium pools: %w", err)
	}

	// Enrich pool data with token metadata
	enrichedPools := make([]RaydiumPool, 0, len(pools))
	for _, pool := range pools {
		// Get token metadata from Jupiter API or token list
		metadata, err := c.getTokenMetadata(ctx, pool.TokenAddress)
		if err != nil {
			continue // Skip tokens we can't get metadata for
		}

		// Check if it's a meme coin based on metadata
		isMemeCoin := c.isMemeCoin(metadata)
		if !isMemeCoin {
			continue
		}

		enrichedPool := RaydiumPool{
			TokenAddress:   pool.TokenAddress,
			Symbol:         metadata.Symbol,
			Name:           metadata.Name,
			LogoURL:        metadata.LogoURL,
			Price:          pool.Price,
			MarketCap:      pool.MarketCap,
			Volume24h:      pool.Volume24h,
			PriceChange24h: pool.PriceChange24h,
			LastUpdated:    time.Now(),
			IsMemeCoin:     true,
		}

		enrichedPools = append(enrichedPools, enrichedPool)
	}

	return enrichedPools, nil
}

// TokenMetadata represents metadata for a token
type TokenMetadata struct {
	Symbol     string
	Name       string
	LogoURL    string
	Tags       []string
	Extensions map[string]interface{}
}

// fetchRaydiumPools fetches raw pool data from Raydium
func (c *RaydiumClient) fetchRaydiumPools(ctx context.Context) ([]RaydiumPool, error) {
	// TODO: Implement actual Raydium API call
	// For now, return test data
	return []RaydiumPool{}, nil
}

// getTokenMetadata fetches token metadata from Jupiter API or token list
func (c *RaydiumClient) getTokenMetadata(ctx context.Context, tokenAddress string) (*TokenMetadata, error) {
	// First try Jupiter API
	metadata, err := c.getJupiterTokenMetadata(ctx, tokenAddress)
	if err == nil {
		return metadata, nil
	}

	// Fallback to token list
	return c.getTokenListMetadata(ctx, tokenAddress)
}

// getJupiterTokenMetadata fetches token metadata from Jupiter API
func (c *RaydiumClient) getJupiterTokenMetadata(ctx context.Context, tokenAddress string) (*TokenMetadata, error) {
	jupiterEndpoint := "https://token.jup.ag/all"

	// Log the endpoint we're using
	log.Printf("Fetching token metadata from Jupiter API: %s", jupiterEndpoint)

	// TODO: Implement Jupiter API call
	return nil, fmt.Errorf("not implemented: %s", jupiterEndpoint)
}

// getTokenListMetadata fetches token metadata from Solana token list
func (c *RaydiumClient) getTokenListMetadata(ctx context.Context, tokenAddress string) (*TokenMetadata, error) {
	tokenListEndpoint := "https://cdn.jsdelivr.net/gh/solana-labs/token-list@main/src/tokens/solana.tokenlist.json"

	// Log the endpoint we're using
	log.Printf("Fetching token metadata from Solana token list: %s", tokenListEndpoint)

	// TODO: Implement token list fetch
	return nil, fmt.Errorf("not implemented: %s", tokenListEndpoint)
}

// isMemeCoin determines if a token is a meme coin based on its metadata
func (c *RaydiumClient) isMemeCoin(metadata *TokenMetadata) bool {
	// Check tags
	for _, tag := range metadata.Tags {
		if tag == "meme" || tag == "memecoin" {
			return true
		}
	}

	// Check name/symbol for common meme coin indicators
	memeKeywords := []string{
		"doge", "shib", "inu", "pepe", "wojak", "moon", "elon",
		"safe", "baby", "rocket", "chad", "based", "wagmi", "frog",
	}

	nameLower := strings.ToLower(metadata.Name)
	symbolLower := strings.ToLower(metadata.Symbol)

	for _, keyword := range memeKeywords {
		if strings.Contains(nameLower, keyword) || strings.Contains(symbolLower, keyword) {
			return true
		}
	}

	return false
}

// RaydiumPool represents a liquidity pool on Raydium
type RaydiumPool struct {
	TokenAddress   string
	Symbol         string
	Name           string
	LogoURL        string
	Price          blockchain.Amount
	MarketCap      blockchain.Amount
	Volume24h      blockchain.Amount
	PriceChange24h float64
	LastUpdated    time.Time
	IsMemeCoin     bool
}

// filterAndSortMemePools filters out non-meme coins and sorts by volume
func (c *RaydiumClient) filterAndSortMemePools(pools []RaydiumPool, timeFrame time.Duration) []blockchain.MemeCoin {
	var memeCoins []blockchain.MemeCoin

	// Filter meme coins
	for _, pool := range pools {
		if !pool.IsMemeCoin {
			continue
		}

		// Convert pool to MemeCoin
		memeCoin := blockchain.MemeCoin{
			Address:     pool.TokenAddress,
			Symbol:      pool.Symbol,
			Name:        pool.Name,
			LogoURL:     pool.LogoURL,
			Price:       pool.Price,
			MarketCap:   pool.MarketCap,
			Volume24h:   pool.Volume24h,
			Change24h:   pool.PriceChange24h,
			LastUpdated: pool.LastUpdated,
		}

		memeCoins = append(memeCoins, memeCoin)
	}

	// Sort by volume (descending)
	sort.Slice(memeCoins, func(i, j int) bool {
		return memeCoins[i].Volume24h.Value.Cmp(memeCoins[j].Volume24h.Value) > 0
	})

	return memeCoins
}

// SwapTokens executes a token swap on Raydium
func (c *RaydiumClient) SwapTokens(ctx context.Context, req SwapRequest) (*blockchain.Transaction, error) {
	// Validate input
	fromPubKey, err := solana.PublicKeyFromBase58(req.FromAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid from address: %w", err)
	}

	toPubKey, err := solana.PublicKeyFromBase58(req.ToAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid to address: %w", err)
	}

	// Get token accounts
	fromTokenAccount, err := c.getTokenAccount(ctx, fromPubKey, req.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get from token account: %w", err)
	}

	toTokenAccount, err := c.getTokenAccount(ctx, toPubKey, req.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get to token account: %w", err)
	}

	// Build the swap instruction
	swapInstruction := c.buildSwapInstruction(
		fromTokenAccount,
		toTokenAccount,
		req.Amount,
		req.MinimumAmountOut,
	)

	// Build the transaction
	recentBlockhash, err := c.rpcClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{swapInstruction},
		recentBlockhash.Value.Blockhash,
		solana.TransactionPayer(fromPubKey),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Sign and send the transaction
	sig, err := c.rpcClient.SendTransaction(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	// Create transaction record
	transaction := &blockchain.Transaction{
		ID:            sig.String(),
		Network:       blockchain.NetworkSolana,
		Type:          req.Type,
		Status:        blockchain.TransactionStatusPending,
		FromAddress:   req.FromAddress,
		ToAddress:     req.ToAddress,
		Amount:        req.Amount,
		TokenAddress:  req.TokenAddress,
		Signature:     sig.String(),
		BlockHash:     recentBlockhash.Value.Blockhash.String(),
		CreatedAt:     req.Timestamp,
		LastUpdatedAt: req.Timestamp,
	}

	return transaction, nil
}

// SwapRequest represents a request to swap tokens on Raydium
type SwapRequest struct {
	FromAddress      string
	ToAddress        string
	TokenAddress     string
	Amount           blockchain.Amount
	MinimumAmountOut blockchain.Amount
	Type             blockchain.TransactionType
	Timestamp        int64
}

// getTokenAccount gets or creates a token account for a given token
func (c *RaydiumClient) getTokenAccount(ctx context.Context, owner solana.PublicKey, tokenAddress string) (solana.PublicKey, error) {
	tokenMint, err := solana.PublicKeyFromBase58(tokenAddress)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("invalid token address: %w", err)
	}

	// Find associated token account
	tokenAccount, _, err := solana.FindAssociatedTokenAddress(
		owner,
		tokenMint,
	)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to find token account: %w", err)
	}

	// Check if account exists
	_, err = c.rpcClient.GetAccountInfo(ctx, tokenAccount)
	if err != nil {
		// Create account if it doesn't exist
		createIx := ata.NewCreateInstruction(
			owner,
			owner,
			tokenMint,
		).Build()

		recentBlockhash, err := c.rpcClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
		if err != nil {
			return solana.PublicKey{}, fmt.Errorf("failed to get recent blockhash: %w", err)
		}

		tx, err := solana.NewTransaction(
			[]solana.Instruction{createIx},
			recentBlockhash.Value.Blockhash,
			solana.TransactionPayer(owner),
		)
		if err != nil {
			return solana.PublicKey{}, fmt.Errorf("failed to create transaction: %w", err)
		}

		_, err = c.rpcClient.SendTransaction(ctx, tx)
		if err != nil {
			return solana.PublicKey{}, fmt.Errorf("failed to create token account: %w", err)
		}
	}

	return tokenAccount, nil
}

// buildSwapInstruction builds the Raydium swap instruction
func (c *RaydiumClient) buildSwapInstruction(
	fromAccount solana.PublicKey,
	toAccount solana.PublicKey,
	amount blockchain.Amount,
	minimumAmountOut blockchain.Amount,
) solana.Instruction {
	// This is a placeholder for the actual Raydium swap instruction
	// You would need to implement this based on Raydium's smart contract
	// and instruction format

	// Example structure:
	/*
		return &raydium.SwapInstruction{
			TokenProgramID: token.ProgramID,
			FromAccount:    fromAccount,
			ToAccount:      toAccount,
			Amount:         amount.Value.Uint64(),
			MinimumOut:     minimumAmountOut.Value.Uint64(),
		}
	*/

	// For now, return a dummy instruction
	return &solana.GenericInstruction{}
}
