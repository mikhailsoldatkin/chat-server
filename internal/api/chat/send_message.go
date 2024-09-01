package chat

import (
	"context"

	"github.com/mikhailsoldatkin/chat-server/internal/customerrors"
	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendMessage handles sending a message from a user to a chat.
func (i *Implementation) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.client.CheckUsersExist(ctx, []int64{req.FromUser})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = i.chatService.SendMessage(ctx, req)
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &emptypb.Empty{}, nil
}
