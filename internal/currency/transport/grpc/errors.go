package grpc

import (
	"Currency-apiNew2/internal/currency/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch err {
	case domain.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrInvalidArg:
		return status.Error(codes.InvalidArgument, err.Error())
	case domain.ErrServiceNotReady:
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
