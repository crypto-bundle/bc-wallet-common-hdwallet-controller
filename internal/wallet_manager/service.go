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

	err = s.walletPoolSrv.SetWalletUnits(ctx, s.walletPoolInitializerSrv.GetWalletPoolUnits())
	if err != nil {
		return err
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

func (s *Service) GetEnabledWalletsUUID(ctx context.Context) ([]string, error) {
	uuidList, err := s.walletsDataSrv.GetAllEnabledWalletUUIDList(ctx)
	if err != nil {
		return nil, err
	}

	return uuidList, nil
}

func (s *Service) GetAddressByPath(ctx context.Context,
	walletUUID uuid.UUID,
	mnemonicWalletUUID uuid.UUID,
	account, change, index uint32,
) (string, error) {
	return s.walletPoolSrv.GetAddressByPath(ctx, walletUUID, mnemonicWalletUUID, account, change, index)
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
	poolUnit, err := s.walletPoolUnitMakerSrv.CreateWallet(ctx, strategy, title, purpose)
	if err != nil {
		return nil, err
	}

	err = s.walletPoolSrv.AddAndStartWalletUnit(ctx, poolUnit.GetWalletUUID(), poolUnit)
	if err != nil {
		return nil, err
	}

	return poolUnit.GetWalletPublicData(), nil
}

func New(logger *zap.Logger,
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

	return &Service{
		logger: logger,
		cfg:    cfg,

		txStmtManager:          txStmtManager,
		walletsDataSrv:         walletDataSrv,
		mnemonicWalletsDataSrv: mnemonicWalletDataSrv,

		mnemonicGeneratorSrv: mnemonicGeneratorSrv,
		encryptSrv:           encryptSrv,

		walletPoolInitializerSrv: walletPoolInitSrv,
	}, nil
}
