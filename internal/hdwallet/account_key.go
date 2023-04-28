package hdwallet

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// AccountKey results info account keys
type AccountKey struct {
	Network     *chaincfg.Params
	ExtendedKey *hdkeychain.ExtendedKey
	Private     string
	Public      string
}

// Init account keys
func (a *AccountKey) Init() error {
	a.Private = a.ExtendedKey.String()
	w, _ := StringWallet(a.Private, a.Network.HDPrivateKeyID, a.Network.HDPublicKeyID)
	pub, err := w.Pub().String()
	if err != nil {
		return err
	}

	a.Public = pub

	return nil
}
