package chat

import (
	"context"

	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
)

// SendMessage handles sending a message to a chat from a chat.
func (s *serv) SendMessage(ctx context.Context, req *pb.SendMessageRequest) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.chatRepository.SendMessage(ctx, req)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
