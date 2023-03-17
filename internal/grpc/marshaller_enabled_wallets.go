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
		Wallets:      make([]*pbApi.WalletIdentity, walletCount),
		WalletsCount: walletCount,
	}

	wg := sync.WaitGroup{}
	wg.Add(int(walletCount))
	for i := uint32(0); i != walletCount; i++ {
		go func(index uint32) {
			walletData := walletsData[index]
			mnemonicWalletsCount := len(walletData.MnemonicWallets)

			walletIdentity := &pbApi.WalletIdentity{
				WalletUUID:             walletData.UUID.String(),
				Title:                  walletData.Title,
				Purpose:                walletData.Purpose,
				Strategy:               pbApi.WalletMakerStrategy(walletData.Strategy),
				MnemonicWalletCount:    uint32(mnemonicWalletsCount),
				MnemonicWalletIdentity: make([]*pbApi.MnemonicWalletIdentity, mnemonicWalletsCount),
			}

			for j := 0; j != mnemonicWalletsCount; j++ {
				walletIdentity.MnemonicWalletIdentity[j] = &pbApi.MnemonicWalletIdentity{
					WalletUUID: walletData.MnemonicWallets[j].UUID.String(),
					IsHot:      walletData.MnemonicWallets[j].IsHotWallet,
				}
			}

			response.Wallets[index] = walletIdentity

			wg.Done()
		}(i)
	}
	wg.Wait()

	return response, nil
}