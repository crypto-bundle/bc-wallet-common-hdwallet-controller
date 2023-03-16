/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// Package hdwallet Main
package hdwallet

import (
	"encoding/hex"
	"errors"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/config"

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
