package grpc

import (
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

func (m *grpcMarshaller) MarshallWalletInfo(
	walletData *types.PublicWalletData,
) *pbApi.WalletData {
	mnemonicWalletsCount := len(walletData.MnemonicWallets)

	walletInfo := &pbApi.WalletData{
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
		walletInfo.MnemonicWallets[j] = &pbApi.MnemonicWalletData{
			Identity: &pbApi.MnemonicWalletIdentity{
				WalletUUID: walletData.MnemonicWallets[j].UUID.String(),
				WalletHash: walletData.MnemonicWallets[j].Hash,
			},
			IsHot: walletData.MnemonicWallets[j].IsHotWallet,
		}
	}

	return walletInfo
}
