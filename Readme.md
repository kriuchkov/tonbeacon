# TonBeacon

> [!Warning]
> The project is in the early development stage and should not be used in. The codebase is subject to change, and the .
> documentation may be incomplete or outdated.

**TonBeacon** is a self-hosted B2B solution providing infrastructure for generating user wallets on the TON (The Open Network) blockchain, tracking incoming transactions, and automatically centralizing funds to a master wallet. The primary goal of the project is to make the TON blockchain accessible for financial solutions and extend beyond the Telegram ecosystem, enabling B2B applications to offer custodial wallets with account management and fund aggregation for their users without relying on memo phrases or compromising the master wallet.



## Project Goals

- **Wallet Generation**: Create unique subwallets for each user based on a single master key.
- **Transaction Tracking**: Monitor incoming transactions using the Outbox pattern with idempotent delivery to Kafka.
- **Funds Centralization**: Run a collector that periodically consolidates funds from subwallets to the master wallet.

## Key Features

- **Self-Hosted**: Deploy on your own servers with full autonomy.
- **API**: gRPC and HTTP interfaces for managing wallets and retrieving transaction data.
- **Scalability**: Supports Kafka for asynchronous event processing and PostgreSQL for data storage.
- **Security**: Configuration via environment variables and secrets, with graceful shutdown for proper termination.

## Architecture

<p align="center">
 <img src="https://raw.githubusercontent.com/kriuchkov/tonbeacon/refs/heads/master/docs/.images/logo.svg">
</p>

The project uses a hexagonal architecture to separate business logic from infrastructure:

- **Domain**: Core with business logic (wallet generation, transaction tracking, collector).
- **Adapters**: Integration with TON (tonutils-go), PostgreSQL (Bun), Kafka, gRPC/HTTP.
- **Outbox**: Guaranteed event delivery to Kafka with idempotency via unique keys.
- **Collector**: Periodic process for transferring funds to the master wallet.

### Components

- **PostgreSQL**: Stores wallet, transaction, and Outbox event data.
- **Kafka**: Message queue for asynchronous transaction processing (KRaft mode, no Zookeeper).
- **Flyway**: Database migration management.
- **Go Application**: Main service with gRPC/HTTP APIs.

## Requirements

- Go 1.24+
- Docker and Docker Compose (for local deployment)
- PostgreSQL 15+
- Kafka 3.8.0+

## License

MIT License
