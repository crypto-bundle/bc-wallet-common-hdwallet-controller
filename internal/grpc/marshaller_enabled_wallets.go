package grpc

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
)

func (m *grpcMarshaller) MarshallGetEnabledWallets(
	walletsData []*entities.MnemonicWallet,
) (*pbApi.GetEnabledWalletsResponse, error) {
	walletCount := uint32(len(walletsData))

	response := &pbApi.GetEnabledWalletsResponse{
		WalletsCount:     walletCount,
		WalletIdentities: make([]*pbCommon.MnemonicWalletIdentity, walletCount),
		Bookmarks:        make(map[string]uint32, walletCount),
	}

	for i := uint32(0); i != walletCount; i++ {
		walletData := walletsData[i]
		if walletData == nil {
			continue
		}

		walletUUID := walletData.UUID.String()

		response.WalletIdentities[i] = &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletUUID,
			WalletHash: walletData.MnemonicHash,
		}
		response.Bookmarks[walletUUID] = i
	}

	return response, nil
}
