package wallet_manager

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) AddNewWallet(ctx context.Context) (*entities.MnemonicWallet, error) {
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

	resp, err := s.hdWalletClientSvc.GenerateMnemonic(ctx, &hdwallet.GenerateMnemonicRequest{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: toSaveItem.UUID.String(),
		},
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

func (s *Service) saveWallet(ctx context.Context,
	walletItem *entities.MnemonicWallet,
	hdWalletInfo *pbCommon.MnemonicWalletIdentity,
	encryptedData []byte,
) (wallet *entities.MnemonicWallet, err error) {
	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		walletItem.MnemonicHash = hdWalletInfo.WalletHash
		walletItem.VaultEncryptedHash = fmt.Sprintf("%x", sha256.Sum256(encryptedData))
		walletItem.VaultEncrypted = encryptedData

		savedItem, clbErr := s.mnemonicWalletsDataSvc.AddNewMnemonicWallet(txStmtCtx,
			walletItem)
		if clbErr != nil {
			s.logger.Error("unable to save mnemonic wallet item in persistent store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, walletItem.UUID.String()))

			return clbErr
		}

		wallet = savedItem

		return nil
	})
	if err != nil {
		s.logger.Error("unable to save new wallet", zap.Error(err))

		return nil, err
	}

	return
}
