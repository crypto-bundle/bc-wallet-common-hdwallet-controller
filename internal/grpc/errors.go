package grpc

import "errors"

var (
	ErrMissedRequiredData              = errors.New("missed required data")
	ErrDataIsNotValid                  = errors.New("received data is not valid")
	ErrUnableReadGrpcMetadata          = errors.New("unable to read grpc metadata")
	ErrUnableGetWalletUUIDFromMetadata = errors.New("unable to get wallet uuid from metadata")
)
