package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// poolInitializer is struct with receiver-methods
// for initial load and prepare saved in database wallets and mnemonics
type poolInitializer struct {
	logger *zap.Logger
	cfg    configService

	walletsDataSrv         walletsDataService
	mnemonicWalletsDataSrv mnemonicWalletsDataService
	txStmtManager          transactionalStatementManager

	encryptSrv encryptService

	wallets        []*entities.Wallet
	walletUUIDList []string
	walletUnits    map[uuid.UUID]WalletPoolUnitService

	mnemonicItems []*entities.MnemonicWallet

	groupedByUUID map[uuid.UUID][]walletPoolMnemonicUnitService
}

func (pi *poolInitializer) LoadAndInitWallets(ctx context.Context) error {
	err := pi.loadWallets(ctx)
	if err != nil {
		return err
	}

	err = pi.prepareMnemonicUnits(ctx)
	if err != nil {
		return err
	}

	err = pi.prepareWalletPoolUnits(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (pi *poolInitializer) loadWallets(ctx context.Context) error {
	wallets, err := pi.walletsDataSrv.GetAllEnabledWallets(ctx)
	if err != nil {
		return err
	}

	if wallets == nil {
		return nil
	}

	walletsUUIDList := make([]string, len(wallets))
	for i := 0; i != len(wallets); i++ {
		walletsUUIDList[i] = wallets[i].UUID.String()
	}

	mnemonicWallets, err := pi.mnemonicWalletsDataSrv.GetMnemonicWalletsByUUIDList(ctx, walletsUUIDList)
	if err != nil {
		return err
	}

	pi.walletUUIDList = walletsUUIDList
	pi.mnemonicItems = mnemonicWallets

	return nil
}

func (pi *poolInitializer) prepareMnemonicUnits(_ context.Context) error {
	groupedByWalletUUID := make(map[uuid.UUID][]walletPoolMnemonicUnitService, len(pi.walletUUIDList))

	// prepare map of mnemonic units grouped by wallet uuid
	for i, _ := range pi.mnemonicItems {
		mnemonicWallet := pi.mnemonicItems[i]
		walletUUID := mnemonicWallet.WalletUUID

		bucket, isExists := pi.groupedByUUID[walletUUID]
		if !isExists {
			pi.groupedByUUID[walletUUID] = make([]walletPoolMnemonicUnitService, 0)
			bucket = pi.groupedByUUID[walletUUID]
		}

		mnemonicUnit := newMnemonicWalletPoolUnit(pi.logger, pi.cfg, mnemonicWallet.UnloadInterval,
			walletUUID, pi.encryptSrv, pi.mnemonicWalletsDataSrv, mnemonicWallet)

		bucket = append(bucket, mnemonicUnit)
	}

	pi.groupedByUUID = groupedByWalletUUID

	return nil
}

// prepareWalletPoolUnits is method for create and save in memory wallet pool units
func (pi *poolInitializer) prepareWalletPoolUnits(_ context.Context) error {
	for i := 0; i != len(pi.wallets); i++ {
		walletItem := pi.wallets[i]
		bucket := pi.groupedByUUID[walletItem.UUID]
		mnemonicCount := len(bucket)

		if mnemonicCount == 1 && walletItem.Strategy == types.WalletMakerSingleMnemonicStrategy {
			walletPoolUnit := newSingleMnemonicWalletPoolUnit(pi.logger, walletItem.UUID,
				walletItem.Title, walletItem.Purpose)
			addErr := walletPoolUnit.AddMnemonicUnit(bucket[0])
			if addErr != nil {
				return addErr
			}

			pi.walletUnits[walletItem.UUID] = walletPoolUnit

			continue
		}

		if mnemonicCount > 1 && walletItem.Strategy == types.WalletMakerMultipleMnemonicStrategy {
			walletPoolUnit := newMultipleMnemonicWalletPoolUnit(pi.logger, walletItem.UUID,
				walletItem.Title, walletItem.Purpose)
			for j := 0; j != mnemonicCount; j++ {
				addErr := walletPoolUnit.AddMnemonicUnit(bucket[i])
				if addErr != nil {
					return addErr
				}
			}

			pi.walletUnits[walletItem.UUID] = walletPoolUnit
		}
	}

	return nil
}

func (pi *poolInitializer) GetWalletPoolUnits() map[uuid.UUID]WalletPoolUnitService {
	return pi.walletUnits
}

func newWalletPoolInitializer(logger *zap.Logger,
	cfg configService,
	encryptSrv encryptService,
	walletDataSrv walletsDataService,
	mnemonicWalletDataSrv mnemonicWalletsDataService,
	txStmtManager transactionalStatementManager,
) (*poolInitializer, error) {
	return &poolInitializer{
		logger: logger,
		cfg:    cfg,

		txStmtManager:          txStmtManager,
		walletsDataSrv:         walletDataSrv,
		mnemonicWalletsDataSrv: mnemonicWalletDataSrv,

		encryptSrv: encryptSrv,
	}, nil
}
