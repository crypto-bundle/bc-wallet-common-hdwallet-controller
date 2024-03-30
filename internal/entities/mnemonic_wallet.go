package entities

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
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
