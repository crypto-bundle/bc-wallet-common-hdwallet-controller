/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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

	anyRawData := &anypb.Any{}
	err := proto.Unmarshal(signReqItem.AccountData, anyRawData)
	if err != nil {
		return err
	}

	_, err = h.hdWalletSvc.LoadAccount(ctx, &pbHdwallet.LoadAccountRequest{
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: sessionItem.MnemonicWalletUUID,
		},
		AccountIdentifier: &pbCommon.AccountIdentity{
			Parameters: anyRawData,
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
