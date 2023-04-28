package wallet_manager

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/hdwallet"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	"github.com/google/uuid"
)

type configService interface {
	GetDefaultHotWalletUnloadInterval() time.Duration
	GetDefaultWalletUnloadInterval() time.Duration

	GetMnemonicsCountPerWallet() uint8
}

type mnemonicWalletsDataService interface {
	AddNewMnemonicWallet(ctx context.Context, wallet *entities.MnemonicWallet) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletUUID(ctx context.Context, walletUUID uuid.UUID) (*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDList(ctx context.Context,
		UUIDList []string,
	) ([]*entities.MnemonicWallet, error)
	GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetAllNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
}

type walletsDataService interface {
	AddNewWallet(ctx context.Context, wallet *entities.Wallet) (*entities.Wallet, error)
	GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*entities.Wallet, error)
	GetAllEnabledWallets(ctx context.Context) ([]*entities.Wallet, error)
	GetAllEnabledWalletUUIDList(ctx context.Context) ([]string, error)
	SetEnabledToWalletByUUID(ctx context.Context, uuid string) error
}

type walletPoolUnitMakerService interface {
	CreateDisabledWallet(ctx context.Context,
		strategy types.WalletMakerStrategy,
		title, purpose string,
	) (WalletPoolUnitService, error)
}

type WalletPoolUnitService interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error

	AddMnemonicUnit(unit walletPoolMnemonicUnitService) error
	GetWalletUUID() uuid.UUID
	GetWalletTitle() string
	GetWalletPurpose() string
	GetWalletPublicData() *types.PublicWalletData
	GetAddressByPath(ctx context.Context,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
	) (string, error)
	GetAddressesByPathByRange(ctx context.Context,
		mnemonicWalletUUID uuid.UUID,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error
	SignTransaction(ctx context.Context,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
		transactionData []byte,
	) (*types.PublicSignTxData, error)
}

type walletPoolMnemonicUnitService interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
	LoadWallet(ctx context.Context) error
	UnloadWallet(ctx context.Context) error
	GetPublicData() *types.PublicMnemonicWalletData

	IsHotWalletUnit() bool
	GetMnemonicUUID() uuid.UUID
	GetAddressByPath(ctx context.Context,
		account, change, index uint32,
	) (string, error)
	GetAddressesByPathByRange(ctx context.Context,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error
	SignTransaction(ctx context.Context,
		account, change, index uint32,
		transactionData []byte,
	) (*types.PublicSignTxData, error)
}

type walletPoolInitService interface {
	LoadAndInitWallets(ctx context.Context) error
	GetWalletPoolUnits() map[uuid.UUID]WalletPoolUnitService
}

type walletPoolService interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error

	SetWalletUnits(ctx context.Context,
		walletUnits map[uuid.UUID]WalletPoolUnitService,
	) error
	AddAndStartWalletUnit(ctx context.Context,
		walletUUID uuid.UUID,
		walletUnit WalletPoolUnitService,
	) error
	GetAddressByPath(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		account, change, index uint32,
	) (string, error)
	GetAddressesByPathByRange(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicWalletUUID uuid.UUID,
		rangeIterable types.AddrRangeIterable,
		marshallerCallback func(accountIndex, internalIndex, addressIdx, position uint32, address string),
	) error
	GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*types.PublicWalletData, error)
	GetEnabledWallets(ctx context.Context) ([]*types.PublicWalletData, error)
	SignTransaction(ctx context.Context,
		walletUUID uuid.UUID,
		mnemonicUUID uuid.UUID,
		account, change, index uint32,
		transactionData []byte,
	) (*types.PublicSignTxData, error)
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

	NewTronWallet(account, change, address uint32) (*hdwallet.Tron, error)

	ClearSecrets()
}

type mnemonicGenerator interface {
	Generate(ctx context.Context) (string, error)
}

type encryptService interface {
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(encMsg []byte) ([]byte, error)
}
