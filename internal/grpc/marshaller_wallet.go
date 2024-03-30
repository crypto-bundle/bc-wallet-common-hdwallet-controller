package grpc

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
)

func (m *grpcMarshaller) MarshallWalletInfo(
	walletData *types.PublicWalletData,
) *pbApi.MnemonicWalletData {
	mnemonicWalletsCount := len(walletData.MnemonicWallets)
	nonMnemonicWalletCount := mnemonicWalletsCount - 1

	walletInfo := &pbApi.WalletData{
		Identity: &pbApi.MnemonicWalletIdentity{
			WalletUUID: walletData.UUID.String(),
		},
		Title:           walletData.Title,
		Purpose:         walletData.Purpose,
		MnemonicWallets: make([]*pbApi.MnemonicWalletData, mnemonicWalletsCount),
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
