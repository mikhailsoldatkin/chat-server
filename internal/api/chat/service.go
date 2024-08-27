package chat

import (
	"github.com/mikhailsoldatkin/chat-server/internal/service"
	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
)

// Implementation provides methods for handling chat-related gRPC requests.
type Implementation struct {
	pb.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewImplementation creates a new instance of Implementation with the given chat service.
func NewImplementation(
	chatService service.ChatService,
) *Implementation {
	return &Implementation{
		chatService: chatService,
	}
}
