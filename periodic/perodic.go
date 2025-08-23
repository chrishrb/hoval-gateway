package periodic

import (
	"context"
	"log/slog"
	"time"

	"github.com/chrishrb/hoval-gateway/store"
	"k8s.io/utils/clock"
)

func RunPeriodic(ctx context.Context, engine store.Engine, clock clock.PassiveClock, runEvery time.Duration) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down run periodic")
			return
		case <-time.After(runEvery):
			// TODO: implement periodic tasks here
		}
	}
}
