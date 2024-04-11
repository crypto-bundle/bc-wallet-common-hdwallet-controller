package grpc

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
)

func (m *grpcMarshaller) MarshallGetAddressData(
	mnemonicWallet *entities.MnemonicWallet,
	pbAddressData *pbCommon.DerivationAddressIdentity,
) (*pbApi.DerivationAddressResponse, error) {
	return &pbApi.DerivationAddressResponse{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: mnemonicWallet.UUID.String(),
			WalletHash: mnemonicWallet.MnemonicHash,
		},
		AddressIdentity: pbAddressData,
	}, nil
}
