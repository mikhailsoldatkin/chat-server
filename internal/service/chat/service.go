package chat

import (
	"github.com/mikhailsoldatkin/chat-server/internal/client/db"
	"github.com/mikhailsoldatkin/chat-server/internal/repository"
	"github.com/mikhailsoldatkin/chat-server/internal/service"
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
