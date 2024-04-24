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

	"github.com/google/uuid"
)

//go:generate easyjson mnemonic_wallet.go

// MnemonicWallet struct for storing in pg database
// easyjson:json
type MnemonicWallet struct {
	ID           uint32    `db:"id" json:"id"`
	UUID         uuid.UUID `db:"uuid" json:"uuid"`
	MnemonicHash string    `db:"mnemonic_hash" json:"mnemonic_hash"`

	Status         types.MnemonicWalletStatus `db:"status" json:"status"`
	UnloadInterval time.Duration              `db:"unload_interval" json:"unload_interval"`

	VaultEncrypted     []byte `db:"vault_encrypted" json:"vault_encrypted"`
	VaultEncryptedHash string `db:"vault_encrypted_hash" json:"vault_encrypted_hash"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

// MnemonicWalletFull struct for storing full wallet info with list of active sessions
// easyjson:json
type MnemonicWalletFull struct {
	Wallet   MnemonicWallet           `json:"wallet"`
	Sessions []*MnemonicWalletSession `json:"sessions"`
}
