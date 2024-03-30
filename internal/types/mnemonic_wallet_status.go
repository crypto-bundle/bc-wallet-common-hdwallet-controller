package types

type MnemonicWalletStatus uint8

const (
	MnemonicWalletStatusCreatedName  = "created"
	MnemonicWalletStatusEnabledName  = "enabled"
	MnemonicWalletStatusDisabledName = "disabled"
)

const (
	MnemonicWalletStatusPlaceholder MnemonicWalletStatus = iota
	MnemonicWalletStatusCreated
	MnemonicWalletStatusEnabled
	MnemonicWalletStatusDisabled
)

func (d MnemonicWalletStatus) String() string {
	switch d {
	case MnemonicWalletStatusCreated:
		return MnemonicWalletStatusCreatedName
	case MnemonicWalletStatusEnabled:
		return MnemonicWalletStatusEnabledName
	case MnemonicWalletStatusDisabled:
		return MnemonicWalletStatusDisabledName
	case MnemonicWalletStatusPlaceholder:
		fallthrough
	default:
		return ""
	}
}
