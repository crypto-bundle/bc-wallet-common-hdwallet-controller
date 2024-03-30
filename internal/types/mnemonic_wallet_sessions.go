package types

type MnemonicWalletSessionStatus uint8

const (
	MnemonicWalletSessionStatusPreparedName = "prepared"
	MnemonicWalletSessionStatusClosedName   = "closed"
)

const (
	MnemonicWalletSessionStatusPlaceholder MnemonicWalletSessionStatus = iota
	MnemonicWalletSessionStatusPrepared
	MnemonicWalletSessionStatusClosed
)

func (d MnemonicWalletSessionStatus) String() string {
	switch d {
	case MnemonicWalletSessionStatusPrepared:
		return MnemonicWalletSessionStatusPreparedName
	case MnemonicWalletSessionStatusClosed:
		return MnemonicWalletSessionStatusClosedName
	case MnemonicWalletSessionStatusPlaceholder:
		fallthrough
	default:
		return ""
	}
}
