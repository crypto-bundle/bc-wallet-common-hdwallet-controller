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
	"fmt"

	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shengdoushi/base58"
)

// Tron parent
type Tron struct {
	purpose  int
	coinType int

	account       uint32
	change        uint32
	addressNumber uint32

	*BTC
}

// NewTronWallet create new wallet
func (w *Wallet) NewTronWallet(account, change, address uint32) (*Tron, error) {
	accountKey, key, err := w.GetChildKey(Bip44Purpose,
		TronCoinNumber, account, change, address)
	if err != nil {
		return nil, err
	}

	return &Tron{
		purpose:       Bip44Purpose,
		coinType:      TronCoinNumber,
		account:       account,
		change:        change,
		addressNumber: address,
		BTC: &BTC{
			ExtendedKey:      key,
			AccountKey:       accountKey,
			blockChainParams: w.Network,
		},
	}, err
}

// GetAddress get address with 0x
func (e *Tron) GetAddress() (string, error) {
	addr := crypto.PubkeyToAddress(*e.ExtendedKey.PublicECDSA)

	addrTrxBytes := make([]byte, 0)
	addrTrxBytes = append(addrTrxBytes, TronBytePrefix)
	addrTrxBytes = append(addrTrxBytes, addr.Bytes()...)

	crc := calcCheckSum(addrTrxBytes)

	addrTrxBytes = append(addrTrxBytes, crc...)

	//nolint:gocritic // ok. Just reminder hot to generate address_hex for TronKey table
	// addrTrxHex := hex.EncodeToString(addrTrxBytes)

	addrTrx := base58.Encode(addrTrxBytes, base58.BitcoinAlphabet)

	return addrTrx, nil
}

// GetPubKey get key with 0x
func (e *Tron) GetPubKey() string {
	return e.BTC.GetPubKey()
}

// GetPrvKey get key with 0x
func (e *Tron) GetPrvKey() (string, error) {
	return hex.EncodeToString(e.ExtendedKey.PrivateECDSA.D.Bytes()), nil
}

// GetPath ...
func (e *Tron) GetPath() string {
	return fmt.Sprintf("m/%d'/%d'/%d'/%d/%d",
		e.GetPurpose(), e.GetCoinType(), e.account, e.change, e.addressNumber)
}

// GetPurpose ...
func (e *Tron) GetPurpose() int {
	return e.purpose
}

// GetCoinType ...
func (e *Tron) GetCoinType() int {
	return TronCoinNumber
}
