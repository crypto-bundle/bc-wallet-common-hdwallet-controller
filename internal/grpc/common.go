package grpc

import (
	"context"

	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/entities"
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
