package wallet_manager

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"sync"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/hdwallet"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	tronCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type MnemonicWalletUnit struct {
	logger *zap.Logger

	mu          sync.Mutex
	onAirTicker *time.Ticker

	cfgSrv                 configService
	hdWalletSrv            hdWalleter
	cryptoSrv              encryptService
	mnemonicWalletsDataSrv mnemonicWalletsDataService

	isWalletLoaded      bool
	mnemonicWalletUUID  uuid.UUID
	unloadTimerInterval time.Duration
	walletEntity        *entities.MnemonicWallet
	// privateKeyPool is pool of derivation addresses private keys
	// map key - string with derivation path
	// map value - ecdsa.PrivateKey
	privateKeyPool map[string]*ecdsa.PrivateKey
}

func (u *MnemonicWalletUnit) Init(ctx context.Context) error {
	return nil
}

func (u *MnemonicWalletUnit) Run(ctx context.Context) error {
	err := u.loadWallet(ctx)
	if err != nil {
		return err
	}

	u.onAirTicker = time.NewTicker(u.unloadTimerInterval)
	go u.run(ctx)

	return nil
}

func (u *MnemonicWalletUnit) run(ctx context.Context) {
	for {
		select {
		case tick, _ := <-u.onAirTicker.C:
			err := u.onUnloadTimerTick(context.Background())
			if err != nil {
				u.logger.Error("unable to unload logger by ticker", zap.Error(err),
					zap.Time(app.TickerEventTriggerTimeTag, tick))
			}

		case <-ctx.Done():
			u.onAirTicker.Stop()

			err := u.Shutdown(ctx)
			if err != nil {
				u.logger.Error("unable to shutdown by ctx cancel", zap.Error(err))
			}
		}
	}
}

func (u *MnemonicWalletUnit) onUnloadTimerTick(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.isWalletLoaded {
		return nil
	}

	err := u.unloadWallet(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (u *MnemonicWalletUnit) GetMnemonicUUID() uuid.UUID {
	return u.mnemonicWalletUUID
}

func (u *MnemonicWalletUnit) IsHotWalletUnit() bool {
	return u.walletEntity.IsHotWallet
}

func (u *MnemonicWalletUnit) GetPublicData() *types.PublicMnemonicWalletData {
	return &types.PublicMnemonicWalletData{
		UUID:        u.mnemonicWalletUUID,
		IsHotWallet: u.walletEntity.IsHotWallet,
	}
}

func (u *MnemonicWalletUnit) SignTransaction(ctx context.Context,
	account, change, index uint32,
	transaction *tronCore.Transaction,
) (*tronCore.Transaction, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.isWalletLoaded {
		u.onAirTicker.Reset(u.unloadTimerInterval)
		return u.signTransaction(ctx, account, change, index, transaction)
	}

	err := u.loadWallet(ctx)
	if err != nil {
		return nil, err
	}

	return u.signTransaction(ctx, account, change, index, transaction)

}

func (u *MnemonicWalletUnit) signTransaction(ctx context.Context,
	account, change, index uint32,
	transaction *tronCore.Transaction,
) (*tronCore.Transaction, error) {
	key := fmt.Sprintf("%d'/%d/%d", account, change, index)
	privateKey, isExists := u.privateKeyPool[key]
	if !isExists {
		tronWallet, walletErr := u.hdWalletSrv.NewTronWallet(account, change, index)
		if walletErr != nil {
			return nil, walletErr
		}

		wif, walletErr := tronWallet.GetAccountWIF()
		if walletErr != nil {
			return nil, walletErr
		}

		u.privateKeyPool[key] = wif.PrivKey.ToECDSA()
		privateKey = wif.PrivKey.ToECDSA()
	}

	rawData, err := proto.Marshal(transaction.GetRawData())
	if err != nil {
		return nil, err
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	contractList := transaction.GetRawData().GetContract()

	for range contractList {
		signature, signErr := crypto.Sign(hash, privateKey)
		if signErr != nil {
			return nil, signErr
		}

		transaction.Signature = append(transaction.Signature, signature)
	}

	return transaction, ErrMethodUnimplemented
}

func (u *MnemonicWalletUnit) GetAddressByPath(ctx context.Context,
	account, change, index uint32,
) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.isWalletLoaded {
		u.onAirTicker.Reset(u.unloadTimerInterval)
		return u.getAddressByPath(ctx, account, change, index)
	}

	err := u.loadWallet(ctx)
	if err != nil {
		return "", err
	}

	return u.getAddressByPath(ctx, account, change, index)
}

func (u *MnemonicWalletUnit) GetAddressesByPathByRange(ctx context.Context,
	accountIndex uint32,
	internalIndex uint32,
	addressIndexFrom uint32,
	addressIndexTo uint32,
) ([]*types.PublicDerivationAddressData, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.isWalletLoaded {
		u.onAirTicker.Reset(u.unloadTimerInterval)
		return u.getAddressesByPathByRange(ctx, accountIndex, internalIndex, addressIndexFrom, addressIndexTo)
	}

	err := u.loadWallet(ctx)
	if err != nil {
		return nil, err
	}

	return u.getAddressesByPathByRange(ctx, accountIndex, internalIndex, addressIndexFrom, addressIndexTo)
}

func (u *MnemonicWalletUnit) getAddressesByPathByRange(ctx context.Context,
	accountIndex uint32,
	internalIndex uint32,
	addressIndexFrom uint32,
	addressIndexTo uint32,
) ([]*types.PublicDerivationAddressData, error) {
	var err error
	rangeSize := addressIndexTo - addressIndexFrom
	result := make([]*types.PublicDerivationAddressData, rangeSize+1)
	wg := sync.WaitGroup{}
	wg.Add(int(rangeSize) + 1)

	for i, j := addressIndexFrom, uint32(0); i <= addressIndexTo; i++ {
		go func(i, j uint32) {
			defer wg.Done()

			address, getAddrErr := u.getAddressByPath(ctx, accountIndex,
				internalIndex, i)
			if getAddrErr != nil {
				u.logger.Error("unable to get address by path", zap.Error(getAddrErr),
					zap.Uint32(app.HDWalletAccountIndexTag, accountIndex),
					zap.Uint32(app.HDWalletInternalIndexTag, internalIndex),
					zap.Uint32(app.HDWalletAddressIndexTag, i))

				err = getAddrErr
				return
			}

			result[j] = &types.PublicDerivationAddressData{
				AccountIndex:  accountIndex,
				InternalIndex: internalIndex,
				AddressIndex:  i,
				Address:       address,
			}
		}(i, j)

		j++
	}

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *MnemonicWalletUnit) getAddressByPath(_ context.Context,
	account, change, index uint32,
) (string, error) {
	tronWallet, err := u.hdWalletSrv.NewTronWallet(account, change, index)
	if err != nil {
		return "", err
	}

	blockchainAddress, err := tronWallet.GetAddress()
	if err != nil {
		return "", err
	}

	return blockchainAddress, nil
}

func (u *MnemonicWalletUnit) LoadWallet(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.isWalletLoaded {
		u.onAirTicker.Reset(u.unloadTimerInterval)
	}

	return u.loadWallet(ctx)
}

func (u *MnemonicWalletUnit) loadWallet(ctx context.Context) error {
	walletEntity, err := u.mnemonicWalletsDataSrv.GetMnemonicWalletUUID(ctx, u.mnemonicWalletUUID.String())
	if err != nil {
		return err
	}
	u.walletEntity = walletEntity

	mnemonicBytes, err := u.cryptoSrv.Decrypt(u.walletEntity.RsaEncrypted)
	if err != nil {
		return err
	}

	mnemonicSum256 := sha256.Sum256(mnemonicBytes)
	if hex.EncodeToString(mnemonicSum256[:]) != u.walletEntity.MnemonicHash {
		return ErrWrongMnemonicHash
	}

	hdWallet, creatErr := hdwallet.NewFromString(string(mnemonicBytes))
	if creatErr != nil {
		return creatErr
	}
	u.hdWalletSrv = hdWallet

	return nil
}

func (u *MnemonicWalletUnit) UnloadWallet(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if !u.isWalletLoaded {
		return nil
	}

	return u.unloadWallet(ctx)
}

func (u *MnemonicWalletUnit) unloadWallet(ctx context.Context) error {
	u.hdWalletSrv = nil
	u.walletEntity = nil

	for key := range u.privateKeyPool {
		delete(u.privateKeyPool, key)
	}

	u.isWalletLoaded = false

	return nil
}

func (u *MnemonicWalletUnit) Shutdown(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.unloadWallet(ctx)
	if err != nil {
		u.logger.Error("unable to unload wallet", zap.Error(err))

		return err
	}

	return nil
}

func newMnemonicWalletPoolUnit(logger *zap.Logger,
	cfg configService,
	unloadInterval time.Duration,
	walletUUID uuid.UUID,
	cryptoSrv encryptService,
	mnemonicWalletDataSrv mnemonicWalletsDataService,
	walletItem *entities.MnemonicWallet,
) *MnemonicWalletUnit {
	return &MnemonicWalletUnit{
		logger: logger.With(zap.String(app.WalletUUIDTag, walletUUID.String())),
		mu:     sync.Mutex{},

		onAirTicker: nil, // that field will be field @ run stage

		hdWalletSrv: nil, // that field will be field @ load wallet stage

		cryptoSrv:              cryptoSrv,
		mnemonicWalletsDataSrv: mnemonicWalletDataSrv,

		isWalletLoaded:      false,
		mnemonicWalletUUID:  walletUUID,
		unloadTimerInterval: unloadInterval,
		walletEntity:        walletItem,
		privateKeyPool:      make(map[string]*ecdsa.PrivateKey, 0),
	}
}
