package wallet

import (
	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/entities"
	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/hdwallet"
	"context"
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
