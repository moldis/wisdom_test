package cmd

import (
	"github.com/moldis/wisdom_test/assets"
	"github.com/moldis/wisdom_test/config"
	"github.com/moldis/wisdom_test/internal/pow/zkpow"
	srv "github.com/moldis/wisdom_test/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	server.Flags().String("addr", defaultAddress, "Server address")
	server.Flags().String("secret", defaultSecret, "POW secret")
	server.Flags().Int64("difficulty", int64(1), "POW difficulty")

	rootCmd.AddCommand(server)
}

var server = &cobra.Command{
	Use:   "server",
	Short: "Server to serve wisdom quotes",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		log.Info().Msgf("starting the server")

		addr, _ := cmd.Flags().GetString("addr")
		secret, _ := cmd.Flags().GetString("secret")
		difficulty, _ := cmd.Flags().GetInt64("difficulty")

		powProvider := zkpow.NewZKPoW(1, secret)
		quotes := assets.NewQuote()
		cfg := config.NewConfig(addr, difficulty)

		server := srv.NewTcpServer(cfg, powProvider, quotes)

		go server.Start(ctx)

		log.Info().Msgf("server started %s", addr)
		<-ctx.Done()
	},
}
