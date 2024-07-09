package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"

	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
)

const (
	grpcPort   = 50052
	serverName = "chat"
)

type server struct {
	pb.UnimplementedChatV1Server
}

// Create handles the creation of a new chat with provided users.
func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	log.Printf(color.GreenString("chat 21 created, users: %v", req.Users))
	return &pb.CreateResponse{Id: 21}, nil
}

// Delete removes a chat by ID.
func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf(color.GreenString("chat %v deleted", req.GetId()))
	return &emptypb.Empty{}, nil
}

// SendMessage handles sending message to chat from user.
func (s *server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf(color.GreenString("message %q from user %v sent to chat %v", req.Text, req.FromUser, req.ChatId))
	return &emptypb.Empty{}, nil
}

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		log.Fatalf(color.RedString("failed to listen: %v", err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterChatV1Server(s, &server{})

	log.Printf(color.GreenString("%v server listening at %v", serverName, lis.Addr()))

	if err = s.Serve(lis); err != nil {
		log.Fatalf(color.RedString("failed to serve: %v", err))
	}
}
