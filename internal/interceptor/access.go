package interceptor

import (
	"context"
	"log"

	accessProto "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

//func Interceptor(client accessProto.AccessV1Client) func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
//	return func(
//		ctx context.Context,
//		req interface{},
//		info *grpc.UnaryServerInfo,
//		handler grpc.UnaryHandler,
//	) (any, error) {
//
//		md, _ := metadata.FromIncomingContext(ctx)
//		outgoingCtx := metadata.NewOutgoingContext(ctx, md)
//
//		_, err := client.Check(outgoingCtx, &accessProto.CheckRequest{EndpointAddress: info.FullMethod})
//		if err != nil {
//			return nil, err
//		}
//
//		return handler(ctx, req)
//	}
//}

func Interceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {

	log.Println("Interceptor called...")
	log.Println(info.FullMethod)

	address := "localhost:50051"
	//conn, err := grpc.NewClient(
	//	address,
	//	grpc.WithTransportCredentials(creds),
	//)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//closer.Add(conn.Close)

	accessClient := accessProto.NewAccessV1Client(conn)

	md, _ := metadata.FromIncomingContext(ctx)
	outgoingCtx := metadata.NewOutgoingContext(ctx, md)

	_, err = accessClient.Check(outgoingCtx, &accessProto.CheckRequest{EndpointAddress: info.FullMethod})
	if err != nil {
		log.Println("accessClient error", err)
		return nil, err
	}

	return handler(ctx, req)
}
