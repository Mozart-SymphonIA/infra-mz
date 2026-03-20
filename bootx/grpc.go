package bootx

import (
	"context"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// GRPCRun encapsulates the boilerplate for running a gRPC server with graceful shutdown.
func GRPCRun(shutdownCtx context.Context, server *grpc.Server, lis net.Listener) error {
	g, ctx := errgroup.WithContext(shutdownCtx)

	g.Go(func() error {
		return server.Serve(lis)
	})

	g.Go(func() error {
		<-ctx.Done()
		server.GracefulStop()
		return nil
	})

	return g.Wait()
}
