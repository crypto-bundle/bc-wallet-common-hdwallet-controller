package types

import (
	"github.com/google/uuid"
)

type WalletMakerStrategy uint8

const (
	WalletMakerSingleMnemonicStrategyName   = "single_mnemonic_strategy"
	WalletMakerMultipleMnemonicStrategyName = "multiple_mnemonic_strategy"
)

const (
	WalletMakerPlaceholderMnemonicStrategy WalletMakerStrategy = 0
	WalletMakerSingleMnemonicStrategy      WalletMakerStrategy = 1
	WalletMakerMultipleMnemonicStrategy    WalletMakerStrategy = 2
)

func (d WalletMakerStrategy) String() string {
	switch d {
	case WalletMakerSingleMnemonicStrategy:
		return WalletMakerSingleMnemonicStrategyName
	case WalletMakerMultipleMnemonicStrategy:
		return WalletMakerMultipleMnemonicStrategyName
	case WalletMakerPlaceholderMnemonicStrategy:
		fallthrough
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
	UUID                  uuid.UUID
	Title                 string
	Purpose               string
	Strategy              WalletMakerStrategy
	MnemonicWallets       []*PublicMnemonicWalletData
	MnemonicWalletsByUUID map[uuid.UUID]*PublicMnemonicWalletData
}

type PublicSignTxData struct {
	WalletUUID   uuid.UUID
	MnemonicUUID uuid.UUID
	MnemonicHash string
	AddressData  *PublicDerivationAddressData
	SignedTx     []byte
}
