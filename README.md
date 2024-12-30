# Meme Trader

A trading platform for meme coins on the Solana blockchain.

## Features

- Real-time meme coin price tracking from multiple sources:
  - DexScreener
  - CoinGecko
  - Jupiter
- Robust blockchain provider translation layer:
  - Automatic fallback between providers
  - Rate limiting and health monitoring
  - Priority-based provider selection
  - Thread-safe operations
- Price history tracking
- Secure wallet management with encrypted private keys
- Support for multiple DEX integrations:
  - Raydium
  - Orca

## Project Structure

```
.
├── cmd/
│   └── api/              # API server entry point
├── internal/
│   ├── api/             # API handlers and middleware
│   ├── blockchain/      # Blockchain service and providers
│   │   ├── solana/     # Solana-specific implementation
│   │   └── types.go    # Common blockchain interfaces
│   ├── repository/      # Database models and queries
│   └── services/        # Business logic
├── migrations/          # Database migrations
└── docker-compose.yml   # Docker configuration
```

## Blockchain Provider Layer

The application uses a robust provider translation layer that abstracts away differences between various blockchain providers. This ensures high availability and consistent behavior across different providers.

### Key Features

1. **Provider Management**
   - Register multiple providers per network
   - Configure provider priorities
   - Automatic provider health monitoring
   - Rate limit management per provider

2. **Automatic Fallback**
   - Seamless switching between providers
   - Handles rate limits and downtime
   - Prioritizes providers based on health and configuration

3. **Health Monitoring**
   - Tracks provider status (Healthy/Degraded/Unhealthy)
   - Monitors consecutive errors
   - Automatic recovery after failures

4. **Rate Limiting**
   - Configurable requests per time window
   - Automatic rate limit detection
   - Graceful fallback when limits are reached

### Usage Example

```go
// Initialize the blockchain service
service := blockchain.NewService()

// Register primary provider
service.RegisterProvider(primaryProvider)

// Register backup provider with custom configuration
service.RegisterProviderWithConfig(backupProvider, blockchain.ProviderConfig{
    Priority:           2,
    RequestsPerWindow:  1000,
    WindowDuration:     time.Minute,
    HealthCheckPeriod:  time.Minute,
    MaxConsecutiveErrs: 3,
})

// Use the service - fallback is automatic
wallet, err := service.CreateWallet(ctx, blockchain.NetworkSolana)
```

## API Endpoints

### Meme Coins

- `GET /api/v1/memecoins` - Get list of top meme coins
  - Query parameters:
    - `limit` (optional) - Number of coins to return (default: 50)
  - Response includes:
    - Basic coin information (symbol, name)
    - Current price and market data
    - Logo URL and description
    - 24h price changes

- `GET /api/v1/memecoins/{id}` - Get detailed information about a specific coin
  - Response includes:
    - All coin information
    - Price history data

- `POST /api/v1/memecoins/update` - Trigger update of meme coin data
  - Fetches latest data from all providers
  - Updates database with new information

## Setup

1. Install dependencies:
```bash
make install
```

2. Start the backend services:
```bash
make start-backend
```

3. Start the mobile app:
```bash
make start-mobile
```

4. Update meme coin data:
```bash
make update-memecoins
```

## Development

### Running Tests

```bash
# Run all backend tests
make test-backend

# Run blockchain provider tests only
make test-blockchain

# Run specific provider tests (e.g., Solana)
make test-providers

# Run tests with coverage report
make test-coverage

# Run mobile app tests
make test-mobile
```

### Test Coverage

The project maintains comprehensive test coverage for the blockchain provider layer:
- Provider manager tests (registration, fallback, rate limiting)
- Health check and recovery tests
- Priority-based selection tests
- Thread safety tests
- Individual provider implementation tests

To view the test coverage report:
1. Run `make test-coverage`
2. Open `coverage.html` in your browser

### Database Migrations

The project uses SQL migrations for database schema management. Migrations are automatically applied when starting the backend services.

To manually run migrations:

```bash
migrate -path migrations -database "postgres://user:password@localhost:5432/memetrader?sslmode=disable" up
```

To revert migrations:

```bash
migrate -path migrations -database "postgres://user:password@localhost:5432/memetrader?sslmode=disable" down
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 