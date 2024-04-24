package wallet_manager

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Pool struct {
	logger *zap.Logger
	cfg    configService

	runTimeCtx context.Context

	walletsDataSrv         walletsDataService
	mnemonicWalletsDataSrv mnemonicWalletsDataService
	encryptSrv             encryptService

	walletUnitsCount uint
	walletUnits      map[uuid.UUID]WalletPoolUnitService
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
	p.runTimeCtx = ctx

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
	p.walletUnitsCount = uint(len(walletUnits))

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
	p.walletUnitsCount++

	return nil
}

func (p *Pool) AddAndStartWalletUnit(_ context.Context,
	walletUUID uuid.UUID,
	walletUnit WalletPoolUnitService,
) error {
	_, isExists := p.walletUnits[walletUUID]
	if isExists {
		return ErrPassedWalletAlreadyExists
	}

	err := walletUnit.Init(p.runTimeCtx)
	if err != nil {
		return err
	}

	err = walletUnit.Run(p.runTimeCtx)
	if err != nil {
		return err
	}

	p.walletUnits[walletUUID] = walletUnit
	p.walletUnitsCount++

	return nil
}

func (p *Pool) GetAddressByPath(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicWalletUUID uuid.UUID,
	account, change, index uint32,
) (string, error) {
	poolUnit, isExists := p.walletUnits[walletUUID]
	if !isExists {
		return "", ErrPassedWalletNotFound
	}

	return poolUnit.GetAddressByPath(ctx, mnemonicWalletUUID, account, change, index)
}

func (p *Pool) GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*types.PublicWalletData, error) {
	poolUnit, isExists := p.walletUnits[walletUUID]
	if !isExists {
		return nil, nil
	}

	return poolUnit.GetWalletPublicData(), nil
}

func (p *Pool) GetEnabledWallets(ctx context.Context) ([]*types.PublicWalletData, error) {
	if p.walletUnitsCount == 0 {
		return nil, nil
	}

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
	rangeIterable types.AddrRangeIterable,
	marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
) error {
	poolUnit, isExists := p.walletUnits[walletUUID]
	if !isExists {
		return ErrPassedWalletNotFound
	}

	return poolUnit.GetAddressesByPathByRange(ctx, mnemonicWalletUUID,
		rangeIterable, marshallerCallback)
}

func (p *Pool) SignTransaction(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicUUID uuid.UUID,
	account, change, index uint32,
	transactionData []byte,
) (*types.PublicSignTxData, error) {
	poolUnit, isExists := p.walletUnits[walletUUID]
	if !isExists {
		p.logger.Error("wallet is not exists in wallet pool", zap.String(app.WalletUUIDTag, walletUUID.String()))
		return nil, ErrPassedWalletNotFound
	}

	return poolUnit.SignTransaction(ctx, mnemonicUUID, account, change, index, transactionData)
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
		walletUnitsCount:       0,
	}
}
