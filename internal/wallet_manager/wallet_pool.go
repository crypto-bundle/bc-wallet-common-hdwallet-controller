package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	tronCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Pool struct {
	logger *zap.Logger
	cfg    configService

	walletsDataSrv         walletsDataService
	mnemonicWalletsDataSrv mnemonicWalletsDataService
	encryptSrv             encryptService

	walletUnits map[uuid.UUID]WalletPoolUnitService
}

func (p *Pool) Init(ctx context.Context) error {
	for _, walletUnit := range p.walletUnits {
		initErr := walletUnit.Init(ctx)
		if initErr != nil {
			return initErr
		}
	}

	return nil
}

func (p *Pool) Run(ctx context.Context) error {
	for _, walletUnit := range p.walletUnits {
		initErr := walletUnit.Run(ctx)
		if initErr != nil {
			return initErr
		}
	}

	return nil
}

func (p *Pool) Shutdown(ctx context.Context) error {
	for _, walletUnit := range p.walletUnits {
		initErr := walletUnit.Shutdown(ctx)
		if initErr != nil {
			return initErr
		}
	}

	return nil
}

func (p *Pool) SetWalletUnits(ctx context.Context,
	walletUnits map[uuid.UUID]WalletPoolUnitService,
) error {
	if len(p.walletUnits) > 0 {
		return ErrWalletPoolIsNotEmpty
	}

	if len(walletUnits) == 0 {
		return ErrPassedWalletPoolUnitIsEmpty
	}

	p.walletUnits = walletUnits

	return nil
}

func (p *Pool) AddAWalletUnit(ctx context.Context,
	walletUUID uuid.UUID,
	walletUnit WalletPoolUnitService,
) error {
	_, isExists := p.walletUnits[walletUUID]
	if isExists {
		return ErrPassedWalletAlreadyExists
	}

	p.walletUnits[walletUUID] = walletUnit

	return nil
}

func (p *Pool) AddAndStartWalletUnit(ctx context.Context,
	walletUUID uuid.UUID,
	walletUnit WalletPoolUnitService,
) error {
	_, isExists := p.walletUnits[walletUUID]
	if isExists {
		return ErrPassedWalletAlreadyExists
	}

	p.walletUnits[walletUUID] = walletUnit

	err := walletUnit.Init(ctx)
	if err != nil {
		return err
	}

	err = walletUnit.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pool) GetAddressByPath(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicWalletUUID uuid.UUID,
	account, change, index uint32,
) (string, error) {
	poolUnit, isExists := p.walletUnits[walletUUID]
	if isExists {
		return "", ErrPassedWalletNotFound
	}

	return poolUnit.GetAddressByPath(ctx, mnemonicWalletUUID, account, change, index)
}

func (p *Pool) GetEnabledWallets(ctx context.Context) ([]*types.PublicWalletData, error) {
	result := make([]*types.PublicWalletData, len(p.walletUnits))
	i := 0
	for _, walletUnit := range p.walletUnits {
		result[i] = walletUnit.GetWalletPublicData()
		i++
	}

	return result, nil
}

func (p *Pool) GetAddressesByPathByRange(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicWalletUUID uuid.UUID,
	accountIndex uint32,
	internalIndex uint32,
	addressIndexFrom uint32,
	addressIndexTo uint32,
) ([]*types.PublicDerivationAddressData, error) {
	poolUnit, isExists := p.walletUnits[walletUUID]
	if isExists {
		return nil, ErrPassedWalletNotFound
	}

	return poolUnit.GetAddressesByPathByRange(ctx, mnemonicWalletUUID,
		accountIndex, internalIndex,
		addressIndexFrom, addressIndexTo)
}

func (p *Pool) SignTransaction(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicUUID uuid.UUID,
	account, change, index uint32,
	transaction *tronCore.Transaction,
) (*types.PublicSignTxData, error) {
	poolUnit, isExists := p.walletUnits[walletUUID]
	if isExists {
		return nil, ErrPassedWalletNotFound
	}

	return poolUnit.SignTransaction(ctx, mnemonicUUID, account, change, index, transaction)
}

func newWalletPool(logger *zap.Logger,
	cfg configService,
	walletsDataSrv walletsDataService,
	mnemonicWalletsDataSrv mnemonicWalletsDataService,
	encryptSrv encryptService,
) *Pool {
	return &Pool{
		logger:                 logger,
		cfg:                    cfg,
		walletsDataSrv:         walletsDataSrv,
		mnemonicWalletsDataSrv: mnemonicWalletsDataSrv,
		encryptSrv:             encryptSrv,
		walletUnits:            make(map[uuid.UUID]WalletPoolUnitService, 0),
	}
}
