package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"

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
	AddNewWallet(ctx context.Context) (*entities.MnemonicWallet, error)
	ImportWallet(ctx context.Context, mnemonicData []byte) (*entities.MnemonicWallet, error)
	EnableWalletByUUID(ctx context.Context,
		walletUUID string,
	) (*entities.MnemonicWallet, error)
	DisableWalletByUUID(ctx context.Context,
		walletUUID string,
	) (*entities.MnemonicWallet, error)
	DisableWalletsByUUIDList(ctx context.Context,
		walletUUIDs []string,
	) (count uint, list []string, err error)
	GetEnabledWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetWalletByUUID(ctx context.Context, walletUUID string) (*entities.MnemonicWallet, error)
	GetAddress(ctx context.Context,
		mnemonicUUID string,
		account, change, index uint32,
		sessionUUID string,
	) (ownerWallet *entities.MnemonicWallet, address *string, err error)
	GetAddressesByRange(ctx context.Context,
		mnemonicUUID string,
		sessionUUID string,
		addrRanges []*pbCommon.RangeRequestUnit,
	) (ownerWallet *entities.MnemonicWallet, list []*pbCommon.DerivationAddressIdentity, err error)

	StartWalletSession(ctx context.Context,
		walletUUID string,
	) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error)
	CloseWalletSession(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error)
	GetWalletSessionInfo(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error)
	GetWalletSessionsByWalletUUID(ctx context.Context,
		walletUUID string,
	) (wallet *entities.MnemonicWallet, list []*entities.MnemonicWalletSession, err error)
}

type signManagerService interface {
	GetActiveSignRequest(ctx context.Context,
		signUUID string,
	) (signReqItem *entities.SignRequest, err error)
	PrepareSignRequest(ctx context.Context,
		mnemonicUUID string,
		purposeUUID string,
		account, change, index uint32,
	) (signerAddr *pbCommon.DerivationAddressIdentity, request *entities.SignRequest, err error)
	ExecuteSignRequest(ctx context.Context,
		signReqItem *entities.SignRequest,
		transactionData []byte,
	) (signerAddr *pbCommon.DerivationAddressIdentity, signedData []byte, err error)
	CloseSignRequest(ctx context.Context,
		signReqUUID string,
	) (*entities.SignRequest, error)
	CloseSignRequestBySession(ctx context.Context,
		sessionUUID string,
	) (count uint, list []*entities.SignRequest, err error)
	CloseSignRequestByWallet(ctx context.Context,
		walletUUID string,
	) (count uint, list []*entities.SignRequest, err error)
	CloseSignRequestByMultipleWallets(ctx context.Context,
		walletUUIDs []string,
	) (uint, []*entities.SignRequest, error)
}

type marshallerService interface {
	MarshallCreateWalletData(wallet *entities.MnemonicWallet) (*pbApi.AddNewWalletResponse, error)
	MarshallGetAddressData(
		walletPublicData *entities.MnemonicWallet,
		addressPublicData *pbCommon.DerivationAddressIdentity,
	) (*pbApi.DerivationAddressResponse, error)
	MarshallGetAddressByRange(
		walletPublicData *entities.MnemonicWallet,
		addressesData []*pbCommon.DerivationAddressIdentity,
		size uint64,
	) (*pbApi.DerivationAddressByRangeResponse, error)
	MarshallGetEnabledWallets([]*entities.MnemonicWallet) (*pbApi.GetEnabledWalletsResponse, error)
	MarshallWalletSessions(
		sessionsList []*entities.MnemonicWalletSession,
	) []*pbApi.SessionInfo
}

type addWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.AddNewWalletRequest) (*pbApi.AddNewWalletResponse, error)
}

type importWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.ImportWalletRequest) (*pbApi.ImportWalletResponse, error)
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

type prepareSignRequestHandlerService interface {
	Handle(ctx context.Context, request *pbApi.PrepareSignRequestReq) (*pbApi.PrepareSignRequestResponse, error)
}

type executeSignRequestHandlerService interface {
	Handle(ctx context.Context, request *pbApi.ExecuteSignRequestReq) (*pbApi.ExecuteSignRequestResponse, error)
}

type getEnabledWalletsHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetEnabledWalletsRequest) (*pbApi.GetEnabledWalletsResponse, error)
}

type startWalletSessionHandlerService interface {
	Handle(ctx context.Context, request *pbApi.StartWalletSessionRequest) (*pbApi.StartWalletSessionResponse, error)
}

type closeWalletSessionHandlerService interface {
	Handle(ctx context.Context, request *pbApi.CloseWalletSessionsRequest) (*pbApi.CloseWalletSessionsResponse, error)
}

type getWalletSessionHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetWalletSessionRequest) (*pbApi.GetWalletSessionResponse, error)
}

type getWalletSessionsHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetWalletSessionsRequest) (*pbApi.GetWalletSessionsResponse, error)
}
