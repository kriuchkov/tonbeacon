CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    wallet_id SERIAL UNIQUE NOT NULL,
    ton_address VARCHAR(255) NULL,
    is_closed BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE outbox_events (
	id BIGSERIAL PRIMARY KEY,
	event_type TEXT NOT NULL,
	payload TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
	processed BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    account_addr TEXT,
    lt BIGINT NOT NULL,
    prev_tx_hash TEXT NOT NULL,
    prev_tx_lt BIGINT NOT NULL,
    sender TEXT NOT NULL,
    receiver TEXT NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    total_fees DOUBLE PRECISION NOT NULL,
    exit_code INT NOT NULL,
    success BOOLEAN NOT NULL,
    message_type TEXT NOT NULL,
    bounce BOOLEAN NOT NULL,
    bounced BOOLEAN NOT NULL,
    body TEXT NOT NULL,
    block_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    account_status TEXT NOT NULL,
    compute_gas_used INT NOT NULL,
    description TEXT
);

CREATE INDEX idx_transactions_sender ON transactions (sender);
CREATE INDEX idx_transactions_receiver ON transactions (receiver);
