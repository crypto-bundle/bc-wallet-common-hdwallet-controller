/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package hdwallet

import (
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/config"
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
