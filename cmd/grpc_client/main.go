package main

import (
	"context"
	"fmt"
	"log"

	pbAccess "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	address     = "localhost:50051"
	accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ4NTA3OTgsInVzZXJuYW1lIjoibWlzaGEiLCJyb2xlIjoiQURNSU4ifQ.XoripxJ4V3KDBVbSmC0FYdGTDm77oseJeTH1WoNzDsM"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("cert/service.pem", "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := pbAccess.NewAccessV1Client(conn)

	ctx := context.Background()
	md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err = c.Check(ctx, &pbAccess.CheckRequest{EndpointAddress: "/chat_v1.ChatV1/Create"})
	if err != nil {
		log.Fatalf("failed to Check(): %v", err)
	}

	fmt.Println("access granted")
}
