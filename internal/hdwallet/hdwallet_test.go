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
