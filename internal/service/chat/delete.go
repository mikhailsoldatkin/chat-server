package chat

import (
	"context"
)

// Delete removes a chat from the system by ID.
func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.chatRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
