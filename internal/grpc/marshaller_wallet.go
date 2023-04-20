package grpc

import (
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

func (m *grpcMarshaller) MarshallWalletInfo(
	walletData *types.PublicWalletData,
) *pbApi.WalletData {
	mnemonicWalletsCount := len(walletData.MnemonicWallets)
	nonMnemonicWalletCount := mnemonicWalletsCount - 1

	walletInfo := &pbApi.WalletData{
		Identity: &pbApi.WalletIdentity{
			WalletUUID: walletData.UUID.String(),
		},
		Title:               walletData.Title,
		Purpose:             walletData.Purpose,
		Strategy:            pbApi.WalletMakerStrategy(walletData.Strategy),
		MnemonicWalletCount: uint32(mnemonicWalletsCount),
		MnemonicWallets:     make([]*pbApi.MnemonicWalletData, mnemonicWalletsCount),
		Bookmarks: &pbApi.WalletBookmarks{
			HotWalletIndex:      0,
			NonHotWalletIndexes: make([]uint32, nonMnemonicWalletCount),
		},
	}

	for j := 0; j != mnemonicWalletsCount; j++ {
		mnemoWallet := walletData.MnemonicWallets[j]

		walletInfo.MnemonicWallets[j] = &pbApi.MnemonicWalletData{
			Identity: &pbApi.MnemonicWalletIdentity{
				WalletUUID: mnemoWallet.UUID.String(),
				WalletHash: mnemoWallet.Hash,
			},
			IsHot: mnemoWallet.IsHotWallet,
		}

		if mnemoWallet.IsHotWallet {
			walletInfo.Bookmarks.HotWalletIndex = uint32(j)

			continue
		}

		walletInfo.Bookmarks.NonHotWalletIndexes[nonMnemonicWalletCount-1] = uint32(j)
		nonMnemonicWalletCount--
	}

	return walletInfo
}
