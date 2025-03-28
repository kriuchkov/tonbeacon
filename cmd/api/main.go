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
	walletutils "github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/kriuchkov/tonbeacon/adapters/grpc"
	"github.com/kriuchkov/tonbeacon/adapters/repository"
	"github.com/kriuchkov/tonbeacon/adapters/ton"
	"github.com/kriuchkov/tonbeacon/pkg/common"
	"github.com/kriuchkov/tonbeacon/ports/account"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cfg, err := LoadConfig()
	if err != nil {
		log.Panic().Err(err).Msg("config loading")
	}
	log.Info().Any("conf", cfg).Msg("config loaded")

	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Database.DSN()))), pgdialect.New())

	if err = db.PingContext(ctx); err != nil {
		log.Panic().Err(err).Msg("db connection")
	}

	repositoryAdapter := repository.New(db)

	liteClient, err := common.SetupLiteClient(ctx, cfg.IsMainnet)
	if err != nil {
		log.Warn().Err(err).Msg("lite client setup")
		os.Exit(1)
	}

	masterWallet, err := walletutils.FromSeed(liteClient, cfg.Master.GetSeed(), cfg.Master.Version)
	if err != nil {
		log.Panic().Err(err).Msg("master wallet creation")
	}

	accountSvc := account.New(account.Options{
		WalletManager:   ton.NewWalletAdapter(liteClient, masterWallet),
		TxManager:       repository.NewTxRepository(db),
		DatabaseManager: repositoryAdapter,
	})

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Panic().Err(err).Msg("grpc server listen")
	}

	log.Info().Str("port", cfg.GRPCPort).Msg("grpc server started")

	go func() {
		grpcServer := grpc.NewTonBeacon(accountSvc)
		if err = grpcServer.Run(lis); err != nil {
			log.Panic().Err(err).Msg("grpc server run")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("shutting down")
}
