package client

import "context"

// AuthClient defines the interface for interacting with the gRPC authentication service.
type AuthClient interface {
	CheckAccess(ctx context.Context, endpoint string) error
	CheckUsersExist(ctx context.Context, ids []int64) error
}
