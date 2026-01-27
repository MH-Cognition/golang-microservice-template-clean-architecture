package apperrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCStatus(err error) error {
	appErr := From(err)

	var code codes.Code
	switch appErr.Code {
	case InvalidInput, InvalidJSON:
		code = codes.InvalidArgument
	case Unauthorized:
		code = codes.Unauthenticated
	case Forbidden:
		code = codes.PermissionDenied
	case NotFound:
		code = codes.NotFound
	case Conflict:
		code = codes.AlreadyExists
	case Timeout:
		code = codes.DeadlineExceeded
	case Dependency:
		code = codes.Unavailable
	default:
		code = codes.Internal
	}

	return status.Error(code, appErr.Message)
}
