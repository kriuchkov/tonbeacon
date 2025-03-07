package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	walletutils "github.com/xssnick/tonutils-go/ton/wallet"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "tonbeacon",
		Short: "TON Beacon CLI tool",
	}

	var wordCount int
	generateSeedCmd := &cobra.Command{
		Use:   "generate-seed",
		Short: "Generate a new seed phrase (mnemonic)",
		Long:  "Generate a new cryptographically secure random seed phrase with specified number of words",
		Run: func(_ *cobra.Command, args []string) {
			phrase := walletutils.NewSeed()
			fmt.Println("Your seed phrase:")
			fmt.Println(phrase)
			fmt.Println("\nWARNING: Store this seed phrase securely. Anyone with access to this phrase will have access to your wallet.")
		},
	}

	generateSeedCmd.Flags().IntVarP(&wordCount, "words", "w", 24, "Number of words in the seed phrase (12, 15, 18, 21, or 24)")
	rootCmd.AddCommand(generateSeedCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}
