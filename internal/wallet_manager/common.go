package wallet_manager

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
)

type configService interface {
	GetDefaultHotWalletUnloadInterval() time.Duration
	GetDefaultWalletUnloadInterval() time.Duration

	GetMnemonicsCountPerWallet() uint8
}

type mnemonicWalletsCacheStoreService interface {
	SetMnemonicWalletItem(ctx context.Context,
		walletItem *entities.MnemonicWallet,
	) error
	GetAllWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetMnemonicWalletByUUID(ctx context.Context,
		MnemonicWalletUUID string,
	) (*entities.MnemonicWallet, error)
	GetMnemonicWalletInfoByUUID(ctx context.Context,
		MnemonicWalletUUID string,
	) (wallet *entities.MnemonicWallet, sessions []*entities.MnemonicWalletSession, err error)
	GetMnemonicWalletSessionInfoByUUID(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) (wallet *entities.MnemonicWallet, session *entities.MnemonicWalletSession, err error)
	SetMnemonicWalletSessionItem(ctx context.Context,
		sessionItem *entities.MnemonicWalletSession,
	) error
	FullUnsetMnemonicWallet(ctx context.Context,
		mnemonicWalletUUID string,
	) error
}

type mnemonicWalletsDataService interface {
	AddNewMnemonicWallet(ctx context.Context,
		wallet *entities.MnemonicWallet,
	) (*entities.MnemonicWallet, error)
	UpdateWalletStatus(ctx context.Context,
		walletUUID string,
		newStatus types.MnemonicWalletStatus,
	) (*entities.MnemonicWallet, error)
	UpdateMultipleWalletsStatus(ctx context.Context,
		walletUUID []string,
		newStatus types.MnemonicWalletStatus,
	) (uint, []*entities.MnemonicWallet, error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByUUID(ctx context.Context, uuid string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletsByStatus(ctx context.Context,
		status types.MnemonicWalletStatus,
	) ([]*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDList(ctx context.Context,
		UUIDList []string,
	) ([]*entities.MnemonicWallet, error)

	UpdateWalletSessionStatusByWalletUUID(ctx context.Context,
		walletUUID string,
		status types.MnemonicWalletSessionStatus,
	) error
	GetWalletSessionByUUID(ctx context.Context,
		sessionUUID string,
	) (*entities.MnemonicWalletSession, error)
	GetActiveWalletSessionsByWalletUUID(ctx context.Context, walletUUID string) (
		count uint, list []*entities.MnemonicWalletSession, err error,
	)
}

type mnemonicWalletConfig interface {
	GetMnemonicWalletPurpose() string
	GetMnemonicWalletHash() string
	IsHotWallet() bool
}

type transactionalStatementManager interface {
	BeginContextualTxStatement(ctx context.Context) (context.Context, error)
	CommitContextualTxStatement(ctx context.Context) error
	RollbackContextualTxStatement(ctx context.Context) error
	BeginTxWithRollbackOnError(ctx context.Context,
		callback func(txStmtCtx context.Context) error,
	) error
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

	NewTronWallet(account, change, address uint32) ([]byte, error)

	ClearSecrets()
}
