package wallet_manager

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) ImportWallet(ctx context.Context, importedData []byte) (*entities.MnemonicWallet, error) {
	decryptedData, err := s.transitEncryptorSvc.Decrypt(importedData)
	if err != nil {
		return nil, err
	}

	encryptedMnemonicData, err := s.appEncryptorSvc.Encrypt(decryptedData)
	if err != nil {
		return nil, err
	}

	mnemonicHash := fmt.Sprintf("%x", sha256.Sum256(decryptedData))
	vaultEncryptedHash := fmt.Sprintf("%x", sha256.Sum256(encryptedMnemonicData))

	toSaveItem := &entities.MnemonicWallet{
		UUID:               uuid.New(),
		MnemonicHash:       mnemonicHash,
		Status:             types.MnemonicWalletStatusCreated,
		UnloadInterval:     s.cfg.GetDefaultWalletUnloadInterval(),
		VaultEncrypted:     encryptedMnemonicData,
		VaultEncryptedHash: vaultEncryptedHash,
		CreatedAt:          time.Now(),
		UpdatedAt:          nil,
	}

	resp, err := s.hdWalletClientSvc.LoadMnemonic(ctx, &hdwallet.LoadMnemonicRequest{
		MnemonicIdentity: &common.MnemonicWalletIdentity{
			WalletUUID: toSaveItem.UUID.String(),
		},
		EncryptedMnemonicData: encryptedMnemonicData,
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		s.logger.Error("missing resp in load mnemonic request", zap.Error(ErrMissingHdWalletResp),
			zap.String(app.MnemonicWalletUUIDTag, toSaveItem.UUID.String()))

		return nil, ErrMissingHdWalletResp
	}

	unLoadResp, err := s.hdWalletClientSvc.UnLoadMnemonic(ctx, &hdwallet.UnLoadMnemonicRequest{
		MnemonicIdentity: &common.MnemonicWalletIdentity{
			WalletUUID: toSaveItem.UUID.String(),
		},
	})
	if err != nil {
		return nil, err
	}

	if unLoadResp == nil {
		s.logger.Error("missing resp in unload mnemonic request", zap.Error(ErrMissingHdWalletResp),
			zap.String(app.MnemonicWalletUUIDTag, toSaveItem.UUID.String()))

		return nil, ErrMissingHdWalletResp
	}

	return s.saveWallet(ctx, toSaveItem, resp.MnemonicIdentity, encryptedMnemonicData)
}
