package entities

import (
	"github.com/google/uuid"
	"time"
)

//go:generate easyjson access_token.go

// AccessToken struct for storing in pg database
// easyjson:json
type AccessToken struct {
	ID   uint32    `db:"id" json:"id"`
	UUID uuid.UUID `db:"uuid" json:"uuid"`

	WalletUUID uuid.UUID `db:"wallet_uuid" json:"wallet_uuid"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	ExpiredAt *time.Time `db:"expired_at" json:"expired_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}
