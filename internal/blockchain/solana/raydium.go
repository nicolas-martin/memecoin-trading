package solana

import (
	"context"
	"fmt"
	"meme-trader/internal/blockchain"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

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

	tx := solana.NewTransaction(
		[]solana.Instruction{swapInstruction},
		recentBlockhash.Value.Blockhash,
		solana.TransactionPayer(fromPubKey),
	)

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
		createIx := token.NewCreateAssociatedTokenAccountInstruction(
			owner,
			owner,
			tokenMint,
			tokenAccount,
		).Build()

		recentBlockhash, err := c.rpcClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
		if err != nil {
			return solana.PublicKey{}, fmt.Errorf("failed to get recent blockhash: %w", err)
		}

		tx := solana.NewTransaction(
			[]solana.Instruction{createIx},
			recentBlockhash.Value.Blockhash,
			solana.TransactionPayer(owner),
		)

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
