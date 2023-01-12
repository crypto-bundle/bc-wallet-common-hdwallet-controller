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

package hdwallet

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
)

var (
	// DefaultNetwork for generate masterKey
	// nolint:gochecknoglobals // its library function
	DefaultNetwork = &chaincfg.MainNetParams
	// ZeroQuote base zero
	// nolint:gochecknoglobals // its library function
	ZeroQuote uint32 = 0x80000000
)

// Key struct
type Key struct {
	// ExtendedKey hdwallet
	ExtendedKey *hdkeychain.ExtendedKey

	// Network chain params
	Network *chaincfg.Params

	// Private for btc child's
	Private *btcec.PrivateKey
	// Public for btc child's
	Public *btcec.PublicKey

	// PrivateECDSA for eth child's and tokens's
	PrivateECDSA *ecdsa.PrivateKey
	// PrivateECDSA for eth child's and tokens's
	PublicECDSA *ecdsa.PublicKey
}

// NewKey generate new extended key
func NewKey(seed []byte) (*Key, error) {
	extendedKey, err := hdkeychain.NewMaster(seed, DefaultNetwork)
	if err != nil {
		return nil, err
	}
	return newKey(extendedKey, DefaultNetwork)
}

func newKey(extendedKey *hdkeychain.ExtendedKey, network *chaincfg.Params) (*Key, error) {
	key := &Key{ExtendedKey: extendedKey, Network: network}
	if err := key.init(); err != nil {
		return nil, err
	}

	return key, nil
}

func (k *Key) init() error {
	var err error

	k.Private, err = k.ExtendedKey.ECPrivKey()
	if err != nil {
		return err
	}

	k.Public, err = k.ExtendedKey.ECPubKey()
	if err != nil {
		return err
	}

	k.PrivateECDSA = k.Private.ToECDSA()
	k.PublicECDSA = &k.PrivateECDSA.PublicKey
	return nil
}

// GetPath return path in bip44 style
func (k *Key) GetPath(purpose, coinType, account, change, addressIndex uint32) []uint32 {
	purpose = ZeroQuote + purpose
	coinType = ZeroQuote + coinType
	account = ZeroQuote + account
	return []uint32{
		purpose,
		coinType,
		account,
		change,
		addressIndex,
	}
}

// GetChildKey path for address
func (k *Key) GetChildKey(network *chaincfg.Params, purpose, coinType, account,
	change, addressIndex uint32) (*AccountKey, *Key, error) {
	var err error
	k.ExtendedKey.SetNet(network)

	extendedKey := k.ExtendedKey
	accountKey := extendedKey
	for i, v := range k.GetPath(purpose, coinType, account, change, addressIndex) {
		extendedKey, err = extendedKey.Child(v)
		if err != nil {
			return nil, nil, err
		}
		if i == 2 {
			accountKey = extendedKey
		}
	}

	acc := &AccountKey{ExtendedKey: accountKey}

	err = acc.Init(network)
	if err != nil {
		return nil, nil, err
	}

	key, err := newKey(extendedKey, network)

	return acc, key, err
}

// PublicHex generate public key to string by hex
func (k *Key) PublicHex() string {
	return hex.EncodeToString(k.Public.SerializeCompressed())
}

// PublicHash generate public key by hash160
func (k *Key) PublicHash() ([]byte, error) {
	address, err := k.ExtendedKey.Address(k.Network)
	if err != nil {
		return nil, err
	}

	return address.ScriptAddress(), nil
}

// AddressP2WPKH generate public key to p2wpkh style address
func (k *Key) AddressP2PKH() (string, error) {
	pubHash, err := k.PublicHash()
	if err != nil {
		return "", err
	}

	addr, err := btcutil.NewAddressPubKeyHash(pubHash, k.Network)
	if err != nil {
		return "", err
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", err
	}

	addr1, err := btcutil.NewAddressScriptHash(script, k.Network)
	if err != nil {
		return "", err
	}

	return addr1.EncodeAddress(), nil
}

// AddressP2WPKH generate public key to p2wpkh style address
func (k *Key) AddressP2WPKH() (string, error) {
	pubHash, err := k.PublicHash()
	if err != nil {
		return "", err
	}

	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubHash, k.Network)
	if err != nil {
		return "", err
	}

	return addr.EncodeAddress(), nil
}

// AddressP2WPKHInP2SH generate public key to p2wpkh nested within p2sh style address
func (k *Key) AddressP2WPKHInP2SH() (string, error) {
	pubHash, err := k.PublicHash()
	if err != nil {
		return "", err
	}

	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubHash, k.Network)
	if err != nil {
		return "", err
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", err
	}

	addr1, err := btcutil.NewAddressScriptHash(script, k.Network)
	if err != nil {
		return "", err
	}

	return addr1.EncodeAddress(), nil
}
