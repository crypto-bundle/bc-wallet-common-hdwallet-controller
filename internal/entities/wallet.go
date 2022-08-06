package entities

import (
	"time"

	"github.com/google/uuid"
)

// MnemonicWallet struct for storing in pg database
type MnemonicWallet struct {
	ID              uint32     `db:"id"`
	Title           string     `db:"title"`
	UUID            uuid.UUID  `db:"wallet_uuid"`
	Hash            string     `db:"hash"`
	Purpose         string     `db:"purpose"`
	IsHotWallet     bool       `db:"is_hot"`
	IsWalletEnabled bool       `db:"is_enabled"`
	RsaEncrypted    string     `db:"rsa_encrypted"`
	CreatedAt       *time.Time `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}
