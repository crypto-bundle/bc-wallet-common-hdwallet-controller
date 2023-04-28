package wallet_manager

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletUnitMaker struct {
	logger *zap.Logger
	cfg    configService

	walletsDataSrv         walletsDataService
	mnemonicWalletsDataSrv mnemonicWalletsDataService
	txStmtManager          transactionalStatementManager
	mnemonicGeneratorSrv   mnemonicGenerator
	encryptSrv             encryptService
}

func (m *WalletUnitMaker) CreateDisabledWallet(ctx context.Context,
	strategy types.WalletMakerStrategy,
	title, purpose string,
) (WalletPoolUnitService, error) {
	switch strategy {
	case types.WalletMakerSingleMnemonicStrategy:
		return m.createSingleMnemonicWallet(ctx, title, purpose)
	case types.WalletMakerMultipleMnemonicStrategy:
		return m.createMultipleMnemonicWallet(ctx, title, purpose)
	default:
		return m.createMultipleMnemonicWallet(ctx, title, purpose)
	}
}

func (m *WalletUnitMaker) createSingleMnemonicWallet(ctx context.Context,
	title, purpose string,
) (*singleMnemonicWalletUnit, error) {
	walletEntity := &entities.Wallet{
		Title:     title,
		UUID:      uuid.New(),
		Purpose:   purpose,
		Strategy:  types.WalletMakerSingleMnemonicStrategy,
		IsEnabled: false, // temporary until wallet not successfully init
		CreatedAt: time.Now(),
		UpdatedAt: nil, // temporary until wallet not successfully init
	}

	var mnemonicItem *entities.MnemonicWallet = nil

	walletPoolUnit := newSingleMnemonicWalletPoolUnit(m.logger, walletEntity.UUID,
		walletEntity.Title, walletEntity.Purpose)

	err := m.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		_, txStmtErr := m.walletsDataSrv.AddNewWallet(ctx, walletEntity)
		if txStmtErr != nil {
			return txStmtErr
		}

		mnemonicItem, txStmtErr = m.createNewMnemonicWallet(ctx, walletEntity.UUID, true)
		if txStmtErr != nil {
			return txStmtErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	mnemonicUnit := newMnemonicWalletPoolUnit(m.logger, m.cfg,
		m.cfg.GetDefaultHotWalletUnloadInterval(), walletEntity.UUID, m.encryptSrv,
		m.mnemonicWalletsDataSrv, mnemonicItem)

	err = walletPoolUnit.AddMnemonicUnit(mnemonicUnit)
	if err != nil {
		return nil, err
	}

	return walletPoolUnit, nil
}

func (m *WalletUnitMaker) createMultipleMnemonicWallet(ctx context.Context,
	title, purpose string,
) (*multipleMnemonicWalletUnit, error) {
	walletEntity := &entities.Wallet{
		Title:     title,
		UUID:      uuid.New(),
		Purpose:   purpose,
		Strategy:  types.WalletMakerMultipleMnemonicStrategy,
		IsEnabled: false, // temporary until wallet not successfully init
		CreatedAt: time.Now(),
		UpdatedAt: nil, // temporary until wallet not successfully init
	}
	var hotMnemonicItem *entities.MnemonicWallet = nil
	mnemonicItems := make([]*entities.MnemonicWallet, m.cfg.GetMnemonicsCountPerWallet())

	walletPoolUnit := newMultipleMnemonicWalletPoolUnit(m.logger, m.cfg, m.encryptSrv, m.walletsDataSrv,
		m.mnemonicWalletsDataSrv,
		walletEntity.UUID,
		walletEntity.Title, walletEntity.Purpose)

	err := m.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		_, txStmtErr := m.walletsDataSrv.AddNewWallet(txStmtCtx, walletEntity)
		if txStmtErr != nil {
			return txStmtErr
		}

		hotMnemonicItem, txStmtErr = m.createNewMnemonicWallet(txStmtCtx, walletEntity.UUID, true)
		if txStmtErr != nil {
			return txStmtErr
		}

		mnemonicItems[0] = hotMnemonicItem

		for i, j := uint8(0), 1; i != m.cfg.GetMnemonicsCountPerWallet()-1; i++ {
			nonHotMnemonicItem, createItem := m.createNewMnemonicWallet(txStmtCtx, walletEntity.UUID, false)
			if createItem != nil {
				return createItem
			}

			mnemonicItems[j] = nonHotMnemonicItem
			j++
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	for i := uint8(0); i != m.cfg.GetMnemonicsCountPerWallet(); i++ {
		mnemonicUnit := newMnemonicWalletPoolUnit(m.logger, m.cfg,
			m.cfg.GetDefaultWalletUnloadInterval(), walletEntity.UUID, m.encryptSrv,
			m.mnemonicWalletsDataSrv, mnemonicItems[i])

		addErr := walletPoolUnit.AddMnemonicUnit(mnemonicUnit)
		if addErr != nil {
			return nil, addErr
		}
	}

	return walletPoolUnit, nil
}

func (m *WalletUnitMaker) createNewMnemonicWallet(ctx context.Context,
	walletUUID uuid.UUID,
	isHotWallet bool,
) (*entities.MnemonicWallet, error) {
	newWalletMnemonic, err := m.mnemonicGeneratorSrv.Generate(ctx)
	if err != nil {
		return nil, err
	}

	encMnemonic, err := m.encryptSrv.Encrypt([]byte(newWalletMnemonic))
	if err != nil {
		return nil, err
	}

	mnemonicWalletEntity := &entities.MnemonicWallet{
		UUID:               uuid.New(),
		WalletUUID:         walletUUID,
		MnemonicHash:       fmt.Sprintf("%x", sha256.Sum256([]byte(newWalletMnemonic))),
		IsHotWallet:        isHotWallet,
		RsaEncrypted:       []byte("nil"), // temporary nil string TODO: add rsa encryption via merchant service
		RsaEncryptedHash:   "nil",         // temporary nil string TODO: add rsa encryption via merchant service
		VaultEncrypted:     encMnemonic,
		VaultEncryptedHash: fmt.Sprintf("%x", sha256.Sum256([]byte(encMnemonic))),
		UnloadInterval:     m.cfg.GetDefaultWalletUnloadInterval(),
	}

	mnemonicWalletEntity, err = m.mnemonicWalletsDataSrv.AddNewMnemonicWallet(ctx, mnemonicWalletEntity)
	if err != nil {
		return nil, err
	}

	return mnemonicWalletEntity, err
}

func newWalletMaker(logger *zap.Logger,
	cfg configService,
	walletsDataSrv walletsDataService,
	mnemonicWalletsDataSrv mnemonicWalletsDataService,
	txStmtManager transactionalStatementManager,
	mnemonicGeneratorSrv mnemonicGenerator,
	encryptSrv encryptService,
) *WalletUnitMaker {
	return &WalletUnitMaker{
		logger:                 logger,
		cfg:                    cfg,
		walletsDataSrv:         walletsDataSrv,
		mnemonicWalletsDataSrv: mnemonicWalletsDataSrv,
		txStmtManager:          txStmtManager,
		mnemonicGeneratorSrv:   mnemonicGeneratorSrv,
		encryptSrv:             encryptSrv,
	}
}
