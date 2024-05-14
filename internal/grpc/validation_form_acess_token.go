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

package grpc

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrMissingTokenUUIDIdentity = errors.New("missing uuid in token data")
	ErrMismatchedUUIDIdentity   = errors.New("token identity mismatched")
)

const (
	TokenUUIDLabel    = "token_uuid"
	TokenExpiredLabel = "token_expired_at"
)

type accessTokenValidationForm struct {
	JWTSvc jwtService
}

func (f *accessTokenValidationForm) Validate(tokenData []byte) (*uuid.UUID, *time.Time, error) {
	data, err := f.JWTSvc.GetTokenData(string(tokenData))
	if err != nil {
		return nil, nil, err
	}

	tokenUUIDStr, isExist := data[TokenUUIDLabel]
	if !isExist {
		return nil, nil, ErrMissingTokenUUIDIdentity
	}

	tokenUUIDRaw, err := uuid.Parse(tokenUUIDStr)
	if err != nil {
		return nil, nil, err
	}

	tokenExpiredAtStr, isExist := data[TokenExpiredLabel]
	if !isExist {
		return nil, nil, ErrMissingTokenUUIDIdentity
	}

	expiredAt, err := time.Parse(time.Layout, tokenExpiredAtStr)
	if err != nil {
		return nil, nil, err
	}

	return &tokenUUIDRaw, &expiredAt, nil
}
