package cmd

import (
	"context"
	"log/slog"

	"github.com/chrishrb/hoval-gateway/api/pubsub"
	"github.com/chrishrb/hoval-gateway/config"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.DefaultConfig
		if configFile != "" {
			err := cfg.LoadFromFile(configFile)
			if err != nil {
				return err
			}
		}

		settings, err := config.Configure(context.Background(), &cfg)
		if err != nil {
			return err
		}

		errCh := make(chan error, 1)

		// Connect to CAN bus and start consuming messages
		transportConn, err := settings.TransportConsumer.Consume(context.Background(), settings.ConsumeHandler)
		if err != nil {
			errCh <- err
		}

		// Connect to PubSub broker and start listening for messages
		var pubSubConn pubsub.Connection
		if settings.PubSubListener != nil {
			pubSubConn, err = settings.PubSubListener.Connect(context.Background(), settings.SendHandler)
			if err != nil {
				errCh <- err
			}
		}

		// Start periodic requests
		periodicRequester := settings.PeriodicRequester
		periodicRequester.Run(context.Background())

		slog.Info("hoval-gateway started")

		err = <-errCh

		if transportConn != nil {
			err := transportConn.Close()
			if err != nil {
				slog.Warn("closing transport connection", "error", err)
			}
		}

		if pubSubConn != nil {
			err := pubSubConn.Disconnect(context.Background())
			if err != nil {
				slog.Warn("disconnecting from broker", "error", err)
			}
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().StringVarP(&configFile, "config-file", "c", "/config/config.toml",
		"The config file to use")
}
