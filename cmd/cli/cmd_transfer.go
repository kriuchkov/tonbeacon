package main

import (
	"context"
	"fmt"
	"os"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	walletutils "github.com/xssnick/tonutils-go/ton/wallet"
)

func cmdTransfer(ctx context.Context) *cobra.Command {
	command := cobra.Command{
		Use:   "transfer",
		Short: "Transfer funds between two wallets",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := LoadConfig()
			if err != nil {
				log.Warn().Err(err).Msg("config loading")
				os.Exit(64)
			}

			destAddr, _ := cmd.Flags().GetString("to")
			amount, _ := cmd.Flags().GetString("amount")
			mainnet, _ := cmd.Flags().GetBool("mainnet")

			log.Debug().Str("to", destAddr).Str("amount", amount).Bool("is_mainnet", mainnet).Msg("transfer")

			liteClient, err := setupLiteClient(ctx, mainnet)
			if err != nil {
				log.Error().Err(err).Msg("liteclient setup")
				os.Exit(1)
			}

			masterWallet, err := walletutils.FromSeed(liteClient, cfg.Master.GetSeed(), cfg.Master.Version)
			if err != nil {
				log.Warn().Err(err).Msg("master wallet creation")
				os.Exit(64)
			}

			log.Debug().Str("address", masterWallet.Address().String()).Msg("master wallet")

			toAddress, err := address.ParseAddr(destAddr)
			if err != nil {
				log.Error().Err(err).Str("address", destAddr).Msg("parse destination address")
				os.Exit(1)
			}

			bounce := false

			transfer, err := masterWallet.BuildTransfer(toAddress, tlb.MustFromTON(amount), bounce, "")
			if err != nil {
				log.Error().Err(err).Msg("build transfer")
				os.Exit(1)
			}

			fmt.Println("Sending transaction... ")

			tx, block, err := masterWallet.SendWaitTransaction(ctx, transfer)
			if err != nil {
				log.Error().Err(err).Msg("send transaction")
				os.Exit(1)
			}

			fmt.Println("Transaction sent, txid:", tx.String())
			fmt.Println("Block:", block.SeqNo)
		},
	}

	command.Flags().StringP("to", "t", "", "Destination address")
	command.Flags().StringP("amount", "a", "0.0", "Amount to transfer in nano TON")
	command.Flags().Bool("mainnet", false, "Use testnet")
	return &command
}
