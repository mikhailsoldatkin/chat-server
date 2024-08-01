package customerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ConvertError converts an error into a gRPC error with the appropriate status code.
func ConvertError(err error) error {
	var notFoundErr *NotFoundError
	var userNotInChatErr *UserNotInChatError

	switch {
	case errors.As(err, &notFoundErr):
		return status.Errorf(codes.NotFound, notFoundErr.Error())
	case errors.As(err, &userNotInChatErr):
		return status.Errorf(codes.NotFound, userNotInChatErr.Error())
	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
