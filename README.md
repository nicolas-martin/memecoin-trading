# Meme Trader

A platform for trading meme coins on various blockchains, starting with Solana.

## Features

- Multi-blockchain support (currently Solana)
- Wallet creation and management
- Buy/sell transactions using Raydium DEX
- Transaction history tracking
- Real-time price updates
- Mobile app integration

## Architecture

The platform uses a modular architecture that allows for easy integration of additional blockchains:

```
internal/
  blockchain/           # Core blockchain interfaces and types
    types.go           # Common types and interfaces
    service.go         # Generic blockchain service
    solana/            # Solana-specific implementation
      provider.go      # Solana provider implementation
      raydium.go       # Raydium DEX integration
```

### Core Components

1. **Blockchain Service**: Generic interface for blockchain operations
   - Wallet management
   - Transaction handling
   - Balance queries

2. **Provider Interface**: Blockchain-specific implementations
   - Network information
   - Wallet operations
   - Transaction operations

3. **DEX Integration**: Decentralized exchange integration
   - Token swaps
   - Price discovery
   - Liquidity management

## Database Schema

### Wallets Table
```sql
CREATE TABLE wallets (
    id TEXT PRIMARY KEY,
    network TEXT NOT NULL,
    address TEXT NOT NULL,
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    last_updated_at BIGINT NOT NULL,
    UNIQUE(network, address)
)
```

### Transactions Table
```sql
CREATE TABLE blockchain_transactions (
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
```

## API Endpoints

### Wallet Operations
- `POST /api/v1/wallets` - Create a new wallet
- `GET /api/v1/wallets/{network}/{address}` - Get wallet details
- `GET /api/v1/wallets/{network}/{address}/balance` - Get wallet balance

### Transaction Operations
- `POST /api/v1/transactions/buy` - Execute a buy transaction
- `POST /api/v1/transactions/sell` - Execute a sell transaction
- `GET /api/v1/transactions/{network}/{txID}` - Get transaction details
- `GET /api/v1/wallets/{network}/{address}/transactions` - Get wallet transactions

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Start the services:
```bash
docker-compose up -d
```

4. Initialize the database:
```bash
curl -X POST http://localhost:8080/api/v1/memecoins/update
```

## Development

### Adding a New Blockchain

1. Create a new package under `internal/blockchain/{chain}`
2. Implement the `Provider` interface
3. Add necessary types and utilities
4. Register the provider with the blockchain service

Example:
```go
provider := ethereum.NewProvider(false)
service := blockchain.NewService()
service.RegisterProvider(provider)
```

### Testing

Run the test suite:
```bash
go test ./...
```

## Security Considerations

1. **Private Key Storage**: Private keys are stored encrypted in the database
2. **Transaction Validation**: All transactions are validated before execution
3. **Rate Limiting**: API endpoints are rate-limited
4. **Input Validation**: All user inputs are validated and sanitized

## Future Improvements

1. Add support for more blockchains:
   - Ethereum
   - Polygon
   - BNB Chain

2. Enhance DEX integration:
   - Multiple DEX support per chain
   - Automatic route optimization
   - Price impact calculation

3. Security enhancements:
   - Hardware wallet support
   - Multi-signature wallets
   - Advanced encryption for private keys

4. Performance optimizations:
   - Caching layer
   - Transaction batching
   - Parallel processing

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License 