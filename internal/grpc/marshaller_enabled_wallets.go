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
			walletData := walletsData[index]
			mnemonicWalletsCount := len(walletData.MnemonicWallets)

			walletIdentity := &pbApi.WalletData{
				Identity: &pbApi.WalletIdentity{
					WalletUUID: walletData.UUID.String(),
				},
				Title:               walletData.Title,
				Purpose:             walletData.Purpose,
				Strategy:            pbApi.WalletMakerStrategy(walletData.Strategy),
				MnemonicWalletCount: uint32(mnemonicWalletsCount),
				MnemonicWallets:     make([]*pbApi.MnemonicWalletData, mnemonicWalletsCount),
			}

			for j := 0; j != mnemonicWalletsCount; j++ {
				walletIdentity.MnemonicWallets[j] = &pbApi.MnemonicWalletData{
					Identity: &pbApi.MnemonicWalletIdentity{
						WalletUUID: walletData.MnemonicWallets[j].UUID.String(),
						WalletHash: walletData.MnemonicWallets[j].Hash,
					},
					IsHot: walletData.MnemonicWallets[j].IsHotWallet,
				}
			}

			response.Wallets[index] = walletIdentity

			wg.Done()
		}(i)
	}
	wg.Wait()

	return response, nil
}
