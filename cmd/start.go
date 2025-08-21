package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chrishrb/hoval-gateway/hoval/service"
	"github.com/chrishrb/hoval-gateway/store/inmemory"
	"github.com/chrishrb/hoval-gateway/transport/mock"
	"github.com/spf13/cobra"
	"k8s.io/utils/clock"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dummyBus := mock.NewMockBus()
		consumer := mock.NewConsumer(dummyBus)
		// sender := mock.NewSender(dummyBus)

		store := inmemory.NewStore(clock.RealClock{})
		consumeSvc := service.NewConsumeService(store)

		// Get messages from the can bus
		errCh := make(chan error, 1)
		canConn, err := consumer.Consume(context.Background(), consumeSvc)
		if err != nil {
			errCh <- err
		}

		// Add some data
		err = consumer.ReadFromFile(context.Background(), "transport/mock/testfiles/hoval_data_1.log")
		if err != nil {
			errCh <- err
		}

		if canConn != nil {
			err := canConn.Close()
			if err != nil {
				slog.Warn("disconnecting from can", "err", err)
			}
		}

		// Get all store data
		for _, d := range store.ListDevices() {
			fmt.Printf("device: %d, %s\n", d.Address, d.GetName())
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
}
