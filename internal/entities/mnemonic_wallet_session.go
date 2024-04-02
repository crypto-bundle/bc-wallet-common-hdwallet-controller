package entities

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	"time"
)

//go:generate easyjson mnemonic_wallet_session.go

// MnemonicWalletSession struct for storing in pg database
// easyjson:json
type MnemonicWalletSession struct {
	ID   uint32 `db:"id" json:"id"`
	UUID string `db:"uuid" json:"uuid"`

	AccessTokenUUID    uint32 `db:"access_token_uuid" json:"access_token_uuid"`
	MnemonicWalletUUID string `db:"mnemonic_wallet_uuid" json:"mnemonic_wallet_uuid"`

	Status types.MnemonicWalletSessionStatus `db:"status" json:"status"`

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
