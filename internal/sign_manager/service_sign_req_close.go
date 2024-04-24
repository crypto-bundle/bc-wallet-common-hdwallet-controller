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
)

func (s *Service) CloseSignRequest(ctx context.Context,
	signReqUUID string,
) error {
	err := s.signReqDataSvc.UpdateSignRequestItemStatus(ctx, signReqUUID,
		types.SignRequestStatusPrepared)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CloseSignRequestBySession(ctx context.Context,
	sessionUUID string,
) (uint, []*entities.SignRequest, error) {
	count, signReqItems, err := s.signReqDataSvc.UpdateSignRequestItemStatusBySessionUUID(ctx, sessionUUID,
		types.SignRequestStatusClosed)
	if err != nil {
		return 0, nil, err
	}

	if signReqItems == nil {
		return 0, nil, nil
	}

	return count, signReqItems, nil
}

func (s *Service) CloseSignRequestByMultipleWallets(ctx context.Context,
	walletUUIDs []string,
) (uint, []*entities.SignRequest, error) {
	count, signReqItems, err := s.signReqDataSvc.UpdateSignRequestItemStatusByWalletsUUIDList(ctx, walletUUIDs,
		types.SignRequestStatusClosed)
	if err != nil {
		return 0, nil, err
	}

	if signReqItems == nil {
		return 0, nil, nil
	}

	return count, signReqItems, nil
}

func (s *Service) CloseSignRequestByWallet(ctx context.Context,
	walletUUID string,
) (uint, []*entities.SignRequest, error) {
	count, signReqItems, err := s.signReqDataSvc.UpdateSignRequestItemStatusByWalletUUID(ctx, walletUUID,
		types.SignRequestStatusClosed)
	if err != nil {
		return 0, nil, err
	}

	if signReqItems == nil {
		return 0, nil, nil
	}

	return count, signReqItems, nil
}
