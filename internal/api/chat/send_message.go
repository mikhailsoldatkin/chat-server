package chat

import (
	"context"

	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendMessage handles sending a message from a user to a chat.
func (i *Implementation) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.chatService.SendMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
