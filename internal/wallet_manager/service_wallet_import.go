package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/app"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/hdwallet"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) ImportWallet(ctx context.Context, mnemonicData []byte) (*entities.MnemonicWallet, error) {
	toSaveItem := &entities.MnemonicWallet{
		UUID:               uuid.New(),
		MnemonicHash:       "",
		Status:             types.MnemonicWalletStatusCreated,
		UnloadInterval:     s.cfg.GetDefaultWalletUnloadInterval(),
		VaultEncrypted:     nil,
		VaultEncryptedHash: "",
		CreatedAt:          time.Now(),
		UpdatedAt:          nil,
	}

	resp, err := s.hdwalletClientSvc.EncryptMnemonic(ctx, &hdwallet.EncryptMnemonicRequest{
		MnemonicIdentity: &common.MnemonicWalletIdentity{
			WalletUUID: toSaveItem.UUID.String(),
		},
		MnemonicData: mnemonicData,
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		s.logger.Error("missing resp in generate mnemonic request", zap.Error(ErrMissingHdWalletResp),
			zap.String(app.MnemonicWalletUUIDTag, toSaveItem.UUID.String()))

		return nil, ErrMissingHdWalletResp
	}

	return s.saveWallet(ctx, toSaveItem, resp.MnemonicIdentity, resp.EncryptedMnemonicData)
}
