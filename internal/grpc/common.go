package grpc

import (
	"bc-wallet-eth-hdwallet/internal/entities"
	"context"
	"math"
)

const (
	MethodNameTag = "method_name"

	DefaultServerMaxReceiveMessageSize = math.MaxInt32
	DefaultServerMaxSendMessageSize    = 1024 * 1024 * 24
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
}
