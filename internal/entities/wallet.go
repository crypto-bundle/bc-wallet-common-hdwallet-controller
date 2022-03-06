package entities

import (
	"time"

	"github.com/google/uuid"
)

// MnemonicWallet struct for storing in pg database
type MnemonicWallet struct {
	ID            uint32     `db:"id"`
	Title         string     `db:"title"`
	UUID          uuid.UUID  `db:"uuid"`
	Hash          string     `db:"hash"`
	Purpose       string     `db:"purpose"`
	IsHotWallet   bool       `db:"is_hot"`
	EncryptedData []byte     `db:"encrypted_data"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}
