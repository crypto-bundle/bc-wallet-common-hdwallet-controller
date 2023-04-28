package grpc

import (
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

func (m *grpcMarshaller) MarshallGetAddressByRange(
	walletPublicData *types.PublicWalletData,
	mnemonicWalletPublicData *types.PublicMnemonicWalletData,
	addressesData []*pbApi.DerivationAddressIdentity,
	size uint64,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	response := &pbApi.DerivationAddressByRangeResponse{
		WalletIdentity: &pbApi.WalletIdentity{
			WalletUUID: walletPublicData.UUID.String(),
		},
		MnemonicIdentity: &pbApi.MnemonicWalletIdentity{
			WalletUUID: mnemonicWalletPublicData.UUID.String(),
			WalletHash: mnemonicWalletPublicData.Hash,
		},
		AddressIdentitiesCount: size,
		AddressIdentities:      addressesData,
	}

	return response, nil
}
