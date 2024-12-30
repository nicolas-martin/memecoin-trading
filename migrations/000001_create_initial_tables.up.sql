-- Create memecoins table
CREATE TABLE IF NOT EXISTS memecoins (
    id TEXT PRIMARY KEY,
    symbol TEXT NOT NULL,
    name TEXT NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    market_cap DOUBLE PRECISION,
    volume_24h DOUBLE PRECISION,
    price_change_24h DOUBLE PRECISION,
    price_change_percentage_24h DOUBLE PRECISION,
    contract_address TEXT NOT NULL,
    data_provider TEXT NOT NULL,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    logo_url TEXT,
    description TEXT
);

-- Create price_history table
CREATE TABLE IF NOT EXISTS price_history (
    coin_id TEXT REFERENCES memecoins(id),
    price DOUBLE PRECISION NOT NULL,
    volume DOUBLE PRECISION,
    timestamp BIGINT NOT NULL,
    PRIMARY KEY (coin_id, timestamp)
); 