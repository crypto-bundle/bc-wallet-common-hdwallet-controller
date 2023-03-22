/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	tronCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
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
		err = s.walletPoolSrv.SetWalletUnits(ctx, s.walletPoolInitializerSrv.GetWalletPoolUnits())
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

func (s *Service) Shutdown(ctx context.Context) error {
	err := s.walletPoolSrv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
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
	accountIndex uint32,
	internalIndex uint32,
	addressIndexFrom uint32,
	addressIndexTo uint32,
) ([]*types.PublicDerivationAddressData, error) {
	return s.walletPoolSrv.GetAddressesByPathByRange(ctx, walletUUID, mnemonicWalletUUID,
		accountIndex, internalIndex, addressIndexFrom, addressIndexTo)
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
	transaction *tronCore.Transaction,
) (*types.PublicSignTxData, error) {
	return s.walletPoolSrv.SignTransaction(ctx, walletUUID, mnemonicUUID,
		account, change, index, transaction)
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
