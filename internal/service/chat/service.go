package chat

import (
	"context"

	"github.com/mikhailsoldatkin/chat-server/internal/repository"
	"github.com/mikhailsoldatkin/chat-server/internal/service"
	"github.com/mikhailsoldatkin/platform_common/pkg/db"
)

var _ service.ChatService = (*serv)(nil)

type serv struct {
	chatRepository repository.ChatRepository
	txManager      db.TxManager
}

// NewService creates a new instance of the chat service.
func NewService(chatRepository repository.ChatRepository, txManager db.TxManager) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}

// No-op implementation for TxManager
type noOpTxManager struct{}

func (noOpTxManager) ReadCommitted(ctx context.Context, f db.Handler) error {
	return f(ctx)
}

// NewMockService creates a new mock instance of the chat service.
func NewMockService(deps ...any) service.ChatService {
	srv := serv{
		txManager: noOpTxManager{},
	}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.ChatRepository:
			srv.chatRepository = s

		}
	}

	return &srv
}
