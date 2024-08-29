package interceptor

import (
	"context"

	pb "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor creates a gRPC server interceptor that checks access using the provided AccessV1Client.
func AuthInterceptor(c pb.AccessV1Client) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		ctx = metadata.NewOutgoingContext(ctx, md)

		_, err := c.Check(ctx, &pb.CheckRequest{Endpoint: info.FullMethod})
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
