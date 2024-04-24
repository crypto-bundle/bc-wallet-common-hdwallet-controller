/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
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
 *
 */

package entities

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"time"
)

//go:generate easyjson mnemonic_wallet_session.go

// MnemonicWalletSession struct for storing in pg database
// easyjson:json
type MnemonicWalletSession struct {
	ID   uint32 `db:"id" json:"id"`
	UUID string `db:"uuid" json:"uuid"`

	AccessTokenUUID    string `db:"access_token_uuid" json:"access_token_uuid"`
	MnemonicWalletUUID string `db:"mnemonic_wallet_uuid" json:"mnemonic_wallet_uuid"`

	Status types.MnemonicWalletSessionStatus `db:"status" json:"status"`

	StartedAt time.Time `db:"started_at" json:"started_at"`
	ExpiredAt time.Time `db:"expired_at" json:"expired_at"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

func (v *MnemonicWalletSession) IsSessionActive() bool {
	if v == nil {
		return false
	}

	if v.ExpiredAt.Before(time.Now()) {
		return false
	}

	if v.Status != types.MnemonicWalletSessionStatusPrepared {
		return false
	}

	return true
}
