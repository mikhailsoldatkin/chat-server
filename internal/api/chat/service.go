package chat

import (
	"context"

	"github.com/mikhailsoldatkin/chat-server/internal/client"
	"github.com/mikhailsoldatkin/chat-server/internal/service"
	pbChat "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
)

// Implementation provides methods for handling chat-related gRPC requests.
type Implementation struct {
	pbChat.UnimplementedChatV1Server
	chatService service.ChatService
	authClient  client.AuthClient
}

// NewImplementation creates a new instance of Implementation with the given chat service and authentication client.
func NewImplementation(
	chatService service.ChatService,
	authClient client.AuthClient,
) *Implementation {
	return &Implementation{
		chatService: chatService,
		authClient:  authClient,
	}
}

// No-op implementation for auth client
type noOpClient struct{}

func (noOpClient) CheckAccess(_ context.Context, _ string) error {
	return nil
}

func (noOpClient) CheckUsersExist(_ context.Context, _ []int64) error {
	return nil
}

// NewMockImplementation creates a new mock instance of Implementation with the given chat service and no op auth client.
func NewMockImplementation(deps ...any) *Implementation {
	impl := &Implementation{
		authClient: noOpClient{},
	}

	for _, v := range deps {
		switch s := v.(type) {
		case service.ChatService:
			impl.chatService = s
		}
	}

	return impl
}
