package grpc

import (
	"bc-wallet-eth-hdwallet/internal/entities"
	"context"
)

type walleter interface {
	GetAddressByPath(ctx context.Context,
		walletUUID string,
		account, change, index uint32,
	) (string, error)

	CreateNewMnemonicWallet(ctx context.Context,
		title string,
		purpose string,
		isHot bool,
	) (*entities.MnemonicWallet, error)

	GetEnabledWalletsUUID(ctx context.Context) ([]string, error)
}
