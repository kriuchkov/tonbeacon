package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "tonbeacon",
		Short: "TON Beacon CLI tool",
	}

	ctx := context.Background()
	rootCmd.AddCommand(cmdGenerateSeed(ctx))
	rootCmd.AddCommand(cmdTransfer(ctx))
	rootCmd.AddCommand(cmdAccount(ctx))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}
