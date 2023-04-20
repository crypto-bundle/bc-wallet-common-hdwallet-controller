package grpc

import (
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
	"sync"
)

func (m *grpcMarshaller) MarshallGetEnabledWallets(
	walletsData []*types.PublicWalletData,
) (*pbApi.GetEnabledWalletsResponse, error) {
	walletCount := uint32(len(walletsData))

	response := &pbApi.GetEnabledWalletsResponse{
		WalletsCount: walletCount,
		Wallets:      make([]*pbApi.WalletData, walletCount),
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

			walletInfo := m.MarshallWalletInfo(walletData)

			response.Wallets[index] = walletInfo
		}(i)
	}
	wg.Wait()

	return response, nil
}
