package hdwallet

import (
	"bc-wallet-eth-hdwallet/internal/app"
	"bc-wallet-eth-hdwallet/internal/config"
)

type ManagerService struct {
	wallets map[string]*Wallet
}

func NewManagerService(cfg confiManagerer) (*ManagerService, error) {
	mnemoWalletConf := cfg.GetMnemonicsWallets()

	srv := &ManagerService{
		wallets: make(map[string]*Wallet, len(mnemoWalletConf)),
	}

	for walletName, mnemoConf := range mnemoWalletConf {
		err := srv.RegisterWallet(walletName, mnemoConf)
		if err != nil {
			return nil, err
		}
	}

	return srv, nil
}

func (s *ManagerService) RegisterWallet(walletName string, mnemoConf *config.MnemonicConfig) error {
	walletRaw, err := Restore(mnemoConf.Mnemonic)
	if err != nil {
		return err
	}

	s.wallets[walletName] = walletRaw

	return nil
}

// NewEthUserWallet create new ETH wallet
func (s *ManagerService) NewEthUserWallet(account, change, addressNumber uint32) (*ETH, error) {
	hdWallet, ok := s.wallets[app.WalletMnemoName]
	if !ok {
		return nil, ErrUnsupportedBlockchain
	}

	return hdWallet.NewEthWallet(account, change, addressNumber)
}
