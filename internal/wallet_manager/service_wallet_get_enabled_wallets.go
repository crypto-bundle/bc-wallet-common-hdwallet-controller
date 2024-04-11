package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
)

func (s *Service) GetEnabledWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	items, err := s.cacheStoreDataSvc.GetAllWallets(ctx)
	if err != nil {
		return nil, err
	}

	if items != nil {
		return items, nil
	}

	items, err = s.mnemonicWalletsDataSvc.GetMnemonicWalletsByStatus(ctx, types.MnemonicWalletStatusEnabled)
	if err != nil {
		return nil, err
	}

	return items, nil
}
