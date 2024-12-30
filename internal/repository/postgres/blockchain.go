package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"meme-trader/internal/blockchain"
)

// CreateBlockchainTables creates the necessary tables for blockchain data
func createBlockchainTables(db *sql.DB) error {
	// Create wallets table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS wallets (
			id TEXT PRIMARY KEY,
			network TEXT NOT NULL,
			address TEXT NOT NULL,
			public_key TEXT NOT NULL,
			private_key TEXT NOT NULL,
			created_at BIGINT NOT NULL,
			last_updated_at BIGINT NOT NULL,
			UNIQUE(network, address)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create wallets table: %w", err)
	}

	// Create transactions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS blockchain_transactions (
			id TEXT PRIMARY KEY,
			network TEXT NOT NULL,
			type TEXT NOT NULL,
			status TEXT NOT NULL,
			from_address TEXT NOT NULL,
			to_address TEXT NOT NULL,
			amount_value TEXT NOT NULL,
			amount_decimals INTEGER NOT NULL,
			token_address TEXT NOT NULL,
			signature TEXT,
			block_hash TEXT,
			block_number BIGINT,
			timestamp BIGINT NOT NULL,
			gas_fee_value TEXT,
			gas_fee_decimals INTEGER,
			error_message TEXT,
			created_at BIGINT NOT NULL,
			last_updated_at BIGINT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create blockchain_transactions table: %w", err)
	}

	return nil
}

// SaveWallet saves a wallet to the database
func (db *Database) SaveWallet(wallet *blockchain.Wallet) error {
	_, err := db.db.Exec(`
		INSERT INTO wallets (
			id, network, address, public_key, private_key,
			created_at, last_updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (network, address) DO UPDATE SET
			public_key = EXCLUDED.public_key,
			private_key = EXCLUDED.private_key,
			last_updated_at = EXCLUDED.last_updated_at
	`,
		wallet.ID, wallet.Network, wallet.Address, wallet.PublicKey,
		wallet.PrivateKey, wallet.CreatedAt, wallet.LastUpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save wallet: %w", err)
	}

	return nil
}

// GetWallet retrieves a wallet by network and address
func (db *Database) GetWallet(network blockchain.Network, address string) (*blockchain.Wallet, error) {
	var wallet blockchain.Wallet
	err := db.db.QueryRow(`
		SELECT id, network, address, public_key, private_key,
			created_at, last_updated_at
		FROM wallets
		WHERE network = $1 AND address = $2
	`, network, address).Scan(
		&wallet.ID, &wallet.Network, &wallet.Address, &wallet.PublicKey,
		&wallet.PrivateKey, &wallet.CreatedAt, &wallet.LastUpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("wallet not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

// SaveTransaction saves a blockchain transaction to the database
func (db *Database) SaveTransaction(tx *blockchain.Transaction) error {
	amountValue, err := json.Marshal(tx.Amount.Value)
	if err != nil {
		return fmt.Errorf("failed to marshal amount value: %w", err)
	}

	var gasFeeValue []byte
	if tx.GasFee.Value != nil {
		gasFeeValue, err = json.Marshal(tx.GasFee.Value)
		if err != nil {
			return fmt.Errorf("failed to marshal gas fee value: %w", err)
		}
	}

	_, err = db.db.Exec(`
		INSERT INTO blockchain_transactions (
			id, network, type, status, from_address, to_address,
			amount_value, amount_decimals, token_address, signature,
			block_hash, block_number, timestamp, gas_fee_value,
			gas_fee_decimals, error_message, created_at, last_updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			signature = EXCLUDED.signature,
			block_hash = EXCLUDED.block_hash,
			block_number = EXCLUDED.block_number,
			gas_fee_value = EXCLUDED.gas_fee_value,
			gas_fee_decimals = EXCLUDED.gas_fee_decimals,
			error_message = EXCLUDED.error_message,
			last_updated_at = EXCLUDED.last_updated_at
	`,
		tx.ID, tx.Network, tx.Type, tx.Status, tx.FromAddress, tx.ToAddress,
		string(amountValue), tx.Amount.Decimals, tx.TokenAddress, tx.Signature,
		tx.BlockHash, tx.BlockNumber, tx.Timestamp,
		string(gasFeeValue), tx.GasFee.Decimals, tx.ErrorMessage,
		tx.CreatedAt, tx.LastUpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	return nil
}

// GetTransaction retrieves a blockchain transaction by ID
func (db *Database) GetTransaction(network blockchain.Network, txID string) (*blockchain.Transaction, error) {
	var tx blockchain.Transaction
	var amountValue, gasFeeValue string

	err := db.db.QueryRow(`
		SELECT id, network, type, status, from_address, to_address,
			amount_value, amount_decimals, token_address, signature,
			block_hash, block_number, timestamp, gas_fee_value,
			gas_fee_decimals, error_message, created_at, last_updated_at
		FROM blockchain_transactions
		WHERE network = $1 AND id = $2
	`, network, txID).Scan(
		&tx.ID, &tx.Network, &tx.Type, &tx.Status, &tx.FromAddress, &tx.ToAddress,
		&amountValue, &tx.Amount.Decimals, &tx.TokenAddress, &tx.Signature,
		&tx.BlockHash, &tx.BlockNumber, &tx.Timestamp,
		&gasFeeValue, &tx.GasFee.Decimals, &tx.ErrorMessage,
		&tx.CreatedAt, &tx.LastUpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transaction not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Parse amount value
	err = json.Unmarshal([]byte(amountValue), &tx.Amount.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal amount value: %w", err)
	}

	// Parse gas fee value if present
	if gasFeeValue != "" {
		err = json.Unmarshal([]byte(gasFeeValue), &tx.GasFee.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal gas fee value: %w", err)
		}
	}

	return &tx, nil
}

// GetTransactions retrieves blockchain transactions for a wallet
func (db *Database) GetTransactions(network blockchain.Network, address string, limit int) ([]blockchain.Transaction, error) {
	rows, err := db.db.Query(`
		SELECT id, network, type, status, from_address, to_address,
			amount_value, amount_decimals, token_address, signature,
			block_hash, block_number, timestamp, gas_fee_value,
			gas_fee_decimals, error_message, created_at, last_updated_at
		FROM blockchain_transactions
		WHERE network = $1 AND (from_address = $2 OR to_address = $2)
		ORDER BY timestamp DESC
		LIMIT $3
	`, network, address, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []blockchain.Transaction
	for rows.Next() {
		var tx blockchain.Transaction
		var amountValue, gasFeeValue string

		err := rows.Scan(
			&tx.ID, &tx.Network, &tx.Type, &tx.Status, &tx.FromAddress, &tx.ToAddress,
			&amountValue, &tx.Amount.Decimals, &tx.TokenAddress, &tx.Signature,
			&tx.BlockHash, &tx.BlockNumber, &tx.Timestamp,
			&gasFeeValue, &tx.GasFee.Decimals, &tx.ErrorMessage,
			&tx.CreatedAt, &tx.LastUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Parse amount value
		err = json.Unmarshal([]byte(amountValue), &tx.Amount.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal amount value: %w", err)
		}

		// Parse gas fee value if present
		if gasFeeValue != "" {
			err = json.Unmarshal([]byte(gasFeeValue), &tx.GasFee.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal gas fee value: %w", err)
			}
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
