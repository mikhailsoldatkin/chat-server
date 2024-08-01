package repository

import (
	"context"

	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
)

// ChatRepository defines the interface for chat-related database operations.
type ChatRepository interface {
	Create(ctx context.Context, users []int64) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, req *pb.SendMessageRequest) error
}
