package chat

import (
	"context"
)

// Delete removes a chat from the system by ID.
func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.chatRepository.Delete(ctx, id)
	})
}
