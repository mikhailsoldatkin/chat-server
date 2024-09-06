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
	if len(req.UsersIds) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no users provided for the chat")
	}

	err := i.authClient.CheckUsersExist(ctx, req.UsersIds)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := i.chatService.Create(ctx, req.GetUsersIds())
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &pb.CreateResponse{Id: id}, nil
}
