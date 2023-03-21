package entities

import (
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	"time"

	"github.com/google/uuid"
)

// Wallet struct for storing in pg database
type Wallet struct {
	ID        uint32                    `db:"id"`
	Title     string                    `db:"title"`
	UUID      uuid.UUID                 `db:"uuid"`
	Purpose   string                    `db:"purpose"`
	IsEnabled bool                      `db:"is_enabled"`
	Strategy  types.WalletMakerStrategy `db:"strategy"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
