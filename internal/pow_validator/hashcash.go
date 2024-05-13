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
	"crypto/sha256"
	"math/big"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type validatorHashCash struct {
	logger *zap.Logger

	target *big.Int

	powProofDataSvc  powProofDataService
	walletManagerSvc walletManagerService
	signReqDataSvc   signRequestDataService

	txStmtManager transactionalStatementManager
}

func (v *validatorHashCash) PreValidate(ctx context.Context,
	hashData []byte,
) bool {
	margin := len(hashData) - 8

	hash := hashData[:margin]

	hashInt := big.NewInt(0).SetBytes(hash)
	cmp := v.target.Cmp(hashInt)
	if cmp < 1 {
		return false
	}

	return true
}

func (v *validatorHashCash) Validate(ctx context.Context,
	hashData []byte,
	message proto.Message,
	accessTokenUUID uuid.UUID,
	clbFunc func(ctx context.Context, req any) (any, error),
) (bool, error) {
	margin := len(hashData) - 8

	nonce := hashData[margin:]
	hash := hashData[:margin]

	hashInt := big.NewInt(0).SetBytes(hash)
	cmp := v.target.Cmp(hashInt)
	if cmp < 1 {
		return false, nil
	}

	messageRaw, err := proto.Marshal(message)
	if err != nil {
		return false, err
	}

	messageInt := big.NewInt(0).SetBytes(messageRaw)
	nonceNumber := big.NewInt(0).SetBytes(nonce).Int64()

	return v.validateForSessionFlow(ctx, hashInt,
		messageInt, nonceNumber, accessTokenUUID, clbFunc)
}

func (v *validatorHashCash) validateForSessionFlow(ctx context.Context,
	originHashInt *big.Int,
	messageInt *big.Int,
	nonceNumber int64,
	entityUUID uuid.UUID,
	clbFunc func(ctx context.Context, req any) (any, error),
) (isValid bool, err error) {
	var lastProof *entities.PowProof
	err = v.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		lastNonceNumber, clbErr := v.powProofDataSvc.GetMaxNoncePowProofByEntityUUIDAndType(txStmtCtx,
			entityUUID.String(), types.PowProofEntityTypeSession)
		if clbErr != nil {
			return clbErr
		}

		if lastNonceNumber != nonceNumber+1 {
			isValid = false

			return nil
		}

		lastProofItem, clbErr := v.powProofDataSvc.GetPowProofByEntityNonceAndType(txStmtCtx,
			lastNonceNumber, types.PowProofEntityTypeSession)
		if clbErr != nil {
			return clbErr
		}

		lastProof = lastProofItem

		var lastProofUUID uuid.UUID
		if lastProof != nil {
			parsedUUID, parseErr := uuid.Parse(lastProof.UUID)
			if parseErr != nil {
				return parseErr
			}

			lastProofUUID = parsedUUID
		}

		isProofDataValid, clbErr := validateProfData(originHashInt, messageInt, lastProofUUID)
		if clbErr != nil {
			return clbErr
		}

		if !isProofDataValid {
			isValid = false

			return nil
		}

		_, clbErr = v.powProofDataSvc.AddNewPowProof(txStmtCtx, &entities.PowProof{
			UUID:        uuid.NewString(),
			EntityUUID:  entityUUID.String(),
			EntityType:  types.PowProofEntityTypeSession,
			EntityNonce: nonceNumber + 1,
			HashData:    originHashInt.Bytes(),
		})
		if clbErr != nil {
			return clbErr
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func validateProfData(originHashInt *big.Int,
	messageInt *big.Int,
	lastProfUUID uuid.UUID,
) (bool, error) {

	rawUUID, _ := lastProfUUID.MarshalBinary()
	lastProfUUIDInt := big.NewInt(0).SetBytes(rawUUID)

	sum := messageInt.Add(messageInt, lastProfUUIDInt)

	hashSum := sha256.New()
	_, err := hashSum.Write(sum.Bytes())
	if err != nil {
		return false, err
	}

	resultHash := hashSum.Sum(nil)
	resultHashInt := big.NewInt(0).SetBytes(resultHash)

	cmpRes := originHashInt.Cmp(resultHashInt)
	if cmpRes != 0 {
		return false, nil
	}

	return true, nil
}

func NewValidatorHashCash(logger *zap.Logger,
	walletManagerSvc walletManagerService,
	signReqDataSvc signRequestDataService,
	powProofDataSvc powProofDataService,
	txStmtManager transactionalStatementManager,
) *validatorHashCash {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-27))

	return &validatorHashCash{
		logger: logger,

		target: target,

		powProofDataSvc:  powProofDataSvc,
		walletManagerSvc: walletManagerSvc,
		signReqDataSvc:   signReqDataSvc,

		txStmtManager: txStmtManager,
	}
}
