package bootx

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

// HTTPRun encapsulates the boilerplate for running an HTTP server with graceful shutdown.
func HTTPRun(shutdownCtx context.Context, server *http.Server, addr string) error {
	g, ctx := errgroup.WithContext(shutdownCtx)

	g.Go(func() error {
		log.Printf("http listen on %s", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		toCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return server.Shutdown(toCtx)
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
