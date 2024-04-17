package wallet_manager

import "errors"

var (
	ErrMissingHdWalletResp         = errors.New("missing hd-wallet api response")
	ErrMnemonicIsNotValid          = errors.New("mnemonic wallet is not valid")
	ErrUnableDecodeGrpcErrorStatus = errors.New("unable to decode grpc error status")
)
