package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kriuchkov/tonbeacon/core/consts"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	liteclientutils "github.com/xssnick/tonutils-go/liteclient"
	tonutils "github.com/xssnick/tonutils-go/ton"
	walletutils "github.com/xssnick/tonutils-go/ton/wallet"
)

func cmdAccount(ctx context.Context) *cobra.Command {
	command := cobra.Command{
		Use:   "account",
		Short: "Account management",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := LoadConfig()
			if err != nil {
				log.Warn().Err(err).Msg("config loading")
				os.Exit(64)
			}

			subWalletID, _ := cmd.Flags().GetUint32("subwallet")
			if err != nil {
				log.Warn().Err(err).Msg("get subwallet id")
				os.Exit(64)
			}

			mainnet, _ := cmd.Flags().GetBool("mainnet")

			configURL := consts.TestNetConfigURL
			if mainnet {
				configURL = consts.MainNetConfigURL
			}

			log.Info().Bool("is_mainnet", mainnet).Uint32("subwallet_id", subWalletID).Msg("config loaded")
			client := liteclientutils.NewConnectionPool()

			if err := client.AddConnectionsFromConfigUrl(ctx, configURL); err != nil {
				log.Warn().Err(err).Msg("add connections from config")
				os.Exit(64)
			}

			liteClient := tonutils.NewAPIClient(client, tonutils.ProofCheckPolicySecure)
			log.Info().Msg("liteclient connected")

			masterWallet, err := walletutils.FromSeed(liteClient, cfg.Master.GetSeed(), cfg.Master.Version)
			if err != nil {
				log.Panic().Err(err).Msg("master wallet creation")
			}

			if subWalletID == 0 {
				fmt.Println("Master wallet address:", masterWallet.Address().String())
				return
			}

			subwallet, err := masterWallet.GetSubwallet(subWalletID)
			if err != nil {
				log.Error().Err(err).Uint32("subwallet_id", subWalletID).Msg("get subwallet")
				os.Exit(1)
			}

			log.Info().Str("addr", subwallet.Address().String()).Uint32("id", subWalletID).Msg("subwallet")
		},
	}

	command.Flags().Uint32P("subwallet", "s", 0, "Subwallet ID")
	command.Flags().Bool("mainnet", false, "Use testnet")
	return &command
}
