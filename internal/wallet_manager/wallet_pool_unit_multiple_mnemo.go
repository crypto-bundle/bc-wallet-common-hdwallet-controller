package wallet_manager

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

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
		UUID:                  u.walletUUID,
		Title:                 u.walletTitle,
		Purpose:               u.walletPurpose,
		Strategy:              types.WalletMakerMultipleMnemonicStrategy,
		MnemonicWallets:       make([]*types.PublicMnemonicWalletData, u.mnemonicUnitsCount),
		MnemonicWalletsByUUID: make(map[uuid.UUID]*types.PublicMnemonicWalletData, u.mnemonicUnitsCount),
	}

	for i := uint8(0); i != u.mnemonicUnitsCount; i++ {
		mnemonicPubData := u.mnemonicUnits[i].GetPublicData()

		publicData.MnemonicWallets[i] = mnemonicPubData
		publicData.MnemonicWalletsByUUID[mnemonicPubData.UUID] = mnemonicPubData
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
	rangeIterable types.AddrRangeIterable,
	marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
) error {
	mnemonicUnit, isExists := u.mnemonicUnitsByUUID[mnemonicWalletUUID]
	if !isExists {
		return ErrPassedMnemonicWalletNotFound
	}

	return mnemonicUnit.GetAddressesByPathByRange(ctx, rangeIterable, marshallerCallback)
}

func (u *multipleMnemonicWalletUnit) SignTransaction(ctx context.Context,
	mnemonicUUID uuid.UUID,
	account, change, index uint32,
	transactionData []byte,
) (*types.PublicSignTxData, error) {
	mnemonicUnit, isExists := u.mnemonicUnitsByUUID[mnemonicUUID]
	if !isExists {
		return nil, ErrPassedMnemonicWalletNotFound
	}

	return mnemonicUnit.SignTransaction(ctx, account, change, index, transactionData)
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
