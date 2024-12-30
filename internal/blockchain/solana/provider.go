package solana

import (
	"context"
	"fmt"
	"math/big"
	"meme-trader/internal/blockchain"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

const (
	// Solana mainnet RPC endpoints (use your preferred endpoint)
	mainnetRPCEndpoint = "https://api.mainnet-beta.solana.com"
	mainnetWSEndpoint  = "wss://api.mainnet-beta.solana.com"

	// Solana devnet RPC endpoints (for testing)
	devnetRPCEndpoint = "https://api.devnet.solana.com"
	devnetWSEndpoint  = "wss://api.devnet.solana.com"

	// Default decimals for SOL
	solDecimals = 9
)

type Provider struct {
	rpcClient     *rpc.Client
	wsClient      *ws.Client
	raydiumClient *RaydiumClient
	network       blockchain.Network
	isDevnet      bool
}

// NewProvider creates a new Solana provider
func NewProvider(isDevnet bool) (*Provider, error) {
	rpcEndpoint := mainnetRPCEndpoint
	wsEndpoint := mainnetWSEndpoint
	if isDevnet {
		rpcEndpoint = devnetRPCEndpoint
		wsEndpoint = devnetWSEndpoint
	}

	rpcClient := rpc.New(rpcEndpoint)
	wsClient, err := ws.Connect(context.Background(), wsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %w", err)
	}

	raydiumClient := NewRaydiumClient(rpcClient, isDevnet)

	return &Provider{
		rpcClient:     rpcClient,
		wsClient:      wsClient,
		raydiumClient: raydiumClient,
		network:       blockchain.NetworkSolana,
		isDevnet:      isDevnet,
	}, nil
}

func (p *Provider) Network() blockchain.Network {
	return p.network
}

func (p *Provider) IsValidAddress(address string) bool {
	_, err := solana.PublicKeyFromBase58(address)
	return err == nil
}

func (p *Provider) CreateWallet(ctx context.Context) (*blockchain.Wallet, error) {
	// Generate a new keypair
	account := solana.NewWallet()
	privateKey := account.PrivateKey
	publicKey := account.PublicKey()

	now := time.Now().Unix()
	wallet := &blockchain.Wallet{
		ID:            publicKey.String(),
		Network:       p.network,
		Address:       publicKey.String(),
		PublicKey:     publicKey.String(),
		PrivateKey:    privateKey.String(),
		CreatedAt:     now,
		LastUpdatedAt: now,
	}

	return wallet, nil
}

func (p *Provider) GetWallet(ctx context.Context, address string) (*blockchain.Wallet, error) {
	if !p.IsValidAddress(address) {
		return nil, fmt.Errorf("invalid Solana address: %s", address)
	}

	// In a real implementation, you would fetch this from your database
	// Here we just return a basic wallet structure
	return &blockchain.Wallet{
		ID:            address,
		Network:       p.network,
		Address:       address,
		PublicKey:     address,
		CreatedAt:     time.Now().Unix(),
		LastUpdatedAt: time.Now().Unix(),
	}, nil
}

func (p *Provider) GetBalance(ctx context.Context, address string) (blockchain.Amount, error) {
	pubKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return blockchain.Amount{}, fmt.Errorf("invalid address: %w", err)
	}

	balance, err := p.rpcClient.GetBalance(
		ctx,
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return blockchain.Amount{}, fmt.Errorf("failed to get balance: %w", err)
	}

	return blockchain.Amount{
		Value:    new(big.Int).SetUint64(balance),
		Decimals: solDecimals,
	}, nil
}

func (p *Provider) Buy(ctx context.Context, req blockchain.BuyRequest) (*blockchain.Transaction, error) {
	if !p.IsValidAddress(req.WalletAddress) {
		return nil, fmt.Errorf("invalid wallet address")
	}
	if !p.IsValidAddress(req.TokenAddress) {
		return nil, fmt.Errorf("invalid token address")
	}

	// Use Raydium to execute the swap
	swapReq := SwapRequest{
		FromAddress:      req.WalletAddress,
		ToAddress:        req.TokenAddress,
		TokenAddress:     req.TokenAddress,
		Amount:           req.Amount,
		MinimumAmountOut: req.MaxPrice,
		Type:             blockchain.TransactionTypeBuy,
		Timestamp:        time.Now().Unix(),
	}

	return p.raydiumClient.SwapTokens(ctx, swapReq)
}

func (p *Provider) Sell(ctx context.Context, req blockchain.SellRequest) (*blockchain.Transaction, error) {
	if !p.IsValidAddress(req.WalletAddress) {
		return nil, fmt.Errorf("invalid wallet address")
	}
	if !p.IsValidAddress(req.TokenAddress) {
		return nil, fmt.Errorf("invalid token address")
	}

	// Use Raydium to execute the swap
	swapReq := SwapRequest{
		FromAddress:      req.WalletAddress,
		ToAddress:        req.TokenAddress,
		TokenAddress:     req.TokenAddress,
		Amount:           req.Amount,
		MinimumAmountOut: req.MinPrice,
		Type:             blockchain.TransactionTypeSell,
		Timestamp:        time.Now().Unix(),
	}

	return p.raydiumClient.SwapTokens(ctx, swapReq)
}

func (p *Provider) GetTransaction(ctx context.Context, txID string) (*blockchain.Transaction, error) {
	signature, err := solana.SignatureFromBase58(txID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	tx, err := p.rpcClient.GetTransaction(ctx, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Convert Solana transaction to our generic Transaction type
	// This is a simplified conversion
	status := blockchain.TransactionStatusConfirmed
	if !tx.Meta.Err.IsNil() {
		status = blockchain.TransactionStatusFailed
	}

	return &blockchain.Transaction{
		ID:            txID,
		Network:       p.network,
		Status:        status,
		BlockHash:     tx.BlockHash.String(),
		BlockNumber:   uint64(tx.Slot),
		Signature:     txID,
		Timestamp:     time.Now().Unix(), // You should get this from the block
		CreatedAt:     time.Now().Unix(),
		LastUpdatedAt: time.Now().Unix(),
	}, nil
}

func (p *Provider) GetTransactions(ctx context.Context, address string, limit int) ([]blockchain.Transaction, error) {
	pubKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	signatures, err := p.rpcClient.GetSignaturesForAddress(ctx, pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get signatures: %w", err)
	}

	var transactions []blockchain.Transaction
	for i, sig := range signatures {
		if i >= limit {
			break
		}

		tx, err := p.GetTransaction(ctx, sig.Signature.String())
		if err != nil {
			continue // Skip failed transactions
		}
		transactions = append(transactions, *tx)
	}

	return transactions, nil
}
