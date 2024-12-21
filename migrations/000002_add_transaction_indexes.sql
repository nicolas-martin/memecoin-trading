-- Up migration
-- Add indexes for frequently queried columns
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_transactions_type_status ON transactions(type, status);
CREATE INDEX idx_transactions_user_wallet_coin ON transactions(user_id, wallet_id, coin_id);

-- Add composite index for leaderboard queries
CREATE INDEX idx_transactions_status_created_at ON transactions(status, created_at DESC);

-- Down migration
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_type_status;
DROP INDEX IF EXISTS idx_transactions_user_wallet_coin;
DROP INDEX IF EXISTS idx_transactions_status_created_at; 