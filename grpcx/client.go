package grpcx

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Connect creates a gRPC connection with default insecure credentials and logging interceptor.
func Connect(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(clientLoggingInterceptor),
	}

	// Append user opts to override defaults if necessary (though here we just append)
	finalOpts := append(defaultOpts, opts...)

	conn, err := grpc.DialContext(ctx, target, finalOpts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func clientLoggingInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	log.Printf("gRPC Request: %s -> %s", cc.Target(), method)
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		log.Printf("gRPC Error: %s -> %s: %v", cc.Target(), method, err)
	}
	return err
}
