package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
)

func (s *Service) GetWalletByUUID(ctx context.Context, walletUUID string) (*entities.MnemonicWallet, error) {
	item, err := s.cacheStoreDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
	if err != nil {
		return nil, err
	}

	if item != nil {
		return item, nil
	}

	item, err = s.mnemonicWalletsDataSvc.GetMnemonicWalletByUUID(ctx, walletUUID)
	if err != nil {
		return nil, err
	}

	return item, nil
}
