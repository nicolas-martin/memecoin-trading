-- Insert sample coins
INSERT INTO coins (symbol, name, contract_address) VALUES
('BTC', 'Bitcoin', NULL),
('ETH', 'Ethereum', NULL),
('SOL', 'Solana', NULL),
('DOGE', 'Dogecoin', NULL),
('PEPE', 'Pepe', '0x6982508145454ce325ddbe47a25d4ec3d2311933');

-- Insert sample coin prices
INSERT INTO coin_prices (coin_id, price, market_cap, volume_24h)
SELECT 
    id,
    CASE 
        WHEN symbol = 'BTC' THEN 50000
        WHEN symbol = 'ETH' THEN 3000
        WHEN symbol = 'SOL' THEN 100
        WHEN symbol = 'DOGE' THEN 0.1
        ELSE 0.000001
    END,
    1000000000,
    100000000
FROM coins; 