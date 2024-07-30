package chat

import (
	"context"

	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
)

// Create handles the creation of a new chat in the system.
func (i *Implementation) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	id, err := i.chatService.Create(ctx, req.GetUsers())
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: id}, nil
}
