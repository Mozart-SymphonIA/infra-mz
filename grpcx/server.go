package grpcx

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// NewServer creates a new gRPC server with default interceptors (Recovery, Logging)
// and enables reflection.
func NewServer(options ...grpc.ServerOption) *grpc.Server {
	// Default options including interceptors
	defaultOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			recoveryInterceptor,
			loggingInterceptor,
		),
	}

	// Append user provided options
	opts := append(defaultOpts, options...)

	srv := grpc.NewServer(opts...)

	// Enable reflection for debugging
	reflection.Register(srv)

	return srv
}

// Listen creates a net.Listener on the specified port.
func Listen(port string) (net.Listener, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", port, err)
	}
	return lis, nil
}

// recoveryInterceptor handles panics in gRPC handlers.
func recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[gRPC] PANIC in %s: %v\n%s", info.FullMethod, r, debug.Stack())
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}

// loggingInterceptor logs gRPC requests.
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	code := codes.OK
	if err != nil {
		if s, ok := status.FromError(err); ok {
			code = s.Code()
		} else {
			code = codes.Unknown
		}
	}

	log.Printf("[gRPC] %s %s %v %v", info.FullMethod, code, duration, err)
	return resp, err
}
