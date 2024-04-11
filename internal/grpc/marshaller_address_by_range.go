package grpc

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
)

func (m *grpcMarshaller) MarshallGetAddressByRange(
	mnemonicWallet *entities.MnemonicWallet,
	addressesData []*pbCommon.DerivationAddressIdentity,
	size uint64,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	response := &pbApi.DerivationAddressByRangeResponse{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: mnemonicWallet.UUID.String(),
			WalletHash: mnemonicWallet.MnemonicHash,
		},
		AddressIdentitiesCount: size,
		AddressIdentities:      addressesData,
	}

	return response, nil
}
