package grpc_hdwallet

import (
	"context"

	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

const (
	MethodNameTag = "method_name"
)

type configService interface {
	IsDev() bool
	IsDebug() bool
	IsLocal() bool

	GetBindPort() string
}

type generateMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.GenerateMnemonicRequest) (*pbApi.GenerateMnemonicResponse, error)
}

type loadMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.LoadMnemonicRequest) (*pbApi.LoadMnemonicResponse, error)
}

type unLoadMnemonicHandlerService interface {
	Handle(ctx context.Context, req *pbApi.UnLoadMnemonicRequest) (*pbApi.UnLoadMnemonicResponse, error)
}

type getDerivationsAddressesHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.DerivationAddressByRangeRequest,
	) (*pbApi.DerivationAddressByRangeResponse, error)
}

type signTransactionHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.SignTransactionRequest,
	) (*pbApi.SignTransactionResponse, error)
}

type getDerivationAddressHandlerService interface {
	Handle(ctx context.Context,
		req *pbApi.DerivationAddressRequest,
	) (*pbApi.DerivationAddressResponse, error)
}
