package grpc

import "errors"

var (
	ErrMissedRequiredData              = errors.New("missed required data")
	ErrUnableReadGrpcMetadata          = errors.New("unable to read grpc metadata")
	ErrUnableGetWalletUUIDFromMetadata = errors.New("unable to get wallet uuid from metadata")
)
