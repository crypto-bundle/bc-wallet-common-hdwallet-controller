package hdwallet

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"
)

// ETH parent
type ETH struct {
	purpose  int
	coinType int

	account       uint32
	change        uint32
	addressNumber uint32

	*BTC
}

// NewEthWallet create new wallet
func (w *Wallet) NewEthWallet(account, change, address uint32) (*ETH, error) {
	blockChainParams := chaincfg.MainNetParams

	accountKey, key, err := w.GetChildKey(&blockChainParams, Bip44Purpose, EthCoinNumber, account, change, address)
	if err != nil {
		return nil, err
	}

	return &ETH{
		purpose:       Bip44Purpose,
		coinType:      EthCoinNumber,
		account:       account,
		change:        change,
		addressNumber: address,
		BTC: &BTC{
			ExtendedKey:      key,
			AccountKey:       accountKey,
			blockChainParams: &blockChainParams,
		},
	}, err
}

// GetAddress get address with 0x
func (e *ETH) GetAddress() (string, error) {
	return crypto.PubkeyToAddress(*e.ExtendedKey.PublicECDSA).Hex(), nil
}

// GetPubKey get key with 0x
func (e *ETH) GetPubKey() string {
	pubKey := "0x" + e.BTC.GetPubKey()
	return pubKey
}

// GetPrvKey get key with 0x
func (e *ETH) GetPrvKey() (string, error) {
	prv := hex.EncodeToString(e.ExtendedKey.PrivateECDSA.D.Bytes())
	return "0x" + prv, nil
}

// GetPath ...
func (e *ETH) GetPath() string {
	return fmt.Sprintf("m/%d'/%d'/%d'/%d/%d",
		e.GetPurpose(), e.GetCoinType(), e.account, e.change, e.addressNumber)
}

// GetPurpose ...
func (e *ETH) GetPurpose() int {
	return e.purpose
}

// GetCoinType
func (e *ETH) GetCoinType() int {
	return EthCoinNumber
}
