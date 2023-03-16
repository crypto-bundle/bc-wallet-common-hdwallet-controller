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

package grpc

import (
	"context"
	"sync"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	"github.com/crypto-bundle/bc-wallet-common/pkg/tracer"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodGetEnabledWallets = "GetEnabledWallets"
)

type GetEnabledWalletsHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv getEnabledWalletsMarshallerService
}

// nolint:funlen // fixme
func (h *GetEnabledWalletsHandler) Handle(ctx context.Context,
	req *pbApi.GetEnabledWalletsRequest,
) (*pbApi.GetEnabledWalletsResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	wallets, err := h.walletSrv.GetEnabledWallets(ctx)
	if err != nil {
		h.l.Error("unable to get enabled mnemonic wallets", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	walletCount := uint32(len(wallets))

	response := &pbApi.GetEnabledWalletsResponse{
		Wallets:      make([]*pbApi.WalletIdentity, walletCount),
		WalletsCount: walletCount,
	}

	wg := sync.WaitGroup{}
	wg.Add(int(walletCount))
	for i := uint32(0); i != walletCount; i++ {
		go func(index uint32) {
			walletData := wallets[index]
			mnemonicWalletsCount := len(walletData.MnemonicWallets)

			walletIdentity := &pbApi.WalletIdentity{
				WalletUUID:             walletData.UUID.String(),
				Title:                  walletData.Title,
				Purpose:                walletData.Purpose,
				Strategy:               pbApi.WalletMakerStrategy(walletData.Strategy),
				MnemonicWalletCount:    uint32(mnemonicWalletsCount),
				MnemonicWalletIdentity: make([]*pbApi.MnemonicWalletIdentity, mnemonicWalletsCount),
			}

			for j := 0; j != mnemonicWalletsCount; j++ {
				walletIdentity.MnemonicWalletIdentity[j] = &pbApi.MnemonicWalletIdentity{
					WalletUUID: walletData.MnemonicWallets[j].UUID.String(),
					IsHot:      walletData.MnemonicWallets[j].IsHotWallet,
				}
			}

			response.Wallets[index] = walletIdentity

			wg.Done()
		}(i)
	}
	wg.Wait()

	return response, nil
}

func MakeGetEnabledWalletsHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
) *GetEnabledWalletsHandler {
	return &GetEnabledWalletsHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodGetEnabledWallets)),
		walletSrv: walletSrv,
	}
}
