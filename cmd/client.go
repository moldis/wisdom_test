package cmd

import (
	"github.com/moldis/wisdom_test/internal/call"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	client.Flags().String("addr", defaultAddress, "Server address")
	rootCmd.AddCommand(client)
}

var client = &cobra.Command{
	Use:   "client",
	Short: "Client for wisdom TCP server",
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")

		caller := call.NewCaller(
			call.WithServerAddr(addr),
			call.WithProtocol("tcp"),
		)

		if err := caller.Run(cmd.Context()); err != nil {
			log.Err(err).Msgf("failed to get a wisdom quote")
		}
	},
}
