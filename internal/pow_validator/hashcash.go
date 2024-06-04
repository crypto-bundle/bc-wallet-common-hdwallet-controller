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
	"go.uber.org/zap"
	"math/big"
)

type validatorHashCash struct {
	logger *zap.Logger

	target *big.Int
}

func (v *validatorHashCash) PreValidate(ctx context.Context,
	hashData []byte,
) bool {
	hashInt := big.NewInt(0).SetBytes(hashData)
	cmp := v.target.Cmp(hashInt)
	if cmp != 1 {
		return false
	}

	return true
}

func (v *validatorHashCash) ValidateByObscurityData(ctx context.Context,
	hashData []byte,
	nonce int64,
	message []byte,
	obscurityData []byte,
) (valid bool, err error) {
	originHashInt := big.NewInt(0).SetBytes(hashData)
	cmp := v.target.Cmp(originHashInt)
	if cmp != 1 {
		return false, nil
	}

	message = append(message, obscurityData...)
	message = append(message, byte(nonce))

	hashSum := sha256.New()
	_, err = hashSum.Write(message)
	if err != nil {
		return false, err
	}

	resultHash := hashSum.Sum(nil)
	resultHashInt := big.NewInt(0).SetBytes(resultHash)

	cmp = originHashInt.Cmp(resultHashInt)
	if cmp != 0 {
		return false, nil
	}

	cmp = v.target.Cmp(resultHashInt)
	if cmp != 1 {
		return false, nil
	}

	return true, nil
}

func NewValidatorHashCash(logger *zap.Logger) *validatorHashCash {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-4))

	return &validatorHashCash{
		logger: logger,

		target: target,
	}
}
