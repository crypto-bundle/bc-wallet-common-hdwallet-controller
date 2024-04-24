package wallet_manager

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type singleMnemonicWalletUnit struct {
	logger *zap.Logger

	walletUUID    uuid.UUID
	walletTitle   string
	walletPurpose string

	cfgSrv         configService
	cryptoSrv      encryptService
	walletsDataSrv walletsDataService

	mnemonicWalletsDataSrv mnemonicWalletsDataService

	mnemonicUnit walletPoolMnemonicUnitService
}

func (u *singleMnemonicWalletUnit) Init(ctx context.Context) error {
	return u.mnemonicUnit.Init(ctx)
}

func (u *singleMnemonicWalletUnit) Run(ctx context.Context) error {
	return u.mnemonicUnit.Run(ctx)
}

func (u *singleMnemonicWalletUnit) Shutdown(ctx context.Context) error {
	return u.mnemonicUnit.Shutdown(ctx)
}

func (u *singleMnemonicWalletUnit) GetWalletUUID() uuid.UUID {
	return u.walletUUID
}

func (u *singleMnemonicWalletUnit) GetWalletTitle() string {
	return u.walletTitle
}

func (u *singleMnemonicWalletUnit) GetWalletPurpose() string {
	return u.walletPurpose
}

func (u *singleMnemonicWalletUnit) GetWalletPublicData() *types.PublicWalletData {
	mnemonicPubData := u.mnemonicUnit.GetPublicData()
	return &types.PublicWalletData{
		UUID:     u.walletUUID,
		Title:    u.walletTitle,
		Purpose:  u.walletPurpose,
		Strategy: types.WalletMakerSingleMnemonicStrategy,
		MnemonicWallets: []*types.PublicMnemonicWalletData{
			mnemonicPubData,
		},
		MnemonicWalletsByUUID: map[uuid.UUID]*types.PublicMnemonicWalletData{
			mnemonicPubData.UUID: mnemonicPubData,
		},
	}
}

func (u *singleMnemonicWalletUnit) SignTransaction(ctx context.Context,
	mnemonicUUID uuid.UUID,
	account, change, index uint32,
	transactionData []byte,
) (*types.PublicSignTxData, error) {
	if mnemonicUUID != u.mnemonicUnit.GetMnemonicUUID() {
		return nil, ErrPassedMnemonicWalletNotFound
	}

	return u.mnemonicUnit.SignTransaction(ctx, account, change, index, transactionData)
}

func (u *singleMnemonicWalletUnit) AddMnemonicUnit(unit walletPoolMnemonicUnitService) error {
	if u.mnemonicUnit != nil {
		return ErrMnemonicAlreadySet
	}

	u.mnemonicUnit = unit

	return nil
}

func (u *singleMnemonicWalletUnit) GetAddressByPath(ctx context.Context,
	mnemonicUUID uuid.UUID,
	account, change, index uint32,
) (string, error) {
	if mnemonicUUID != u.mnemonicUnit.GetMnemonicUUID() {
		return "", ErrPassedMnemonicWalletNotFound
	}

	return u.mnemonicUnit.GetAddressByPath(ctx, account, change, index)
}

func (u *singleMnemonicWalletUnit) GetAddressesByPathByRange(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
	rangeIterable types.AddrRangeIterable,
	marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
) error {
	if mnemonicWalletUUID != u.mnemonicUnit.GetMnemonicUUID() {
		return ErrPassedMnemonicWalletNotFound
	}

	return u.mnemonicUnit.GetAddressesByPathByRange(ctx, rangeIterable, marshallerCallback)
}

func newSingleMnemonicWalletPoolUnit(logger *zap.Logger,
	walletUUID uuid.UUID,
	walletTitle string,
	walletPurpose string,
) *singleMnemonicWalletUnit {
	return &singleMnemonicWalletUnit{
		logger:        logger.With(zap.String(app.WalletUUIDTag, walletUUID.String())),
		walletUUID:    walletUUID,
		walletTitle:   walletTitle,
		walletPurpose: walletPurpose,
	}
}
