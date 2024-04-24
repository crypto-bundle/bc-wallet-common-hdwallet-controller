// Package hdwallet Main
package hdwallet

import (
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	bip39 "github.com/tyler-smith/go-bip39"
)

const (
	Bip49Purpose   = 49
	Bip44Purpose   = 44
	DefaultPurpose = Bip49Purpose

	BtcCoinNumber = 0

	TronCoinNumber = 195
	TronBytePrefix = byte(0x41)
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
func Create(network *chaincfg.Params) (*Wallet, error) {
	entropy, _ := bip39.NewEntropy(256)
	return New(entropy, network)
}

// New hdwallet via entropy
func New(entropy []byte, network *chaincfg.Params) (*Wallet, error) {
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return Restore(mnemonic, network)
}

// NewFromString hdwallet via mnemo string
func NewFromString(mnemo string, network *chaincfg.Params) (*Wallet, error) {
	entropy, err := bip39.EntropyFromMnemonic(mnemo)
	if err != nil {
		return nil, err
	}

	mnemonic, _ := bip39.NewMnemonic(entropy)
	return Restore(mnemonic, network)
}

// Restore mnemonic a Bip32 HD hdwallet for the mnemonic
func Restore(mnemonic string, network *chaincfg.Params) (*Wallet, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	key, err := NewKey(seed)
	if err != nil {
		return nil, err
	}

	key.ExtendedKey.SetNet(network)

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

// ClearSecrets is function clear sensitive secrets data
func (w *Wallet) ClearSecrets() {
	w.mnemonic = "0"

	pattern := []byte{0x1, 0x2, 0x3, 0x4}
	// Copy the pattern into the start of the container
	copy(w.seed, pattern)
	// Incrementally duplicate the pattern throughout the container
	for j := len(pattern); j < len(w.seed); j *= 2 {
		copy(w.seed[j:], w.seed[:j])
	}
}
