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
