package wallet_manager

import (
	"context"
)

func (s *Service) CheckSession(ctx context.Context,
	mnemonicUUID string,
	sessionUUID string,
) (isActive bool, err error) {
	wallet, walletSession, err := s.getWalletAndSession(ctx, mnemonicUUID, sessionUUID)
	if err != nil {
		return false, err
	}

	if wallet == nil {
		return false, nil
	}

	if walletSession == nil {
		return false, nil
	}

	return walletSession.IsSessionActive(), nil
}
