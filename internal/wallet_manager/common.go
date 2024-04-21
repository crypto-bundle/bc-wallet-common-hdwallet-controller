package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"time"
)

type configService interface {
	GetDefaultWalletSessionDelay() time.Duration
	GetDefaultWalletUnloadInterval() time.Duration
}

type mnemonicWalletsCacheStoreService interface {
	SetMnemonicWalletItem(ctx context.Context,
		walletItem *entities.MnemonicWallet,
	) error
	SetMultipleMnemonicWallets(ctx context.Context,
		walletItems []*entities.MnemonicWallet,
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
	UnsetWalletSession(ctx context.Context,
		mnemonicWalletsUUID string,
		sessionsUUID string,
	) error
	UnsetMultipleWallets(ctx context.Context,
		mnemonicWalletsUUIDs []string,
		sessionsUUIDs []string,
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
	) (uint, []string, error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByUUID(ctx context.Context, uuid string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletsByStatus(ctx context.Context,
		status types.MnemonicWalletStatus,
	) ([]*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDList(ctx context.Context,
		UUIDList []string,
	) ([]*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDListAndStatus(ctx context.Context,
		UUIDList []string,
		status []types.MnemonicWalletStatus,
	) ([]string, []*entities.MnemonicWallet, error)

	AddNewWalletSession(ctx context.Context,
		sessionItem *entities.MnemonicWalletSession,
	) (*entities.MnemonicWalletSession, error)
	UpdateWalletSessionStatusByWalletUUID(ctx context.Context,
		walletUUID string,
		status types.MnemonicWalletSessionStatus,
	) error
	UpdateWalletSessionStatusBySessionUUID(ctx context.Context,
		sessionUUID string,
		newStatus types.MnemonicWalletSessionStatus,
	) (result *entities.MnemonicWalletSession, err error)
	UpdateMultipleWalletSessionStatus(ctx context.Context,
		sessionsUUIDs []string,
		newStatus types.MnemonicWalletSessionStatus,
	) (count uint, sessions []string, err error)
	GetWalletSessionByUUID(ctx context.Context,
		sessionUUID string,
	) (*entities.MnemonicWalletSession, error)
	GetActiveWalletSessionsByWalletUUID(ctx context.Context, walletUUID string) (
		count uint, list []*entities.MnemonicWalletSession, err error,
	)
}

type signRequestDataService interface {
	AddSignRequestItem(ctx context.Context,
		toSaveItem *entities.SignRequest,
	) (*entities.SignRequest, error)
	UpdateSignRequestItemStatus(ctx context.Context,
		signReqUUID string,
		newStatus types.SignRequestStatus,
	) error
}

type encryptService interface {
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(encMsg []byte) ([]byte, error)
}

type eventPublisherService interface {
	SendSessionStartEvent(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) error
	SendSessionClosedEvent(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) error
}

type transactionalStatementManager interface {
	BeginContextualTxStatement(ctx context.Context) (context.Context, error)
	CommitContextualTxStatement(ctx context.Context) error
	RollbackContextualTxStatement(ctx context.Context) error
	BeginTxWithRollbackOnError(ctx context.Context,
		callback func(txStmtCtx context.Context) error,
	) error
}
