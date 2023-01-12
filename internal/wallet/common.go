package wallet

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/entities"
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/hdwallet"
)

type config interface {
}

type repo interface {
	AddNewMnemonicWallet(ctx context.Context, wallet *entities.MnemonicWallet) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletUUID(ctx context.Context, uuid string) (*entities.MnemonicWallet, error)
	GetAllEnabledMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetAllEnabledHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetAllEnabledNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
}

type mnemonicWalletConfig interface {
	GetMnemonicWalletPurpose() string
	GetMnemonicWalletHash() string
	IsHotWallet() bool
}

type walleter interface {
	GetAddress() (string, error)
	GetPubKey() string
	GetPrvKey() (string, error)
	GetPath() string
}

type hdWalleter interface {
	PublicHex() string
	PublicHash() ([]byte, error)

	NewEthWallet(account, change, address uint32) (*hdwallet.ETH, error)
}

type mnemonicGenerator interface {
	Generate(ctx context.Context) (string, error)
}

type crypto interface {
	Encrypt(msg string) (string, error)
	Decrypt(encMsg string) ([]byte, error)
}
