package grpc

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"
)

func (m *grpcMarshaller) MarshallCreateWalletData(
	walletData *entities.MnemonicWallet,
) (*pbApi.AddNewWalletResponse, error) {

	resp := &pbApi.AddNewWalletResponse{WalletIdentity: &pbCommon.MnemonicWalletIdentity{
		WalletUUID: walletData.UUID.String(),
		WalletHash: walletData.MnemonicHash,
	}}

	return resp, nil
}
