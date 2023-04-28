package grpc

import (
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

func (m *grpcMarshaller) MarshallSignTransaction(
	publicSignTxData *types.PublicSignTxData,
) (*pbApi.SignTransactionResponse, error) {
	return &pbApi.SignTransactionResponse{
		WalletIdentity: &pbApi.WalletIdentity{
			WalletUUID: publicSignTxData.WalletUUID.String(),
		},
		MnemonicIdentity: &pbApi.MnemonicWalletIdentity{
			WalletUUID: publicSignTxData.MnemonicUUID.String(),
			WalletHash: publicSignTxData.MnemonicHash,
		},
		TxOwnerIdentity: &pbApi.DerivationAddressIdentity{
			AccountIndex:  publicSignTxData.AddressData.AccountIndex,
			InternalIndex: publicSignTxData.AddressData.InternalIndex,
			AddressIndex:  publicSignTxData.AddressData.AddressIndex,
			Address:       publicSignTxData.AddressData.Address,
		},
		SignedTxData: publicSignTxData.SignedTx,
	}, nil
}
