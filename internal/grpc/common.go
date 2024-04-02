package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"
	"github.com/google/uuid"
)

type configService interface {
	IsDev() bool
	IsDebug() bool
	IsLocal() bool

	GetBindPort() string

	GetProviderName() string
	GetNetworkName() string
}

type mnemonicWalletsDataService interface {
	AddNewMnemonicWallet(ctx context.Context, wallet *entities.MnemonicWallet) (*entities.MnemonicWallet, error)
	AddNewMnemonicWallets(ctx context.Context,
		walletsList []*entities.MnemonicWallet,
	) (uint32, []*entities.MnemonicWallet, error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletUUID(ctx context.Context, walletUUID uuid.UUID) (*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDList(ctx context.Context,
		UUIDList []string,
	) ([]*entities.MnemonicWallet, error)
	GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetAllNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
}

type walletSessionDataService interface {
	AddNewWalletSession(ctx context.Context,
		sessionItem *entities.MnemonicWalletSession,
	) (*entities.MnemonicWalletSession, error)
	GetAllActiveSessions(ctx context.Context,
	) (uint, []*entities.MnemonicWalletSession, error)
	GetAllActiveSessionsClb(ctx context.Context,
		onItemCallBack func(item *entities.MnemonicWalletSession) error,
	) (uint, []*entities.MnemonicWalletSession, error)
	GetWalletActiveSessions(ctx context.Context,
		walletUUID string,
	) (uint, []*entities.MnemonicWalletSession, error)
	GetWalletActiveSessionsClb(ctx context.Context,
		walletUUID string,
		onItemCallBack func(item *entities.MnemonicWalletSession) error,
	) (uint, []*entities.MnemonicWalletSession, error)
	GetWalletSessionByUUID(ctx context.Context,
		sessionUUID string,
	) (*entities.MnemonicWalletSession, error)
}

type walletManagerService interface {
	CreateNewWallet(ctx context.Context,
		title string,
		purpose string,
	) (*types.PublicWalletData, error)
	DisableWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*types.PublicWalletData, error)
	DisableWalletsByUUIDList(ctx context.Context,
		walletUUID []string,
	) (uint, []string, error)
	EnableWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*types.PublicWalletData, error)
	GetAddressByPath(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
		account, change, index uint32,
	) (*types.PublicDerivationAddressData, error)

	GetAddressesByPathByRange(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error

	GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*entities.MnemonicWallet, error)
	GetEnabledWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)

	SignTransaction(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
		transactionData []byte,
	) (*types.PublicSignTxData, error)
}

type marshallerService interface {
	MarshallCreateWalletData(*types.PublicWalletData) (*pbApi.AddNewWalletResponse, error)
	MarshallGetAddressData(
		walletPublicData *entities.MnemonicWallet,
		addressPublicData *pbCommon.DerivationAddressIdentity,
	) (*pbApi.DerivationAddressResponse, error)
	MarshallGetAddressByRange(
		walletPublicData *types.PublicWalletData,
		mnemonicWalletPublicData *types.PublicMnemonicWalletData,
		addressesData []*pbCommon.DerivationAddressIdentity,
		size uint64,
	) (*pbApi.DerivationAddressByRangeResponse, error)
	MarshallGetEnabledWallets([]*entities.MnemonicWallet) (*pbApi.GetEnabledWalletsResponse, error)
	MarshallSignTransaction(
		publicSignTxData *types.PublicSignTxData,
	) (*pbApi.SignTransactionResponse, error)
	MarshallWalletInfo(
		walletData *types.PublicWalletData,
	) *pbCommon.WalletData
}

type addWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.AddNewWalletRequest) (*pbApi.AddNewWalletResponse, error)
}

type disableWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.DisableWalletRequest) (*pbApi.DisableWalletResponse, error)
}

type disableWalletsHandlerService interface {
	Handle(ctx context.Context, request *pbApi.DisableWalletsRequest) (*pbApi.DisableWalletsResponse, error)
}

type enableWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.EnableWalletRequest) (*pbApi.EnableWalletResponse, error)
}

type getWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetWalletInfoRequest) (*pbApi.GetWalletInfoResponse, error)
}

type getAddressHandlerService interface {
	Handle(ctx context.Context, request *pbApi.DerivationAddressRequest) (*pbApi.DerivationAddressResponse, error)
}

type getAddressByRangeHandlerService interface {
	Handle(ctx context.Context,
		request *pbApi.DerivationAddressByRangeRequest,
	) (*pbApi.DerivationAddressByRangeResponse, error)
}

type signTransactionRequestHandlerService interface {
	Handle(ctx context.Context, request *pbApi.SignTransactionRequest) (*pbApi.SignTransactionResponse, error)
}

type getEnabledWalletsHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetEnabledWalletsRequest) (*pbApi.GetEnabledWalletsResponse, error)
}

type startWalletSessionHandlerService interface {
	Handle(ctx context.Context, request *pbApi.StartWalletSessionRequest) (*pbApi.StartWalletSessionResponse, error)
}

type getWalletSessionHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetWalletSessionRequest) (*pbApi.GetWalletSessionResponse, error)
}

type getWalletSessionsHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetWalletSessionsRequest) (*pbApi.GetWalletSessionsResponse, error)
}
