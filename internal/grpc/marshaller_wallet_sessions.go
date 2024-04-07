package grpc

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/manager"
)

func (m *grpcMarshaller) MarshallWalletSessions(
	sessionsList []*entities.MnemonicWalletSession,
) []*pbApi.SessionInfo {
	mnemonicSessionsCount := len(sessionsList)
	result := make([]*pbApi.SessionInfo, mnemonicSessionsCount)

	for j := 0; j != mnemonicSessionsCount; j++ {
		sessionItem := sessionsList[j]

		result[j] = &pbApi.SessionInfo{
			SessionIdentity: &pbApi.WalletSessionIdentity{
				SessionUUID: sessionItem.UUID,
			},
			SessionStartedAt: uint64(sessionItem.CreatedAt.Unix()),
			SessionExpiredAt: uint64(sessionItem.ExpiredAt.Unix()),
		}
	}

	return result
}
