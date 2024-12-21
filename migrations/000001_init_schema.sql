-- Up migration
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    address VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE coins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    contract_address VARCHAR(255) UNIQUE NOT NULL,
    logo_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE coin_prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    coin_id UUID NOT NULL REFERENCES coins(id),
    price DECIMAL(20,8) NOT NULL,
    market_cap DECIMAL(30,2),
    volume_24h DECIMAL(30,2),
    price_change_24h DECIMAL(20,8),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    wallet_id UUID NOT NULL REFERENCES wallets(id),
    coin_id UUID NOT NULL REFERENCES coins(id),
    type VARCHAR(20) NOT NULL, -- 'BUY' or 'SELL'
    amount DECIMAL(20,8) NOT NULL,
    price DECIMAL(20,8) NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'PENDING', 'COMPLETED', 'FAILED'
    tx_hash VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_wallet_id ON transactions(wallet_id);
CREATE INDEX idx_transactions_coin_id ON transactions(coin_id);
CREATE INDEX idx_coin_prices_coin_id_timestamp ON coin_prices(coin_id, timestamp);

-- Down migration
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS coin_prices;
DROP TABLE IF EXISTS coins;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS users; 