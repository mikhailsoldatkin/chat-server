package grpc_client

import (
	"context"
	"log"

	pbAccess "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	pbUser "github.com/mikhailsoldatkin/auth/pkg/user_v1"
	"github.com/mikhailsoldatkin/platform_common/pkg/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// AuthClient defines the interface for interacting with the gRPC Authentication service.
type AuthClient interface {
	// CheckAccess verifies if the specified endpoint has the required access permissions.
	// It returns an error if access is denied.
	CheckAccess(ctx context.Context, endpoint string) error

	// CheckUsersExist checks a users from the given id list exist.
	// It returns an error if the users does not exist or if the operation fails.
	CheckUsersExist(ctx context.Context, ids []int64) error
}

type authClientImplementation struct {
	accessClient pbAccess.AccessV1Client
	userClient   pbUser.UserV1Client
}

// NewAuthClient creates a new instance of AuthClient with the provided server address and TLS certificate.
func NewAuthClient(address, certFilePath string) (AuthClient, error) {
	creds, err := credentials.NewClientTLSFromFile(certFilePath, "")
	if err != nil {
		log.Fatalf("failed to load TLS credentials from file: %v", err)
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to create connection: %v", err)
	}
	closer.Add(conn.Close)

	return &authClientImplementation{
		accessClient: pbAccess.NewAccessV1Client(conn),
		userClient:   pbUser.NewUserV1Client(conn),
	}, nil
}

func (cl *authClientImplementation) CheckAccess(ctx context.Context, endpoint string) error {
	_, err := cl.accessClient.Check(ctx, &pbAccess.CheckRequest{Endpoint: endpoint})
	if err != nil {
		return err
	}
	return nil
}

func (cl *authClientImplementation) CheckUsersExist(ctx context.Context, ids []int64) error {
	_, err := cl.userClient.CheckUsersExist(ctx, &pbUser.CheckUsersExistRequest{Ids: ids})
	if err != nil {
		return err
	}
	return nil
}
