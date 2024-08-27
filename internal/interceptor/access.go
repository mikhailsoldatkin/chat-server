package interceptor

import (
	"context"
	"log"

	accessProto "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Interceptor creates a gRPC server interceptor that checks access using the provided AccessV1Client.
func Interceptor(client accessProto.AccessV1Client) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		md, _ := metadata.FromIncomingContext(ctx)
		outgoingCtx := metadata.NewOutgoingContext(ctx, md)

		_, err := client.Check(outgoingCtx, &accessProto.CheckRequest{EndpointAddress: info.FullMethod})
		if err != nil {
			log.Println("access check failed:", err.Error())
			return nil, err
		}

		log.Printf("access to %s endpoint granted\n", info.FullMethod)

		return handler(ctx, req)
	}
}
