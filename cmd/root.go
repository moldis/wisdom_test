package cmd

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var (
	defaultAddress = "0.0.0.0:3333"
	defaultSecret  = "my__own__secret"
)

var rootCmd = &cobra.Command{
	Use:   "wisdom",
	Short: "Wisdom TCP server",
}

func Execute() {
	rootCtx, rootCancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer rootCancel()

	if err := rootCmd.ExecuteContext(rootCtx); err != nil {
		log.Error().Msg("failed to execute a root cmd")
	}
}
