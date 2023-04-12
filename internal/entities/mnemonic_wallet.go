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

package entities

import (
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson mnemonic_wallet.go

// MnemonicWallet struct for storing in pg database
// easyjson:json
type MnemonicWallet struct {
	ID             uint32        `db:"id" json:"id"`
	UUID           uuid.UUID     `db:"uuid" json:"uuid"`
	WalletUUID     uuid.UUID     `db:"wallet_uuid" json:"wallet_uuid"`
	MnemonicHash   string        `db:"mnemonic_hash" json:"mnemonic_hash"`
	IsHotWallet    bool          `db:"is_hot" json:"is_hot"`
	UnloadInterval time.Duration `db:"unload_interval" json:"unload_interval"`

	RsaEncrypted       []byte `db:"rsa_encrypted" json:"rsa_encrypted"`
	RsaEncryptedHash   string `db:"rsa_encrypted_hash" json:"rsa_encrypted_hash"`
	VaultEncrypted     []byte `db:"vault_encrypted" json:"vault_encrypted"`
	VaultEncryptedHash string `db:"vault_encrypted_hash" json:"vault_encrypted_hash"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}
