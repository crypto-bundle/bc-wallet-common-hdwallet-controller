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

package events

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	pbHdwallet "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
)

type signRequestPreparedHandler struct {
	walletCacheDataSvc mnemonicWalletsCacheStoreService
	walletDataSvc      mnemonicWalletsDataService

	signReqDataSvc signRequestDataService

	txStmtManager transactionalStatementManager

	hdWalletSvc pbHdwallet.HdWalletApiClient
}

func (h *signRequestPreparedHandler) Process(ctx context.Context, event *pbApi.SignRequestEvent) error {
	var sessionItem *entities.MnemonicWalletSession
	var signReqItem *entities.SignRequest

	if err := h.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		signReq, clbErr := h.signReqDataSvc.GetSignRequestItemByUUIDAndStatus(ctx, event.SignRequestIdentifier.UUID,
			types.SignRequestStatusPrepared)
		if clbErr != nil {
			return clbErr
		}

		session, clbErr := h.walletDataSvc.GetWalletSessionByUUID(ctx, signReq.SessionUUID)
		if clbErr != nil {
			return clbErr
		}

		signReqItem = signReq
		sessionItem = session

		return nil
	}); err != nil {
		return err
	}

	if sessionItem == nil || signReqItem == nil {
		return nil
	}

	return h.process(ctx, sessionItem, signReqItem)
}

func (h *signRequestPreparedHandler) process(ctx context.Context,
	sessionItem *entities.MnemonicWalletSession,
	signReqItem *entities.SignRequest,
) error {
	if !sessionItem.IsSessionActive() {
		return nil
	}

	accountIdx, internalIdx, addrIdx := uint32(signReqItem.DerivationPath[0]),
		uint32(signReqItem.DerivationPath[1]),
		uint32(signReqItem.DerivationPath[2])

	_, err := h.hdWalletSvc.LoadDerivationAddress(ctx, &pbHdwallet.LoadDerivationAddressRequest{
		MnemonicWalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: sessionItem.MnemonicWalletUUID,
		},
		AddressIdentifier: &pbCommon.DerivationAddressIdentity{
			AccountIndex:  accountIdx,
			InternalIndex: internalIdx,
			AddressIndex:  addrIdx,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func MakeEventSignReqPreparedHandler(walletCacheDataSvc mnemonicWalletsCacheStoreService,
	walletDataSvc mnemonicWalletsDataService,
	signReqDataSvc signRequestDataService,
	hdWalletSvc pbHdwallet.HdWalletApiClient,
	txStmtManager transactionalStatementManager,
) *signRequestPreparedHandler {
	return &signRequestPreparedHandler{
		walletCacheDataSvc: walletCacheDataSvc,
		walletDataSvc:      walletDataSvc,
		signReqDataSvc:     signReqDataSvc,
		txStmtManager:      txStmtManager,
		hdWalletSvc:        hdWalletSvc,
	}
}
