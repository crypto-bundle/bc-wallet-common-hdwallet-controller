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
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"math/big"
)

type validatorHashCash struct {
	logger *zap.Logger

	target *big.Int

	powProofDataSvc powProofDataService
	walletDataSvc   walletDataService
	signReqDataSvc  signRequestDataService

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

func (v *validatorHashCash) ValidateByObscurityData(ctx context.Context,
	hashData []byte,
	message proto.Message,
	obscurityItemUUID uuid.UUID,
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

	rawUUID, _ := obscurityItemUUID.MarshalBinary()
	lastProfUUIDInt := big.NewInt(0).SetBytes(rawUUID)

	sum := messageInt.Add(messageInt, lastProfUUIDInt)

	hashSum := sha256.New()
	_, err = hashSum.Write(sum.Bytes())
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

func validateProfData(originHashInt *big.Int,
	messageInt *big.Int,
	lastSessionUUID uuid.UUID,
) (bool, error) {
}

func NewValidatorHashCash(logger *zap.Logger,
	walletManagerSvc walletDataService,
	powProofDataSvc powProofDataService,
	txStmtManager transactionalStatementManager,
) *validatorHashCash {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-27))

	return &validatorHashCash{
		logger: logger,

		target: target,

		powProofDataSvc: powProofDataSvc,
		walletDataSvc:   walletManagerSvc,

		txStmtManager: txStmtManager,
	}
}
