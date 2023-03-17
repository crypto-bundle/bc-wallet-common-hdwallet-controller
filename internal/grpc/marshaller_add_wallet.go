package grpc

import (
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"

	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

func (m *grpcMarshaller) MarshallCreateWalletData(
	walletData *types.PublicWalletData,
) (*pbApi.AddNewWalletResponse, error) {
	mnemonicsCount := uint32(len(walletData.MnemonicWallets))

	resp := &pbApi.AddNewWalletResponse{Wallet: &pbApi.WalletIdentity{
		WalletUUID:             walletData.UUID.String(),
		Title:                  walletData.Title,
		Purpose:                walletData.Purpose,
		Strategy:               pbApi.WalletMakerStrategy(walletData.Strategy),
		MnemonicWalletCount:    uint32(len(walletData.MnemonicWallets)),
		MnemonicWalletIdentity: make([]*pbApi.MnemonicWalletIdentity, mnemonicsCount),
	}}

	for i := uint32(0); i != mnemonicsCount; i++ {
		mnemonicPublicData := walletData.MnemonicWallets[i]
		resp.Wallet.MnemonicWalletIdentity[i] = &pbApi.MnemonicWalletIdentity{
			WalletUUID: mnemonicPublicData.UUID.String(),
			IsHot:      mnemonicPublicData.IsHotWallet,
		}
	}

	return resp, nil
}
