package types

import (
	tronCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/google/uuid"
)

type WalletMakerStrategy uint8

const (
	WalletMakerSingleMnemonicStrategyName   = "single_mnemonic_strategy"
	WalletMakerMultipleMnemonicStrategyName = "multiple_mnemonic_strategy"
)

const (
	WalletMakerSingleMnemonicStrategy   WalletMakerStrategy = 0
	WalletMakerMultipleMnemonicStrategy WalletMakerStrategy = 1
)

func (d WalletMakerStrategy) String() string {
	switch d {
	case WalletMakerSingleMnemonicStrategy:
		return WalletMakerSingleMnemonicStrategyName
	case WalletMakerMultipleMnemonicStrategy:
		return WalletMakerMultipleMnemonicStrategyName
	default:
		return ""
	}
}

type PublicMnemonicWalletData struct {
	UUID        uuid.UUID
	Hash        string
	IsHotWallet bool
}

type PublicWalletData struct {
	UUID            uuid.UUID
	Title           string
	Purpose         string
	Strategy        WalletMakerStrategy
	MnemonicWallets []*PublicMnemonicWalletData
}

type PublicSignTxData struct {
	WalletUUID   uuid.UUID
	MnemonicUUID uuid.UUID
	MnemonicHash string
	AddressData  *PublicDerivationAddressData
	SignedTx     *tronCore.Transaction
}
