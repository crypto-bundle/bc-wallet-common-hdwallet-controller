package types

type PublicDerivationAddressData struct {
	AccountIndex  uint32
	InternalIndex uint32
	AddressIndex  uint32
	Address       string
}

type PublicDerivationAddressRangeData struct {
	AccountIndex     uint32
	InternalIndex    uint32
	AddressIndexFrom uint32
	AddressIndexTo   uint32
	AddressIndexDiff int32
}

type AddrRangeIterable interface {
	GetNext() *PublicDerivationAddressRangeData
	GetRangesSize() uint32
}
