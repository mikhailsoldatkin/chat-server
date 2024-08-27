package interceptor

import (
	"context"
	"log"

	accessProto "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	address     = "localhost:50051"
	accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ4NDQ5MDMsInVzZXJuYW1lIjoibWlzaGEiLCJyb2xlIjoiQURNSU4ifQ.DclaC744l-sXM52FTTTjUP8KR-BIKCP6ucyQ1VoQbI0"
)

func Interceptor(
	client accessProto.AccessV1Client,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		log.Println("interceptor called...")

		md, _ := metadata.FromIncomingContext(ctx)
		outgoingCtx := metadata.NewOutgoingContext(ctx, md)

		_, err := client.Check(outgoingCtx, &accessProto.CheckRequest{EndpointAddress: info.FullMethod})
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		log.Printf("access to %s endpoint granted\n", info.FullMethod)

		return handler(ctx, req)
	}
}
