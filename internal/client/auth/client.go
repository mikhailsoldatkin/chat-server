package auth

import (
	"context"

	pbAccess "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	pbUser "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"github.com/mikhailsoldatkin/chat-server/internal/client"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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

func (cl *authClient) CheckAccess(ctx context.Context, endpoint string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "auth-check-endpoint-access")
	defer span.Finish()

	span.SetTag("endpoint", endpoint)

	_, err := cl.accessClient.Check(ctx, &pbAccess.CheckRequest{Endpoint: endpoint})
	if err != nil {
		span.SetTag("error", true)
		return errors.WithMessage(err, "checking endpoint access")
	}
	return nil
}

func (cl *authClient) CheckUsersExist(ctx context.Context, ids []int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "auth-check-users-exist")
	defer span.Finish()

	span.SetTag("user-ids", ids)

	_, err := cl.userClient.CheckUsersExist(ctx, &pbUser.CheckUsersExistRequest{Ids: ids})
	if err != nil {
		span.SetTag("error", true)
		return errors.WithMessage(err, "checking users existence")
	}
	return nil
}
