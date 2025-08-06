DROP INDEX IF EXISTS idx_price_history_currency_id_timestamp;
DROP TABLE IF EXISTS price_history;
DROP TABLE IF EXISTS tracked_currencies;
DROP EXTENSION IF EXISTS "uuid-ossp";