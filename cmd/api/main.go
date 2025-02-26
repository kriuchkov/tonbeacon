package main

import (
	"context"
	"database/sql"
	"net"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	liteclientutils "github.com/xssnick/tonutils-go/liteclient"
	tonutils "github.com/xssnick/tonutils-go/ton"
	walletutils "github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/kriuchkov/tonbeacon/adapters/grpc"
	"github.com/kriuchkov/tonbeacon/adapters/repository"
	"github.com/kriuchkov/tonbeacon/adapters/ton"
	"github.com/kriuchkov/tonbeacon/ports/account"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cfg, err := LoadConfig()
	if err != nil {
		log.Panic().Err(err).Msg("config loading")
	}

	client := liteclientutils.NewConnectionPool()

	if err := client.AddConnectionsFromConfigUrl(ctx, "https://tonutils.com/testnet-global.config.json"); err != nil {
		panic(err)
	}

	liteClient := tonutils.NewAPIClient(client, tonutils.ProofCheckPolicySecure)

	masterWallet, err := walletutils.FromSeed(liteClient, cfg.Master.Seed, cfg.Master.Version)
	if err != nil {
		log.Panic().Err(err).Msg("master wallet creation")
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Database.DSN())))
	db := bun.NewDB(sqldb, pgdialect.New())

	repositoryAdapter := repository.New(db)

	accountSvc := account.New(account.Options{
		WalletManager:   ton.NewWalletAdapter(liteClient, masterWallet),
		TxManager:       repository.NewTxRepository(db),
		DatabaseManager: repositoryAdapter,
		EventManager:    nil, //TODO: add kafka manager
	})

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Panic().Err(err).Msg("grpc server listen")
	}

	log.Info().Str("port", cfg.GRPCPort).Msg("grpc server started")

	go func() {
		grpcServer := grpc.NewTonBeacon(accountSvc)
		if err := grpcServer.Run(lis); err != nil {
			log.Panic().Err(err).Msg("grpc server run")
		}
	}()

	<-ctx.Done()

	log.Info().Msg("shutting down")
}
