package periodic

import (
	"context"
	"log/slog"
	"time"

	"github.com/chrishrb/hoval-gateway/api/pubsub"
	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/service"
)

var (
	Operation = hoval.OperationGetRequest
)

type PeriodicRequest struct {
	ReceiverMask   uint32
	FunctionGroup  uint8
	FunctionNumber uint8
	DatapointID    uint16
}

type PeriodicRequester struct {
	svc      *service.SendService
	requests []PeriodicRequest
	runEvery time.Duration
}

func NewPeriodicRequester(svc *service.SendService, requests []PeriodicRequest, runEvery time.Duration) *PeriodicRequester {
	return &PeriodicRequester{
		svc:      svc,
		requests: requests,
		runEvery: runEvery,
	}
}

func (r *PeriodicRequester) Run(ctx context.Context) {
	for _, req := range r.requests {
		go run(ctx, r.svc, req, r.runEvery)
	}
}

func run(ctx context.Context, svc *service.SendService, req PeriodicRequest, runEvery time.Duration) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down run periodic")
			return
		case <-time.After(runEvery):
			msg := &pubsub.Message{
				FunctionGroup:  req.FunctionGroup,
				FunctionNumber: req.FunctionNumber,
				DatapointID:    req.DatapointID,
				Data:           0,
			}

			err := svc.Send(ctx, req.ReceiverMask, Operation, msg)
			if err != nil {
				slog.Error("error sending periodic message", "error", err)
			}
		}
	}
}
