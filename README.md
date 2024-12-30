# Meme Trader

A trading platform for meme coins on the Solana blockchain.

## Features

- Real-time meme coin price tracking from multiple sources:
  - DexScreener
  - CoinGecko
  - Jupiter
- Automatic fallback between data providers for better reliability
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
│   ├── repository/      # Database models and queries
│   └── services/        # Business logic
├── migrations/          # Database migrations
└── docker-compose.yml   # Docker configuration
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
# Run backend tests
make test-backend

# Run mobile app tests
make test-mobile
```

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