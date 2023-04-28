package grpc

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
	"github.com/google/uuid"
)

type configService interface {
	IsDev() bool
	IsDebug() bool
	IsLocal() bool

	GetBindPort() string
}

type walletManagerService interface {
	CreateNewWallet(ctx context.Context,
		strategy types.WalletMakerStrategy,
		title string,
		purpose string,
	) (*types.PublicWalletData, error)
	GetAddressByPath(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		account, change, index uint32,
	) (*types.PublicDerivationAddressData, error)

	GetAddressesByPathByRange(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error

	GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*types.PublicWalletData, error)
	GetEnabledWallets(ctx context.Context) ([]*types.PublicWalletData, error)

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
		walletPublicData *types.PublicWalletData,
		mnemonicWalletPublicData *types.PublicMnemonicWalletData,
		addressPublicData *pbApi.DerivationAddressIdentity,
	) (*pbApi.DerivationAddressResponse, error)
	MarshallGetAddressByRange(
		walletPublicData *types.PublicWalletData,
		mnemonicWalletPublicData *types.PublicMnemonicWalletData,
		addressesData []*pbApi.DerivationAddressIdentity,
		size uint64,
	) (*pbApi.DerivationAddressByRangeResponse, error)
	MarshallGetEnabledWallets([]*types.PublicWalletData) (*pbApi.GetEnabledWalletsResponse, error)
	MarshallSignTransaction(
		publicSignTxData *types.PublicSignTxData,
	) (*pbApi.SignTransactionResponse, error)
	MarshallWalletInfo(
		walletData *types.PublicWalletData,
	) *pbApi.WalletData
}
