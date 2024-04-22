package wallet_manager

import (
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
)

type sessionsByWallet map[string][]string

type sessionsByWalletDataAdapter struct {
	sessionByWallet map[string][]string
	sessionUUIDs    []string
}

func (m *sessionsByWalletDataAdapter) Marshall(item *entities.MnemonicWalletSession) error {
	bucket, isExisted := m.sessionByWallet[item.MnemonicWalletUUID]
	if !isExisted {
		m.sessionByWallet[item.MnemonicWalletUUID] = make([]string, 1)

		m.sessionByWallet[item.MnemonicWalletUUID][0] = item.UUID
	}

	m.sessionByWallet[item.MnemonicWalletUUID] = append(bucket, item.UUID)
	m.sessionUUIDs = append(m.sessionUUIDs, item.UUID)

	return nil
}

func (m *sessionsByWalletDataAdapter) GetGroupedSessions() sessionsByWallet {
	return m.sessionByWallet
}

func (m *sessionsByWalletDataAdapter) GetSessionsUUIDs() []string {
	return m.sessionUUIDs
}

func newSessionsByWalletDataMapper() *sessionsByWalletDataAdapter {
	return &sessionsByWalletDataAdapter{
		sessionByWallet: make(map[string][]string),
	}
}
