package bootx

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

// WorkerRun encapsulates the boilerplate for running a background worker with graceful shutdown.
func WorkerRun(shutdownCtx context.Context, run func(ctx context.Context) error) error {
	g, ctx := errgroup.WithContext(shutdownCtx)

	g.Go(func() error {
		return run(ctx)
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
