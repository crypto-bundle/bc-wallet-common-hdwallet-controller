package grpc

import (
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
	"google.golang.org/protobuf/proto"
)

func (m *grpcMarshaller) MarshallSignTransaction(
	publicSignTxData *types.PublicSignTxData,
) (*pbApi.SignTransactionResponse, error) {
	signedTxRaw, err := proto.Marshal(publicSignTxData.SignedTx)
	if err != nil {
		return nil, err
	}

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
		SignedTxData: signedTxRaw,
	}, nil
}
