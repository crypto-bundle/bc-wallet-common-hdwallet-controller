package grpc

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/manager"
	"sync"
)

func (m *grpcMarshaller) MarshallGetEnabledWallets(
	walletsData []*entities.MnemonicWallet,
) (*pbApi.GetEnabledWalletsResponse, error) {
	walletCount := uint32(len(walletsData))

	response := &pbApi.GetEnabledWalletsResponse{
		WalletsCount:     walletCount,
		WalletIdentities: make([]*pbCommon.MnemonicWalletIdentity, walletCount),
	}

	wg := sync.WaitGroup{}
	wg.Add(int(walletCount))
	for i := uint32(0); i != walletCount; i++ {
		go func(index uint32) {
			defer wg.Done()

			walletData := walletsData[index]
			if walletData == nil {
				return
			}

			response.WalletIdentities[index] = &pbCommon.MnemonicWalletIdentity{
				WalletUUID: walletData.UUID.String(),
				WalletHash: walletData.MnemonicHash,
			}
		}(i)
	}
	wg.Wait()

	return response, nil
}
