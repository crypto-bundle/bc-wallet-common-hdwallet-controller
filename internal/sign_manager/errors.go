package sign_manager

import "errors"

var (
	ErrMissingHdWalletResp         = errors.New("missing hd-wallet api response")
	ErrMissingDerivationPathField  = errors.New("missing derivationPath field sign in request item")
	ErrUnableDecodeGrpcErrorStatus = errors.New("unable to decode grpc error status")
)
