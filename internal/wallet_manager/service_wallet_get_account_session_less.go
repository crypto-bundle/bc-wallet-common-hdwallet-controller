/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
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
 *
 */

package wallet_manager

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbHdwallet "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

func (s *Service) GetAccountWithoutSession(ctx context.Context,
	wallet *entities.MnemonicWallet,
	parameters *anypb.Any,
) (address *string, err error) {
	walletUUID := wallet.UUID.String()

	_, err = s.hdWalletClientSvc.LoadMnemonic(ctx, &pbHdwallet.LoadMnemonicRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletUUID,
			WalletHash: wallet.MnemonicHash,
		},
		TimeToLive:            uint64(time.Second * 8),
		EncryptedMnemonicData: wallet.VaultEncrypted,
	})
	if err != nil {
		s.logger.Error("unable to load mnemonic wallet by hd-wallet service", zap.Error(err),
			zap.String(app.MnemonicWalletUUIDTag, walletUUID))

		return nil, err
	}

	resp, err := s.hdWalletClientSvc.GetAccount(ctx, &pbHdwallet.GetAccountRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletUUID,
		},
		AccountIdentifier: &pbCommon.AccountIdentity{
			Parameters: parameters,
		},
	})
	if err != nil {
		grpcStatus, statusExists := status.FromError(err)
		if !statusExists {
			s.logger.Error("unable get status from error", zap.Error(ErrUnableDecodeGrpcErrorStatus))
			return nil, ErrUnableDecodeGrpcErrorStatus
		}

		switch grpcStatus.Code() {
		case codes.NotFound, codes.ResourceExhausted:
			return nil, nil

		default:
			s.logger.Error("unable to get derivation address",
				zap.Error(err),
				zap.String(app.MnemonicWalletUUIDTag, walletUUID))

			return nil, err
		}
	}

	return &resp.AccountIdentifier.Address, nil
}
