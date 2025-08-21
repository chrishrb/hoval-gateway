package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chrishrb/hoval-gateway/hoval/service"
	"github.com/chrishrb/hoval-gateway/store/inmemory"
	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/chrishrb/hoval-gateway/transport/can"
	"github.com/chrishrb/hoval-gateway/transport/mock"
	"github.com/spf13/cobra"
	"k8s.io/utils/clock"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen on Hoval CAN bus for messages and print to console",
	RunE: func(cmd *cobra.Command, args []string) error {
		canInterface, _ := cmd.Flags().GetString("interface")
		mockFile, _ := cmd.Flags().GetString("mock-file")

		store := inmemory.NewStore(clock.RealClock{})
		consumeSvc := service.NewConsumeService(store)

		handlerFunc := func(ctx context.Context, message *transport.Message) {
			if message == nil {
				slog.Debug("received nil message, ignoring")
				return
			}

			hovalMsg, err := consumeSvc.FromTransportMessage(*message)
			if err != nil {
				slog.Error("failed to convert transport message to hoval message", "error", err)
				return
			}

			if hovalMsg.Datapoint == nil {
				return
			}

			slog.Debug("received hoval message", "message", hovalMsg)
		}

		var canConn transport.Connection
		var err error
		if mockFile != "" {
			slog.Info("Using mock CAN bus", "file", mockFile)

			mockBus := mock.NewMockBus()
			consumer := mock.NewConsumer(mockBus)

			// Get messages from the can bus
			errCh := make(chan error, 1)
			canConn, err = consumer.Consume(context.Background(), transport.MessageHandlerFunc(handlerFunc))
			if err != nil {
				errCh <- err
			}

			// Add some data
			if mockFile != "" {
				err = consumer.ReadFromFile(context.Background(), mockFile)
				if err != nil {
					errCh <- err
				}
			}
		} else {
			if canInterface == "" {
				return fmt.Errorf("CAN interface must be specified")
			}
			slog.Error("Listening on CAN interface", "interface", canInterface)

			consumer := can.NewConsumer(
				can.WithCANDevice[can.Consumer](canInterface),
			)

			// Get messages from the can bus
			errCh := make(chan error, 1)
			canConn, err = consumer.Consume(context.Background(), consumeSvc)
			if err != nil {
				errCh <- err
			}
		}

		if canConn != nil {
			err := canConn.Close()
			if err != nil {
				slog.Warn("disconnecting from can", "err", err)
			}
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	listenCmd.Flags().StringP("interface", "i", "", "CAN interface to listen on (e.g., can0, vcan0)")
	listenCmd.Flags().String("mock-file", "", "File containing CAN messages to read in mock mode")
}
