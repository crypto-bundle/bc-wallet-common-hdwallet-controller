package types

type PublicDerivationAddressData struct {
	PublicWallet   *PublicWalletData
	MnemonicWallet *PublicMnemonicWalletData

	AccountIndex  uint32
	InternalIndex uint32
	AddressIndex  uint32
	Address       string
}
