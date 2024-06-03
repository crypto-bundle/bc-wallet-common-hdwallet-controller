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

package controller

import (
	"context"
	"fmt"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

func (s *ManagerApiClientWrapper) GetWalletAccount(ctx context.Context,
	walletUUID string,
	accountData *anypb.Any,
) (*GetAccountResponse, error) {
	req := &GetAccountRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletUUID,
		},
		SessionIdentifier: nil,
		AccountIdentifier: &pbCommon.AccountIdentity{
			Parameters: accountData,
		},
	}

	resp, err := s.originGRPCClient.GetAccount(ctx, req)
	if err != nil {
		grpcStatus, statusExists := status.FromError(err)
		if !statusExists {
			s.logger.Error("unable get status from error", zap.Error(err))

			return nil, fmt.Errorf("%w: %s", ErrUnableDecodeGrpcErrorStatus, err)
		}

		switch grpcStatus.Code() {
		case codes.NotFound:
			s.logger.Error("wallet not found", zap.Error(err))

			return nil, nil
		case codes.ResourceExhausted:
			s.logger.Error("unable to get account, session expired or wallet not loaded", zap.Error(err))

			return nil, nil

		default:
			return nil, err
		}
	}

	if resp == nil {
		return nil, ErrMissingResponse
	}

	return resp, nil
}
