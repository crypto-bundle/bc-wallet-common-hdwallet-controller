package wallet_manager

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/app"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/hdwallet"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrMissingHdWalletResp = errors.New("missing hd-wallet api response")
)

type Service struct {
	logger *zap.Logger
	cfg    configService

	mnemonicWalletsDataSrv mnemonicWalletsDataService
	cacheStoreDataSvc      mnemonicWalletsCacheStoreService

	hdwalletClientSvc hdwallet.HdWalletApiClient

	txStmtManager transactionalStatementManager
}

func (s *Service) AddNewWallet(ctx context.Context) (*entities.MnemonicWallet, error) {
	var resultItem *entities.MnemonicWallet = nil

	err := s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		toSaveItem := &entities.MnemonicWallet{
			UUID:               uuid.New(),
			MnemonicHash:       "",
			Status:             types.MnemonicWalletStatusCreated,
			UnloadInterval:     s.cfg.GetDefaultWalletUnloadInterval(),
			VaultEncrypted:     nil,
			VaultEncryptedHash: "",
			CreatedAt:          time.Time{},
			UpdatedAt:          nil,
		}

		resp, clbErr := s.hdwalletClientSvc.GenerateMnemonic(txStmtCtx, &hdwallet.GenerateMnemonicRequest{
			MnemonicIdentity: &common.MnemonicWalletIdentity{
				WalletUUID: toSaveItem.UUID.String(),
			},
		})

		if resp == nil {
			s.logger.Error("missing resp in generate mnemonic request", zap.Error(ErrMissingHdWalletResp),
				zap.String(app.MnemonicWalletUUIDTag, toSaveItem.UUID.String()))

			return ErrMissingHdWalletResp
		}

		toSaveItem.MnemonicHash = resp.MnemonicIdentity.WalletHash
		toSaveItem.VaultEncryptedHash = fmt.Sprintf("%x", sha256.Sum256(resp.EncryptedMnemonicData))
		toSaveItem.VaultEncrypted = resp.EncryptedMnemonicData

		savedItem, clbErr := s.mnemonicWalletsDataSrv.AddNewMnemonicWallet(txStmtCtx,
			toSaveItem)
		if clbErr != nil {
			s.logger.Error("unable to save mnemonic wallet item in persistent store", zap.Error(clbErr),
				zap.String(app.MnemonicWalletUUIDTag, toSaveItem.UUID.String()))

			return clbErr
		}

		resultItem = savedItem

		return nil
	})
	if err != nil {
		s.logger.Error("unable to save new wallet", zap.Error(err))

		return nil, err
	}

	return resultItem, nil
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

func NewService(logger *zap.Logger,
	cfg configService,
	walletDataSrv walletsDataService,
	mnemonicWalletDataSrv mnemonicWalletsDataService,
	txStmtManager transactionalStatementManager,
) (*Service, error) {
	return &Service{
		logger: logger,
		cfg:    cfg,

		txStmtManager:          txStmtManager,
		walletsDataSrv:         walletDataSrv,
		mnemonicWalletsDataSrv: mnemonicWalletDataSrv,
	}, nil
}
