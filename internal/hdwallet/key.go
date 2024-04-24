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
func (k *Key) GetChildKey(purpose, coinType,
	account,
	change,
	addressIndex uint32,
) (*AccountKey, *Key, error) {
	var err error
	//k.ExtendedKey.SetNet(network)

	extendedKeyCloned, err := k.ExtendedKey.CloneWithVersion(k.Network.HDPrivateKeyID[:])
	if err != nil {
		return nil, nil, err
	}

	//extendedKey.SetNet(network)
	accountKey := extendedKeyCloned
	for i, v := range k.GetPath(purpose, coinType, account, change, addressIndex) {
		extendedKey, loopErr := extendedKeyCloned.Derive(v)
		if loopErr != nil {
			return nil, nil, loopErr
		}

		if i == 2 {
			accountKey = extendedKey
		}
	}

	acc := &AccountKey{
		ExtendedKey: accountKey,
		Network:     k.Network,
	}

	err = acc.Init()
	if err != nil {
		return nil, nil, err
	}

	key, err := newKey(extendedKeyCloned, k.Network)

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

// AddressP2PKH generate public key to p2wpkh style address
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

// CloneECDSAPrivateKey generate public key to p2wpkh style address
func (k *Key) CloneECDSAPrivateKey() (*ecdsa.PrivateKey, error) {
	clonedX := *k.PrivateECDSA.X
	clonedY := *k.PrivateECDSA.Y
	clonedD := *k.PrivateECDSA.D

	clonedPrivKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: btcec.S256(),
			X:     &clonedX,
			Y:     &clonedY,
		},
		D: &clonedD,
	}

	return &clonedPrivKey, nil
}
