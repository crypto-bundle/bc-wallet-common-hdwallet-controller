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

package controller

import (
	"context"
)

const (
	serverAddressTemplate = "BC_WALLET_%s_HDWALLET_SERVICE_HOST"
	serverPortTemplate    = "BC_WALLET_%s_HDWALLET_SERVICE_PORT"
)

type baseConfigService interface {
	GetProviderName() string
	GetNetworkName() string
}

type hdWalletClientConfigService interface {
	baseConfigService

	GetMaxReceiveMessageSize() int
	GetMaxSendMessageSize() int
	IsPowShieldEnabled() bool
	IsAccessTokenShieldEnabled() bool

	GetServerPort() uint
	GetServerHost() string
	GetServerBindAddress() string
}

type obscurityDataProvider interface {
	GetLastObscurityData(ctx context.Context,
		walletUUID string,
		accessTokenHash string) ([]byte, error)
	AddLastObscurityData(ctx context.Context,
		walletUUID string,
		accessTokenHash string,
		obscurityData []byte,
	) error
}

type accessTokensDataService interface {
	GetAccessTokenForWallet(ctx context.Context, walletUUID string) (*string, error)
}

type transactionalStatementManager interface {
	BeginTxWithRollbackOnError(ctx context.Context,
		callback func(txStmtCtx context.Context) error,
	) error
}

type jwtService interface {
}
