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

package grpc

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodNameEnableWallets = "EnableWallets"
)

type EnableWalletsHandler struct {
	l *zap.Logger

	walletSvc walletManagerService
}

// nolint:funlen // fixme
func (h *EnableWalletsHandler) Handle(ctx context.Context,
	req *pbApi.EnableWalletsRequest,
) (*pbApi.EnableWalletsResponse, error) {
	var err error

	validationForm := &WalletsIdentitiesForm{}
	valid, err := validationForm.LoadAndValidate(req.WalletIdentifiers)
	if err != nil {
		h.l.Error("unable load and validate request values", zap.Error(err))

		if !valid {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "something went wrong")
	}

	enabledWalletsCount, walletsIdentities, err := h.walletSvc.EnableWalletsByUUIDList(ctx, validationForm.WalletUUIDs)
	if err != nil {
		h.l.Error("unable to disable wallets", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if walletsIdentities == nil {
		return nil, status.Error(codes.NotFound, "there are no wallets available to enable")
	}

	bookmarks := make(map[string]uint32, enabledWalletsCount)
	pbWalletsData := make([]*common.MnemonicWalletData, enabledWalletsCount)
	for i := uint(0); i != enabledWalletsCount; i++ {
		item := walletsIdentities[i]
		itemUUID := item.UUID.String()

		pbWalletsData[i] = &common.MnemonicWalletData{
			WalletIdentifier: &common.MnemonicWalletIdentity{
				WalletUUID: itemUUID,
				WalletHash: item.MnemonicHash,
			},
			WalletStatus: common.WalletStatus(item.Status),
		}

		bookmarks[itemUUID] = uint32(i)
	}

	return &pbApi.EnableWalletsResponse{
		WalletsCount: uint32(enabledWalletsCount),
		WalletsData:  pbWalletsData,
		Bookmarks:    bookmarks,
	}, nil
}

func MakeEnableWalletsHandler(loggerEntry *zap.Logger,
	walletSvc walletManagerService,
) *EnableWalletsHandler {
	return &EnableWalletsHandler{
		l:         loggerEntry.With(zap.String(MethodNameTag, MethodNameEnableWallets)),
		walletSvc: walletSvc,
	}
}
