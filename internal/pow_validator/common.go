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

package pow_validator

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
)

type walletDataService interface {
	GetLastWalletSessionIdentityByAccessTokenUUID(ctx context.Context,
		accessTokenUUID string,
	) (resultItem *entities.AccessTokenWalletSession, err error)
}

type signRequestDataService interface {
	GetSignRequestItemByUUIDAndStatus(ctx context.Context,
		signReqUUID string,
		status types.SignRequestStatus,
	) (*entities.SignRequest, error)
}

type powProofDataService interface {
	GetMaxNoncePowProofByEntityUUIDAndType(ctx context.Context,
		entityUUID string,
		entityType types.PowProofEntityType,
	) (lastNonce int64, err error)
	GetPowProofByEntityNonceAndType(ctx context.Context,
		entityNonce int64,
		entityType types.PowProofEntityType,
	) (powProof *entities.PowProof, err error)
	GetPowProofByEntityUUIDAndType(ctx context.Context,
		entityUUID string,
		entityType types.PowProofEntityType,
	) (powProof *entities.PowProof, err error)
	GetPowProofByUUID(ctx context.Context,
		uuid string,
	) (powProof *entities.PowProof, err error)
	AddNewPowProof(ctx context.Context,
		toSaveItem *entities.PowProof,
	) (result *entities.PowProof, err error)
}

type transactionalStatementManager interface {
	BeginContextualTxStatement(ctx context.Context) (context.Context, error)
	CommitContextualTxStatement(ctx context.Context) error
	RollbackContextualTxStatement(ctx context.Context) error
	BeginTxWithRollbackOnError(ctx context.Context,
		callback func(txStmtCtx context.Context) error,
	) error
}
