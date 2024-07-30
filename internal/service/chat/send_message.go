package chat

import (
	"context"

	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
)

// SendMessage handles sending a message to a chat from a chat.
func (s *serv) SendMessage(ctx context.Context, req *pb.SendMessageRequest) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.chatRepository.SendMessage(ctx, req)
		if err != nil {
			return err
		}

		return nil
	})
}
