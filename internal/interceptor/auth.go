package interceptor

import (
	"context"

	"github.com/mikhailsoldatkin/chat-server/internal/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor creates a gRPC server interceptor that checks access using the provided gRPC authentication client.
func AuthInterceptor(cl client.AuthClient) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		ctx = metadata.NewOutgoingContext(ctx, md)

		err := cl.CheckAccess(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
