package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	tronCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/google/uuid"

	"go.uber.org/zap"
)

type multipleMnemonicWalletUnit struct {
	logger *zap.Logger

	walletUUID    uuid.UUID
	walletTitle   string
	walletPurpose string

	cfgSrv         configService
	encryptSrv     encryptService
	walletsDataSrv walletsDataService

	mnemonicWalletsDataSrv mnemonicWalletsDataService

	hotMnemonicUnit     walletPoolMnemonicUnitService
	mnemonicUnitsCount  uint8
	mnemonicUnits       []walletPoolMnemonicUnitService
	mnemonicUnitsByUUID map[uuid.UUID]walletPoolMnemonicUnitService
}

func (u *multipleMnemonicWalletUnit) Init(ctx context.Context) error {
	for _, walletUnit := range u.mnemonicUnits {
		initErr := walletUnit.Init(ctx)
		if initErr != nil {
			return initErr
		}
	}

	return nil
}

func (u *multipleMnemonicWalletUnit) Run(ctx context.Context) error {
	for _, walletUnit := range u.mnemonicUnits {
		runErr := walletUnit.Run(ctx)
		if runErr != nil {
			return runErr
		}
	}

	return nil
}

func (u *multipleMnemonicWalletUnit) Shutdown(ctx context.Context) error {
	for _, walletUnit := range u.mnemonicUnits {
		err := walletUnit.Shutdown(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *multipleMnemonicWalletUnit) GetWalletUUID() uuid.UUID {
	return u.walletUUID
}

func (u *multipleMnemonicWalletUnit) GetWalletTitle() string {
	return u.walletTitle
}

func (u *multipleMnemonicWalletUnit) GetWalletPurpose() string {
	return u.walletPurpose
}

func (u *multipleMnemonicWalletUnit) GetWalletPublicData() *types.PublicWalletData {
	publicData := &types.PublicWalletData{
		UUID:            u.walletUUID,
		Title:           u.walletTitle,
		Purpose:         u.walletPurpose,
		Strategy:        types.WalletMakerSingleMnemonicStrategy,
		MnemonicWallets: make([]*types.PublicMnemonicWalletData, u.mnemonicUnitsCount),
	}

	for i := uint8(0); i != u.mnemonicUnitsCount; i++ {
		publicData.MnemonicWallets[i] = u.mnemonicUnits[i].GetPublicData()
	}

	return publicData
}

func (u *multipleMnemonicWalletUnit) AddMnemonicUnit(unit walletPoolMnemonicUnitService) error {
	u.mnemonicUnits = append(u.mnemonicUnits, unit)

	u.mnemonicUnitsByUUID[unit.GetMnemonicUUID()] = unit

	if unit.IsHotWalletUnit() {
		u.hotMnemonicUnit = unit
	}

	u.mnemonicUnitsCount++

	return nil
}

func (u *multipleMnemonicWalletUnit) GetAddressByPath(ctx context.Context,
	mnemonicUUID uuid.UUID,
	account, change, index uint32,
) (string, error) {
	mnemonicUnit, isExists := u.mnemonicUnitsByUUID[mnemonicUUID]
	if !isExists {
		return "", ErrPassedMnemonicWalletNotFound
	}

	return mnemonicUnit.GetAddressByPath(ctx, account, change, index)
}

func (u *multipleMnemonicWalletUnit) GetAddressesByPathByRange(ctx context.Context,
	mnemonicWalletUUID uuid.UUID,
	accountIndex uint32,
	internalIndex uint32,
	addressIndexFrom uint32,
	addressIndexTo uint32,
) ([]*types.PublicDerivationAddressData, error) {
	mnemonicUnit, isExists := u.mnemonicUnitsByUUID[mnemonicWalletUUID]
	if !isExists {
		return nil, ErrPassedMnemonicWalletNotFound
	}

	return mnemonicUnit.GetAddressesByPathByRange(ctx, accountIndex, internalIndex, addressIndexFrom, addressIndexTo)
}

func (u *multipleMnemonicWalletUnit) SignTransaction(ctx context.Context,
	mnemonicUUID uuid.UUID,
	account, change, index uint32,
	transaction *tronCore.Transaction,
) (*types.PublicSignTxData, error) {
	mnemonicUnit, isExists := u.mnemonicUnitsByUUID[mnemonicUUID]
	if !isExists {
		return nil, ErrPassedMnemonicWalletNotFound
	}

	return mnemonicUnit.SignTransaction(ctx, account, change, index, transaction)
}

func newMultipleMnemonicWalletPoolUnit(logger *zap.Logger,
	cfg configService,
	encryptionSrv encryptService,
	walletDataSrv walletsDataService,
	mnemonicWalletDataSrv mnemonicWalletsDataService,
	walletUUID uuid.UUID,
	walletTitle string,
	walletPurpose string,
) *multipleMnemonicWalletUnit {
	return &multipleMnemonicWalletUnit{
		logger: logger.With(zap.String(app.WalletUUIDTag, walletUUID.String())),
		cfgSrv: cfg,

		walletUUID:    walletUUID,
		walletTitle:   walletTitle,
		walletPurpose: walletPurpose,

		encryptSrv:             encryptionSrv,
		walletsDataSrv:         walletDataSrv,
		mnemonicWalletsDataSrv: mnemonicWalletDataSrv,

		mnemonicUnitsByUUID: map[uuid.UUID]walletPoolMnemonicUnitService{},
		mnemonicUnits:       make([]walletPoolMnemonicUnitService, 0),
		mnemonicUnitsCount:  0,
	}
}
