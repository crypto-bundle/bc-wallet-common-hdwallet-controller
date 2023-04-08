package wallet_manager

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
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

type addressData struct {
	address    string
	privateKey *ecdsa.PrivateKey
}

type MnemonicWalletUnit struct {
	logger *zap.Logger

	mu          sync.Mutex
	onAirTicker *time.Ticker

	cfgSrv                 configService
	hdWalletSrv            hdWalleter
	cryptoSrv              encryptService
	mnemonicWalletsDataSrv mnemonicWalletsDataService

	isWalletLoaded bool
	isHotWallet    bool

	walletUUID          uuid.UUID
	mnemonicWalletUUID  uuid.UUID
	mnemonicWalletHash  string
	unloadTimerInterval time.Duration
	walletEntity        *entities.MnemonicWallet
	// addressPool is pool of derivation addresses with private keys and address
	// map key - string with derivation path
	// map value - ecdsa.PrivateKey and address string
	addressPool map[string]*addressData
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
				continue
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
	return u.isHotWallet
}

func (u *MnemonicWalletUnit) GetPublicData() *types.PublicMnemonicWalletData {
	return &types.PublicMnemonicWalletData{
		UUID:        u.mnemonicWalletUUID,
		IsHotWallet: u.isHotWallet,
		Hash:        u.mnemonicWalletHash,
	}
}

func (u *MnemonicWalletUnit) SignTransaction(ctx context.Context,
	account, change, index uint32,
	transaction *tronCore.Transaction,
) (*types.PublicSignTxData, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.isWalletLoaded {
		defer u.onAirTicker.Reset(u.unloadTimerInterval)
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
) (*types.PublicSignTxData, error) {
	key := fmt.Sprintf("%d'/%d/%d", account, change, index)
	addrData, isExists := u.addressPool[key]
	if !isExists {
		tronWallet, walletErr := u.hdWalletSrv.NewTronWallet(account, change, index)
		if walletErr != nil {
			return nil, walletErr
		}

		wif, walletErr := tronWallet.GetAccountWIF()
		if walletErr != nil {
			return nil, walletErr
		}

		address, walletErr := tronWallet.GetAddress()
		if walletErr != nil {
			return nil, walletErr
		}

		clonedPrivKey := *wif.PrivKey.ToECDSA()

		addrData = &addressData{
			address:    address,
			privateKey: &clonedPrivKey,
		}

		u.addressPool[key] = addrData

		// clear temporary private key
		defer zeroKeyBTCec(wif.PrivKey)
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
		signature, signErr := crypto.Sign(hash, addrData.privateKey)
		if signErr != nil {
			return nil, signErr
		}

		transaction.Signature = append(transaction.Signature, signature)
	}

	return &types.PublicSignTxData{
		WalletUUID:   u.walletEntity.WalletUUID,
		MnemonicUUID: u.mnemonicWalletUUID,
		MnemonicHash: u.walletEntity.MnemonicHash,
		SignedTx:     transaction,
		AddressData: &types.PublicDerivationAddressData{
			AccountIndex:  account,
			InternalIndex: change,
			AddressIndex:  index,
			Address:       addrData.address,
		},
	}, nil
}

func (u *MnemonicWalletUnit) GetAddressByPath(ctx context.Context,
	account, change, index uint32,
) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.isWalletLoaded {
		defer u.onAirTicker.Reset(u.unloadTimerInterval)

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
	marshallerCallback func(addressIdx, position uint32, address string),
) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.isWalletLoaded {
		defer u.onAirTicker.Reset(u.unloadTimerInterval)
		return u.getAddressesByPathByRange(ctx, accountIndex, internalIndex,
			addressIndexFrom, addressIndexTo, marshallerCallback)
	}

	err := u.loadWallet(ctx)
	if err != nil {
		return err
	}

	return u.getAddressesByPathByRange(ctx, accountIndex, internalIndex,
		addressIndexFrom, addressIndexTo, marshallerCallback)
}

func (u *MnemonicWalletUnit) getAddressesByPathByRange(ctx context.Context,
	accountIndex uint32,
	internalIndex uint32,
	addressIndexFrom uint32,
	addressIndexTo uint32,
	marshallerCallback func(addressIdx, position uint32, address string),
) error {
	var err error
	rangeSize := (addressIndexTo - addressIndexFrom) + 1
	wg := sync.WaitGroup{}
	wg.Add(int(rangeSize))

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

			marshallerCallback(i, j, address)

			return
		}(i, j)

		j++
	}

	wg.Wait()

	if err != nil {
		return err
	}

	return nil
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
		defer u.onAirTicker.Reset(u.unloadTimerInterval)
	}

	return u.loadWallet(ctx)
}

func (u *MnemonicWalletUnit) loadWallet(ctx context.Context) error {
	walletEntity, err := u.mnemonicWalletsDataSrv.GetMnemonicWalletUUID(ctx, u.mnemonicWalletUUID.String())
	if err != nil {
		return err
	}
	if walletEntity == nil {
		return ErrPassedMnemonicWalletNotFound
	}

	u.walletEntity = walletEntity

	mnemonicBytes, err := u.cryptoSrv.Decrypt(u.walletEntity.VaultEncrypted)
	if err != nil {
		return err
	}

	mnemonicSum256 := sha256.Sum256(mnemonicBytes)
	if hex.EncodeToString(mnemonicSum256[:]) != u.walletEntity.MnemonicHash {
		return ErrWrongMnemonicHash
	}

	blockChainParams := chaincfg.MainNetParams
	hdWallet, creatErr := hdwallet.NewFromString(string(mnemonicBytes), &blockChainParams)
	if creatErr != nil {
		return creatErr
	}
	u.hdWalletSrv = hdWallet

	u.isWalletLoaded = true

	u.logger.Info("wallet successfully load")

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
	u.hdWalletSrv.ClearSecrets()

	u.hdWalletSrv = nil
	u.walletEntity = nil

	for key, data := range u.addressPool {
		if data == nil {
			continue
		}

		if data.privateKey != nil {
			zeroKey(data.privateKey)
		}

		delete(u.addressPool, key)
	}

	u.isWalletLoaded = false

	u.logger.Info("wallet successfully unload")

	return nil
}

func (u *MnemonicWalletUnit) Shutdown(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if !u.isWalletLoaded {
		return nil
	}

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
	mnemonicWalletItem *entities.MnemonicWallet,
) *MnemonicWalletUnit {
	return &MnemonicWalletUnit{
		logger: logger.With(zap.String(app.WalletUUIDTag, walletUUID.String()),
			zap.String(app.MnemonicWalletUUIDTag, mnemonicWalletItem.UUID.String())),
		mu: sync.Mutex{},

		onAirTicker: nil, // that field will be field @ run stage

		hdWalletSrv: nil, // that field will be field @ load wallet stage

		cfgSrv:                 cfg,
		cryptoSrv:              cryptoSrv,
		mnemonicWalletsDataSrv: mnemonicWalletDataSrv,

		isWalletLoaded:      false,
		isHotWallet:         mnemonicWalletItem.IsHotWallet,
		walletUUID:          walletUUID,
		mnemonicWalletUUID:  mnemonicWalletItem.UUID,
		mnemonicWalletHash:  mnemonicWalletItem.MnemonicHash,
		unloadTimerInterval: unloadInterval,
		walletEntity:        mnemonicWalletItem,
		addressPool:         make(map[string]*addressData, 0),
	}
}
