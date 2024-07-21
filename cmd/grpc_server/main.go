package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mikhailsoldatkin/chat-server/internal/config"
	"github.com/mikhailsoldatkin/chat-server/internal/logger"
	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedChatV1Server
	pool *pgxpool.Pool
}

// Create handles the creation of a new chat with provided users.
func (s *server) Create(_ context.Context, _ *pb.CreateRequest) (*pb.CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method is not implemented")
}

// Delete removes a chat by ID.
func (s *server) Delete(_ context.Context, _ *pb.DeleteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method is not implemented")
}

// SendMessage handles sending message to chat from user.
func (s *server) SendMessage(_ context.Context, _ *pb.SendMessageRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method is not implemented")
}

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, cfg.Database.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterChatV1Server(s, &server{pool: pool})

	logger.Info("%v server listening at %v", cfg.AppName, lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
