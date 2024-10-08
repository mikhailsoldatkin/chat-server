package chat

import (
	"context"
)

// Create creates a new chat in the system with provided users.
func (s *serv) Create(ctx context.Context, users []int64) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.Create(ctx, users)
		if errTx != nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
