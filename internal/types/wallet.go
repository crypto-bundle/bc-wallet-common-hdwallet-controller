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
	IsHotWallet bool
}

type PublicWalletData struct {
	UUID            uuid.UUID
	Title           string
	Purpose         string
	Strategy        WalletMakerStrategy
	MnemonicWallets []*PublicMnemonicWalletData
}
