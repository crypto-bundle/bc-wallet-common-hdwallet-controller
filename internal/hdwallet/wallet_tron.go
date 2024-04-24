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
