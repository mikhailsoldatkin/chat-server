package chat

import (
	"context"

	"github.com/mikhailsoldatkin/chat-server/internal/customerrors"
	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create handles the creation of a new chat in the system.
func (i *Implementation) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	if len(req.Users) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no users provided for the chat")
	}

	id, err := i.chatService.Create(ctx, req.GetUsers())
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &pb.CreateResponse{Id: id}, nil
}
