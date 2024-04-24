package grpc

import (
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

func (m *grpcMarshaller) MarshallCreateWalletData(
	walletData *types.PublicWalletData,
) (*pbApi.AddNewWalletResponse, error) {
	mnemonicsCount := uint32(len(walletData.MnemonicWallets))

	resp := &pbApi.AddNewWalletResponse{Wallet: &pbApi.WalletData{
		Identity:            &pbApi.WalletIdentity{WalletUUID: walletData.UUID.String()},
		Title:               walletData.Title,
		Purpose:             walletData.Purpose,
		Strategy:            pbApi.WalletMakerStrategy(walletData.Strategy),
		MnemonicWalletCount: uint32(len(walletData.MnemonicWallets)),
		MnemonicWallets:     make([]*pbApi.MnemonicWalletData, mnemonicsCount),
		Bookmarks: &pbApi.WalletBookmarks{
			HotWalletIndex: 0,
		},
	}}

	for i := uint32(0); i != mnemonicsCount; i++ {
		mnemonicPublicData := walletData.MnemonicWallets[i]
		resp.Wallet.MnemonicWallets[i] = &pbApi.MnemonicWalletData{
			Identity: &pbApi.MnemonicWalletIdentity{
				WalletUUID: mnemonicPublicData.UUID.String(),
				WalletHash: mnemonicPublicData.Hash,
			},
			IsHot: mnemonicPublicData.IsHotWallet,
		}

		if mnemonicPublicData.IsHotWallet {
			resp.Wallet.Bookmarks.HotWalletIndex = i
		}
	}

	return resp, nil
}
