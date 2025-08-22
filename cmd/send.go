package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/chrishrb/hoval-gateway/config"
	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/service"
	"github.com/chrishrb/hoval-gateway/store/inmemory"
	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/chrishrb/hoval-gateway/transport/can"
	"github.com/spf13/cobra"
	"k8s.io/utils/clock"
)

type mockSender struct{}

func (m *mockSender) Send(ctx context.Context, message *transport.Message) error {
	fmt.Printf("Mock message sent: ID=%X, Data=%X\n", message.ID, message.Data[:message.Length])
	return nil
}

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send hoval messages to a CAN bus",
	RunE: func(cmd *cobra.Command, args []string) error {
		canInterface, _ := cmd.Flags().GetString("interface")
		mockMode, _ := cmd.Flags().GetBool("mock")

		// Get message configuration from flags
		senderID, _ := cmd.Flags().GetUint32("sender")
		receiverMask, _ := cmd.Flags().GetUint32("receiver")
		operationStr, _ := cmd.Flags().GetString("operation")
		datapointName, _ := cmd.Flags().GetString("datapoint")
		data, _ := cmd.Flags().GetFloat64("data")

		// Parse operation
		var operation hoval.Operation
		switch operationStr {
		case "response":
			operation = hoval.OperationResponse
		case "get":
			operation = hoval.OperationGetRequest
		case "set":
			operation = hoval.OperationSetRequest
		default:
			// Try to parse as uint8
			if opVal, err := strconv.ParseUint(operationStr, 10, 8); err == nil {
				operation = hoval.Operation(opVal)
			} else {
				return fmt.Errorf("invalid operation: %s (use 'response', 'get', 'set', or numeric value)", operationStr)
			}
		}

		// Handle datapoint name (can be nil)
		if datapointName == "" {
			return fmt.Errorf("datapoint name must be specified")
		}
		dpName := &datapointName

		// Build message
		msg := &hoval.Message{
			SenderID:      senderID,
			ReceiverMask:  receiverMask,
			OperationID:   operation,
			DatapointName: dpName,
			Data:          data,
		}

		var sender transport.Sender
		if mockMode != false {
			slog.Info("Send message to STDOUT")

			sender = &mockSender{}
		} else {
			if canInterface == "" {
				return fmt.Errorf("CAN interface must be specified")
			}

			sender = can.NewSender(
				can.WithCANDevice[can.Sender](canInterface),
			)
		}

		dpProvider, err := config.NewCsvDatapointProvider("config/datapoints.csv")
		if err != nil {
			return err
		}

		store := inmemory.NewStore(clock.RealClock{})
		svc := service.NewSendService(store, dpProvider, sender)

		// Send messages to the can bus
		err = svc.Send(context.Background(), msg)
		if err != nil {
			return err
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("interface", "i", "", "CAN interface to send on (e.g., can0, vcan0)")
	sendCmd.Flags().Bool("mock", false, "Mock mode for testing (if specified, will use mock instead of real CAN)")

	// Message configuration flags
	sendCmd.Flags().Uint32("sender", 1153, "Sender ID for the message")
	sendCmd.Flags().Uint32("receiver", 0, "Receiver for the message (2047 to broadcast)")
	sendCmd.Flags().StringP("operation", "o", "get", "Operation type (response, get, set, or numeric value)")
	sendCmd.Flags().String("datapoint", "", "Datapoint name")
	sendCmd.Flags().Float64("data", 0, "Data to send")
}
