package wallet_manager

import "errors"

var (
	ErrMissingHdWalletResp         = errors.New("missing hd-wallet api response")
	ErrUnableDecodeGrpcErrorStatus = errors.New("unable to decode grpc error status")
)
