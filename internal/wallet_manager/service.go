package wallet_manager

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	cfg    configService

	walletsDataSrv           walletsDataService
	mnemonicWalletsDataSrv   mnemonicWalletsDataService
	mnemonicGeneratorSrv     mnemonicGenerator
	encryptSrv               encryptService
	walletPoolSrv            walletPoolService
	walletPoolInitializerSrv walletPoolInitService
	walletPoolUnitMakerSrv   walletPoolUnitMakerService

	txStmtManager transactionalStatementManager
}

func (s *Service) Init(ctx context.Context) error {
	err := s.walletPoolInitializerSrv.LoadAndInitWallets(ctx)
	if err != nil {
		return err
	}

	loadedWallets := s.walletPoolInitializerSrv.GetWalletPoolUnits()
	if loadedWallets != nil {
		err = s.walletPoolSrv.SetWalletUnits(ctx, loadedWallets)
		if err != nil {
			return err
		}
	}

	err = s.walletPoolSrv.Init(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Run(ctx context.Context) error {
	err := s.walletPoolSrv.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Shutdown(ctx context.Context) error {
	err := s.walletPoolSrv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*types.PublicWalletData, error) {
	return s.walletPoolSrv.GetWalletByUUID(ctx, walletUUID)
}

func (s *Service) GetEnabledWallets(ctx context.Context) ([]*types.PublicWalletData, error) {
	return s.walletPoolSrv.GetEnabledWallets(ctx)
}

func (s *Service) GetAddressByPath(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicWalletUUID uuid.UUID,
	account, change, index uint32,
) (*types.PublicDerivationAddressData, error) {
	address, err := s.walletPoolSrv.GetAddressByPath(ctx, walletUUID, mnemonicWalletUUID, account, change, index)
	if err != nil {
		return nil, err
	}

	return &types.PublicDerivationAddressData{
		AccountIndex:  account,
		InternalIndex: change,
		AddressIndex:  index,
		Address:       address,
	}, nil
}

func (s *Service) GetAddressesByPathByRange(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicWalletUUID uuid.UUID,
	rangeIterable types.AddrRangeIterable,
	marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
) error {
	return s.walletPoolSrv.GetAddressesByPathByRange(ctx, walletUUID, mnemonicWalletUUID,
		rangeIterable, marshallerCallback)
}

func (s *Service) CreateNewWallet(ctx context.Context,
	strategy types.WalletMakerStrategy,
	title string,
	purpose string,
) (*types.PublicWalletData, error) {
	poolUnit, err := s.walletPoolUnitMakerSrv.CreateDisabledWallet(ctx, strategy, title, purpose)
	if err != nil {
		return nil, err
	}

	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		err = s.walletPoolSrv.AddAndStartWalletUnit(txStmtCtx, poolUnit.GetWalletUUID(), poolUnit)
		if err != nil {
			return err
		}

		err = s.walletsDataSrv.SetEnabledToWalletByUUID(txStmtCtx, poolUnit.GetWalletUUID().String())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return poolUnit.GetWalletPublicData(), nil
}

func (s *Service) SignTransaction(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicUUID uuid.UUID,
	account, change, index uint32,
	transactionData []byte,
) (*types.PublicSignTxData, error) {
	return s.walletPoolSrv.SignTransaction(ctx, walletUUID, mnemonicUUID,
		account, change, index, transactionData)
}

func NewService(logger *zap.Logger,
	cfg configService,
	encryptSrv encryptService,
	walletDataSrv walletsDataService,
	mnemonicWalletDataSrv mnemonicWalletsDataService,
	txStmtManager transactionalStatementManager,
	mnemonicGeneratorSrv mnemonicGenerator,
) (*Service, error) {
	walletPoolInitSrv, err := newWalletPoolInitializer(logger, cfg,
		encryptSrv, walletDataSrv, mnemonicWalletDataSrv,
		txStmtManager)
	if err != nil {
		return nil, err
	}

	walletPoolSrv := newWalletPool(logger, cfg, walletDataSrv, mnemonicWalletDataSrv, encryptSrv)
	walletMaker := newWalletMaker(logger, cfg, walletDataSrv, mnemonicWalletDataSrv,
		txStmtManager, mnemonicGeneratorSrv, encryptSrv)

	return &Service{
		logger: logger,
		cfg:    cfg,

		txStmtManager:          txStmtManager,
		walletsDataSrv:         walletDataSrv,
		mnemonicWalletsDataSrv: mnemonicWalletDataSrv,

		mnemonicGeneratorSrv: mnemonicGeneratorSrv,
		encryptSrv:           encryptSrv,

		walletPoolSrv:            walletPoolSrv,
		walletPoolInitializerSrv: walletPoolInitSrv,
		walletPoolUnitMakerSrv:   walletMaker,
	}, nil
}
