package hdwallet

import (
	"log"
	"testing"
)

func TestTronPubKey(t *testing.T) {
	// test Tron
	prvMagic := [4]byte{0x04, 0x88, 0xad, 0xe4}
	pubMagic := [4]byte{0x04, 0x88, 0xb2, 0x1e}

	prvKey := "xprv9y95ZDyUd1t39HcvxQV9j5UqjhYRFsRXdoUrUVyDKCZo5fjvRDept9qLbHPg9B9jCn1vp9PJm575LCbkUyqv4ybRLRgzUV32DYbXVfMaS31"
	w, err := StringWallet(prvKey, prvMagic, pubMagic)
	if err != nil {
		t.Errorf("%s should have been nil", err.Error())
	}
	pub, err := w.Pub().String()
	if err != nil {
		t.Error(err)
	}

	log.Println(pub)
	if pub != "xpub6C8RxjWNTPSLMmhQ4S2A6DRaHjNufL9P12QTGtNpsY6mxU54xky5Rx9pSZQp6B3gqPyfetfjwwPFeqPceuxyzm78HS4dAPYjFUEwdcpTqGF" {
		t.Errorf("\n%s\nsupposed to be\n", pub)
	}
}
