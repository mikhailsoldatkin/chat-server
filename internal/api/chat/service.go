package chat

import (
	accessProto "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	"github.com/mikhailsoldatkin/chat-server/internal/service"
	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
)

// Implementation provides methods for handling chat-related gRPC requests.
type Implementation struct {
	pb.UnimplementedChatV1Server
	chatService  service.ChatService
	accessClient accessProto.AccessV1Client
}

// NewImplementation creates a new instance of Implementation with the given chat service.
func NewImplementation(
	chatService service.ChatService,
	accessClient accessProto.AccessV1Client,
) *Implementation {
	return &Implementation{
		chatService:  chatService,
		accessClient: accessClient,
	}
}
