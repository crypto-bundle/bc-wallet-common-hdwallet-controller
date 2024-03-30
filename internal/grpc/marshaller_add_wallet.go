package grpc

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"
)

func (m *grpcMarshaller) MarshallCreateWalletData(
	walletData *types.PublicWalletData,
) (*pbApi.AddNewWalletResponse, error) {
	mnemonicsCount := uint32(len(walletData.MnemonicWallets))

	resp := &pbApi.AddNewWalletResponse{Wallet: &pbCommon.WalletData{
		Identity:            &pbCommon.WalletIdentity{WalletUUID: walletData.UUID.String()},
		Title:               walletData.Title,
		Purpose:             walletData.Purpose,
		Strategy:            pbCommon.WalletMakerStrategy(walletData.Strategy),
		MnemonicWalletCount: uint32(len(walletData.MnemonicWallets)),
		MnemonicWallets:     make([]*pbCommon.MnemonicWalletData, mnemonicsCount),
		Bookmarks: &pbCommon.WalletBookmarks{
			HotWalletIndex: 0,
		},
	}}

	for i := uint32(0); i != mnemonicsCount; i++ {
		mnemonicPublicData := walletData.MnemonicWallets[i]
		resp.Wallet.MnemonicWallets[i] = &pbCommon.MnemonicWalletData{
			Identity: &pbCommon.MnemonicWalletIdentity{
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
