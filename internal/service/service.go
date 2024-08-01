package service

import (
	"context"

	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
)

// ChatService defines the interface for chat-related business logic operations.
type ChatService interface {
	Create(ctx context.Context, users []int64) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, req *pb.SendMessageRequest) error
}
