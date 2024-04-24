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

package sign_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"

	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) ExecuteSignRequest(ctx context.Context,
	signReqItem *entities.SignRequest,
	transactionData []byte,
) (signerAddr *pbCommon.DerivationAddressIdentity, signedData []byte, err error) {
	if signReqItem.DerivationPath == nil {
		return nil, nil, ErrMissingDerivationPathField
	}

	err = s.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		clbErr := s.signReqDataSvc.UpdateSignRequestItemStatus(txStmtCtx, signReqItem.UUID,
			types.SignRequestStatusSigned)
		if clbErr != nil {
			return clbErr
		}

		accountIdx, internalIdx, addrIdx := uint32(signReqItem.DerivationPath[0]),
			uint32(signReqItem.DerivationPath[1]),
			uint32(signReqItem.DerivationPath[2])

		signerAddr, signedData, clbErr = s.signTransaction(txStmtCtx, signReqItem.WalletUUID,
			accountIdx, internalIdx, addrIdx, transactionData)
		if clbErr != nil {
			return clbErr
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return
}

func (s *Service) signTransaction(ctx context.Context,
	mnemonicUUID string,
	account, change, index uint32,
	transactionData []byte,
) (signerAddr *pbCommon.DerivationAddressIdentity, signedData []byte, err error) {
	signResp, err := s.hdWalletClientSvc.SignData(ctx, &hdwallet.SignDataRequest{
		MnemonicWalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: mnemonicUUID,
		},
		AddressIdentifier: &pbCommon.DerivationAddressIdentity{
			AccountIndex:  account,
			InternalIndex: change,
			AddressIndex:  index,
		},
		DataForSign: transactionData,
	})
	if err != nil {
		grpcStatus, statusExists := status.FromError(err)
		if !statusExists {
			s.logger.Error("unable get status from error", zap.Error(ErrUnableDecodeGrpcErrorStatus))
			return nil, nil, ErrUnableDecodeGrpcErrorStatus
		}

		switch grpcStatus.Code() {
		case codes.NotFound, codes.ResourceExhausted:
			return nil, nil, nil

		default:
			s.logger.Error("unable get block by hash from bc-adapter",
				zap.Error(ErrUnableDecodeGrpcErrorStatus))

			return nil, nil, err
		}
	}

	return signResp.TxOwnerIdentity, signResp.SignedData, nil
}
