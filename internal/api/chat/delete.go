package chat

import (
	"context"

	"github.com/mikhailsoldatkin/chat-server/internal/customerrors"
	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete removes a chat by ID.
func (i *Implementation) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, customerrors.ConvertError(err)
	}

	return &emptypb.Empty{}, nil
}
