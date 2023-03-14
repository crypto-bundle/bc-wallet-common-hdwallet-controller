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

package wallet

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/entities"
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/hdwallet"
)

type config interface {
}

type repo interface {
	AddNewMnemonicWallet(ctx context.Context, wallet *entities.MnemonicWallet) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletUUID(ctx context.Context, uuid string) (*entities.MnemonicWallet, error)
	GetAllEnabledMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetAllEnabledHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetAllEnabledNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
}

type mnemonicWalletConfig interface {
	GetMnemonicWalletPurpose() string
	GetMnemonicWalletHash() string
	IsHotWallet() bool
}

type walleter interface {
	GetAddress() (string, error)
	GetPubKey() string
	GetPrvKey() (string, error)
	GetPath() string
}

type hdWalleter interface {
	PublicHex() string
	PublicHash() ([]byte, error)

	NewEthWallet(account, change, address uint32) (*hdwallet.Tron, error)
}

type mnemonicGenerator interface {
	Generate(ctx context.Context) (string, error)
}

type crypto interface {
	Encrypt(msg string) (string, error)
	Decrypt(encMsg string) ([]byte, error)
}
