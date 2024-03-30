package grpc

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"
)

func (m *grpcMarshaller) MarshallSignTransaction(
	publicSignTxData *types.PublicSignTxData,
) (*pbApi.SignTransactionResponse, error) {
	return &pbApi.SignTransactionResponse{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: publicSignTxData.MnemonicUUID.String(),
			WalletHash: publicSignTxData.MnemonicHash,
		},
		TxOwnerIdentity: &pbCommon.DerivationAddressIdentity{
			AccountIndex:  publicSignTxData.AddressData.AccountIndex,
			InternalIndex: publicSignTxData.AddressData.InternalIndex,
			AddressIndex:  publicSignTxData.AddressData.AddressIndex,
			Address:       publicSignTxData.AddressData.Address,
		},
		SignedTxData: publicSignTxData.SignedTx,
	}, nil
}
