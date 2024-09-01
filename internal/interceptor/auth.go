package interceptor

import (
	"context"

	grpcClient "github.com/mikhailsoldatkin/chat-server/internal/client/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor creates a gRPC server interceptor that checks access using the provided gRPC client.
func AuthInterceptor(cl grpcClient.AuthClient) grpc.UnaryServerInterceptor {
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
