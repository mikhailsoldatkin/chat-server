package auth

import (
	"context"

	pbAccess "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	pbUser "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"github.com/mikhailsoldatkin/chat-server/internal/client"
)

var _ client.AuthClient = (*authClient)(nil)

type authClient struct {
	accessClient pbAccess.AccessV1Client
	userClient   pbUser.UserV1Client
}

// NewAuthClient creates a new instance of AuthClient with the provided access and user clients.
func NewAuthClient(
	accessClient pbAccess.AccessV1Client,
	userClient pbUser.UserV1Client,
) client.AuthClient {
	return &authClient{
		accessClient: accessClient,
		userClient:   userClient,
	}
}

// CheckAccess verifies if the specified endpoint has the required access permissions.
// It returns an error if access is denied.
func (cl *authClient) CheckAccess(ctx context.Context, endpoint string) error {
	_, err := cl.accessClient.Check(ctx, &pbAccess.CheckRequest{Endpoint: endpoint})
	if err != nil {
		return err
	}
	return nil
}

// CheckUsersExist checks a users from the given id list exist.
// It returns an error if the users does not exist or if the operation fails.
func (cl *authClient) CheckUsersExist(ctx context.Context, ids []int64) error {
	_, err := cl.userClient.CheckUsersExist(ctx, &pbUser.CheckUsersExistRequest{Ids: ids})
	if err != nil {
		return err
	}
	return nil
}
