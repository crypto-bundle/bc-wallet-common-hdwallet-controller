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

package wallet

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/hdwallet"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/mnemonic"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/wallet/repository"

	"github.com/crypto-bundle/bc-wallet-common/pkg/postgres"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger         *zap.Logger
	cfg            config
	repo           repo
	mnemoGenerator mnemonicGenerator
	cryptoSrv      crypto

	// hotWallets - Ethereum hot wallets
	hotWallets map[string]walleter
	// hotHdWallets - hot hdwallets
	hotHdWallets map[string]hdWalleter
}

func (s *Service) Init(ctx context.Context) error {
	wallets, err := s.repo.GetAllEnabledMnemonicWallets(ctx)
	if err != nil {
		return err
	}

	for i, _ := range wallets {
		mnemonicBytes, err := s.cryptoSrv.Decrypt(wallets[i].RsaEncrypted)
		if err != nil {
			return err
		}

		mnemonicSum256 := sha256.Sum256(mnemonicBytes)
		if hex.EncodeToString(mnemonicSum256[:]) != wallets[i].Hash {
			return ErrWrongMnemonicHash
		}

		hdWallet, creatErr := hdwallet.NewFromString(string(mnemonicBytes))
		if creatErr != nil {
			return creatErr
		}

		ethWallet, walletErr := hdWallet.NewEthWallet(uint32(0), uint32(0), uint32(0))
		if walletErr != nil {
			return walletErr
		}

		s.hotWallets[wallets[i].UUID.String()] = ethWallet
		s.hotHdWallets[wallets[i].UUID.String()] = hdWallet
	}

	return nil
}

func (s *Service) Shutdown(_ context.Context) error {
	return nil
}

func (s *Service) GetMnemonicWallet(ctx context.Context,
	walletConfig mnemonicWalletConfig,
) (*entities.MnemonicWallet, error) {
	wallet, err := s.repo.GetMnemonicWalletByHash(ctx, walletConfig.GetMnemonicWalletHash())
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *Service) GetEnabledWalletsUUID(ctx context.Context) ([]string, error) {
	wallets, err := s.repo.GetAllEnabledNonHotMnemonicWallets(ctx)
	if err != nil {
		return nil, err
	}

	uuids := make([]string, len(wallets))

	for i, _ := range wallets {
		uuids[i] = wallets[i].UUID.String()
	}

	return uuids, nil
}

func (s *Service) GetAddressByPath(_ context.Context,
	walletUUID string,
	account, change, index uint32,
) (string, error) {
	pickedHDWallet, ok := s.hotHdWallets[walletUUID]
	if !ok {
		return "", ErrPassedWalletNotFound
	}

	ethWallet, err := pickedHDWallet.NewEthWallet(account, change, index)
	if err != nil {
		return "", err
	}

	ethAddress, err := ethWallet.GetAddress()
	if err != nil {
		return "", err
	}

	return ethAddress, nil
}

func (s *Service) CreateNewMnemonicWallet(ctx context.Context,
	title string,
	purpose string,
	isHot bool,
) (*entities.MnemonicWallet, error) {
	newWalletMnemonic, err := s.mnemoGenerator.Generate(ctx)
	if err != nil {
		return nil, err
	}

	encMnemonic, err := s.cryptoSrv.Encrypt(newWalletMnemonic)
	if err != nil {
		return nil, err
	}

	walletEntity := &entities.MnemonicWallet{
		Title:            title,
		UUID:             uuid.New(),
		Hash:             fmt.Sprintf("%x", sha256.Sum256([]byte(newWalletMnemonic))),
		IsHotWallet:      isHot,
		Purpose:          purpose,
		RsaEncrypted:     encMnemonic,
		RsaEncryptedHash: fmt.Sprintf("%x", sha256.Sum256([]byte(encMnemonic))),
	}

	walletEntity, err = s.repo.AddNewMnemonicWallet(ctx, walletEntity)
	if err != nil {
		return nil, err
	}

	if isHot {
		hdWallet, creatErr := hdwallet.NewFromString(newWalletMnemonic)
		if creatErr != nil {
			return nil, creatErr
		}

		ethWallet, walletErr := hdWallet.NewEthWallet(uint32(0), uint32(0), uint32(0))
		if walletErr != nil {
			return nil, walletErr
		}

		uuidStr := walletEntity.UUID.String()

		s.hotWallets[uuidStr] = ethWallet
		s.hotHdWallets[uuidStr] = hdWallet
	}

	return walletEntity, nil
}

func New(logger *zap.Logger,
	cfg config,
	pgConn *postgres.Connection,
	cryptoSrv crypto,
) (*Service, error) {
	pgRepos, err := repository.NewPostgresStore(logger, pgConn)
	if err != nil {
		return nil, err
	}

	mnemoGenerator, err := mnemonic.NewMnemonicGenerator(logger)
	if err != nil {
		return nil, err
	}

	return &Service{
		logger:         logger,
		cfg:            cfg,
		repo:           pgRepos,
		mnemoGenerator: mnemoGenerator,
		cryptoSrv:      cryptoSrv,

		hotWallets:   make(map[string]walleter, 0),
		hotHdWallets: make(map[string]hdWalleter, 0),
	}, nil
}
