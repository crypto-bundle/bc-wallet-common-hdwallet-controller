package grpc

import (
	"sync"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/types"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
)

func (m *grpcMarshaller) MarshallGetAddressByRange(
	walletPublicData *types.PublicWalletData,
	mnemonicWalletPublicData *types.PublicMnemonicWalletData,
	addressesData []*types.PublicDerivationAddressData,
) (*pbApi.DerivationAddressByRangeResponse, error) {
	rangeSize := uint32(len(addressesData))

	response := &pbApi.DerivationAddressByRangeResponse{
		WalletIdentity: &pbApi.WalletIdentity{
			WalletUUID: walletPublicData.UUID.String(),
		},
		MnemonicIdentity: &pbApi.MnemonicWalletIdentity{
			WalletUUID: mnemonicWalletPublicData.UUID.String(),
			WalletHash: mnemonicWalletPublicData.Hash,
		},
		AddressIdentities: make([]*pbApi.DerivationAddressIdentity, rangeSize+1),
	}

	wg := sync.WaitGroup{}
	wg.Add(int(rangeSize))
	for i := uint32(0); i != rangeSize; i++ {
		go func(index uint32) {
			addrData := addressesData[index]

			response.AddressIdentities[index] = &pbApi.DerivationAddressIdentity{
				AccountIndex:  addrData.AccountIndex,
				InternalIndex: addrData.InternalIndex,
				AddressIndex:  addrData.AddressIndex,
				Address:       addrData.Address,
			}

			wg.Done()
		}(i)
	}
	wg.Wait()

	return response, nil
}
