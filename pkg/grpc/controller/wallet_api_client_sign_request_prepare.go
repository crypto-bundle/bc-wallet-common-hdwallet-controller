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
	"google.golang.org/protobuf/types/known/anypb"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *WalletApiClientWrapper) PrepareSignRequest(ctx context.Context,
	walletUUID string,
	sessionUUID string,
	purposeUUID string,
	accountData *anypb.Any,
) (*PrepareSignRequestResponse, error) {
	req := &PrepareSignRequestReq{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletUUID,
		},
		SessionIdentifier: &WalletSessionIdentity{
			SessionUUID: sessionUUID,
		},
		AccountIdentifier: &pbCommon.AccountIdentity{
			Parameters: accountData,
			Address:    "",
		},
		SignPurposeIdentifier: &SignPurposeIdentity{
			UUID: purposeUUID,
		},
	}

	resp, err := s.originGRPCClient.PrepareSignRequest(ctx, req)
	if err != nil {
		grpcStatus, statusExists := status.FromError(err)
		if !statusExists {
			s.logger.Error("unable get status from error", zap.Error(err))

			return nil, fmt.Errorf("%w: %s", ErrUnableDecodeGrpcErrorStatus, err)
		}

		switch grpcStatus.Code() {
		case codes.NotFound, codes.ResourceExhausted:
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
