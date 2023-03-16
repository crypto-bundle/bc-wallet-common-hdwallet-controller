package hdwallet

import (
	"strconv"
	"strings"
)

type AddressPath struct {
	path string

	accountIndex  uint32
	internalIndex uint32
	addressIndex  uint32
}

type AddressPathParser struct {
}

func (p *AddressPathParser) Parse(path string) *AddressPath {
	idxList := ParseAddressPath(path)
	return &AddressPath{
		path:          path,
		accountIndex:  uint32(idxList[0]),
		internalIndex: uint32(idxList[1]),
		addressIndex:  uint32(idxList[2]),
	}
}

func (p *AddressPath) GetAccountIndex() uint32 {
	return p.accountIndex
}

func (p *AddressPath) GetInternalIndex() uint32 {
	return p.internalIndex
}

func (p *AddressPath) GetAddressIndex() uint32 {
	return p.addressIndex
}

// ParseAddressPath derives. Example m/49'/1'/0'/0/0
func ParseAddressPath(path string) []int {
	var data []int
	parts := strings.Split(path, "/")
	for _, part := range parts {
		// do we have an apostrophe?
		harden := part[len(part)-1:] == "'"
		if harden {
			part = part[:len(part)-1]
		}
		if idx, err := strconv.Atoi(part); err == nil {
			data = append(data, idx)
		}
	}
	return data
}
