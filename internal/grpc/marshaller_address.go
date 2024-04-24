package grpc

import (
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

func (m *grpcMarshaller) MarshallGetAddressData(
	walletPublicData *types.PublicWalletData,
	mnemonicWalletPublicData *types.PublicMnemonicWalletData,
	pbAddressData *pbApi.DerivationAddressIdentity,
) (*pbApi.DerivationAddressResponse, error) {
	return &pbApi.DerivationAddressResponse{
		WalletIdentity: &pbApi.WalletIdentity{
			WalletUUID: walletPublicData.UUID.String(),
		},
		MnemonicIdentity: &pbApi.MnemonicWalletIdentity{
			WalletUUID: mnemonicWalletPublicData.UUID.String(),
			WalletHash: mnemonicWalletPublicData.Hash,
		},
		AddressIdentity: pbAddressData,
	}, nil
}
