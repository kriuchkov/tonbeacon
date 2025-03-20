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
    id TEXT PRIMARY KEY,
    sender TEXT NOT NULL,
    receiver TEXT NOT NULL,
    amount DECIMAL NOT NULL,
    block_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    sender_is_ours BOOLEAN NOT NULL,
    receiver_is_ours BOOLEAN NOT NULL
);
