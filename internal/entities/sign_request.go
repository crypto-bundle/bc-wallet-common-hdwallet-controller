package entities

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	"github.com/lib/pq"
	"time"
)

//go:generate easyjson sign_request.go

// SignRequest struct for storing in pg database
// easyjson:json
type SignRequest struct {
	ID   uint32 `db:"id" json:"id"`
	UUID string `db:"uuid" json:"uuid"`

	WalletUUID  string `db:"wallet_uuid" json:"wallet_uuid"`
	SessionUUID string `db:"session_uuid" json:"session_uuid"`
	PurposeUUID string `db:"purpose_uuid" json:"purpose_uuid"`

	Status types.SignRequestStatus `db:"status" json:"status"`

	DerivationPath pq.Int32Array `db:"derivation_path" json:"derivation_path,omitempty"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
