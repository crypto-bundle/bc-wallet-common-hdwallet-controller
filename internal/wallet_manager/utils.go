package wallet_manager

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/btcec"
)

func zeroKey(key *ecdsa.PrivateKey) {
	d := key.D.Bits()
	for i := range d {
		d[i] = 0
	}

	x := key.X.Bits()
	for i := range x {
		x[i] = 0
	}

	y := key.Y.Bits()
	for i := range y {
		y[i] = 0
	}
}

func zeroPubKey(key *ecdsa.PublicKey) {
	x := key.X.Bits()
	for i := range x {
		x[i] = 0
	}

	y := key.Y.Bits()
	for i := range y {
		y[i] = 0
	}
}

func zeroKeyBTCec(key *btcec.PrivateKey) {
	zeroKey(key.ToECDSA())
}
