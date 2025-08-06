CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tracked_currencies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol VARCHAR(10) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS price_history (
    currency_id UUID NOT NULL,
    price NUMERIC(20, 8) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_currency
        FOREIGN KEY(currency_id)
        REFERENCES tracked_currencies(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_price_history_currency_id_timestamp ON price_history (currency_id, timestamp DESC);