// Package hdwallet Main
package hdwallet

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"bc-wallet-eth-hdwallet/internal/config"

	"github.com/btcsuite/btcd/chaincfg"
	bip39 "github.com/tyler-smith/go-bip39"
)

const (
	Bip49Purpose   = 49
	Bip44Purpose   = 44
	DefaultPurpose = Bip49Purpose

	BtcCoinNumber = 0

	EthCoinNumber = 60
)

var (
	ErrUnsupportedBlockchain = errors.New("unsupported blockchain")
)

// Wallet contains the individual mnemonic and seed
type Wallet struct {
	mnemonic string
	seed     []byte
	*Key
}

type confiManagerer interface {
	IsModeDev() bool
	IsModeTest() bool

	GetMnemonicsWallets() map[string]*config.MnemonicConfig
}

// BTC parent
type BTC struct {
	blockChainParams *chaincfg.Params
	purpose          int
	coinType         int
	account          uint32
	change           uint32
	addressNumber    uint32

	ExtendedKey *Key
	AccountKey  *AccountKey
}

// Create a mnemonic for memorization or user-friendly seeds
func Create() (*Wallet, error) {
	entropy, _ := bip39.NewEntropy(256)
	return New(entropy)
}

// New hdwallet via entropy
func New(entropy []byte) (*Wallet, error) {
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return Restore(mnemonic)
}

// NewFromString hdwallet via mnemo string
func NewFromString(mnemo string) (*Wallet, error) {
	entropy, err := bip39.EntropyFromMnemonic(mnemo)
	if err != nil {
		return nil, err
	}

	mnemonic, _ := bip39.NewMnemonic(entropy)
	return Restore(mnemonic)
}

// Restore mnemonic a Bip32 HD hdwallet for the mnemonic
func Restore(mnemonic string) (*Wallet, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	key, err := NewKey(seed)
	if err != nil {
		return nil, err
	}

	return &Wallet{mnemonic, seed, key}, nil
}

// Seed return seed
func (w *Wallet) Seed() []byte {
	return w.seed
}

// GetSeed return string of seed from byte
func (w *Wallet) GetSeed() string {
	return hex.EncodeToString(w.Seed())
}

// GetMnemonic return mnemonic string
func (w *Wallet) GetMnemonic() string {
	return w.mnemonic
}

// ParseAddressPath derives. Example m/49'/1'/0'/0/0
func ParseAddressPath(path string) []int {
	var data []int
	parts := strings.Split(path, "/")
	for _, part := range parts {
		// do we have an apostrophe?
		harden := part[len(part)-1:] == "'"
		if harden {
			part = part[:len(part)-1]
		}
		if idx, err := strconv.Atoi(part); err == nil {
			data = append(data, idx)
		}
	}
	return data
}
