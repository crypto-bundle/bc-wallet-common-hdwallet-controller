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
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type configService interface {
	GetStageName() string
	GetApplicationName() string
	GetProviderName() string
	GetNetworkName() string

	GetInstanceIdentifier() uuid.UUID

	GetEventChannelName() string
	GetEventChannelWorkersCount() int
	GetEventChannelBufferSize() int
}

type mnemonicWalletsCacheStoreService interface {
	GetMnemonicWalletByUUID(ctx context.Context,
		MnemonicWalletUUID string,
	) (*entities.MnemonicWallet, error)
	GetMnemonicWalletSessionInfoByUUID(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) (wallet *entities.MnemonicWallet, session *entities.MnemonicWalletSession, err error)
}

type mnemonicWalletsDataService interface {
	GetMnemonicWalletByUUID(ctx context.Context, uuid string) (*entities.MnemonicWallet, error)
	GetWalletSessionByUUID(ctx context.Context,
		sessionUUID string,
	) (*entities.MnemonicWalletSession, error)
}

type signRequestDataService interface {
	GetSignRequestItemByUUIDAndStatus(ctx context.Context,
		signReqUUID string,
		status types.SignRequestStatus,
	) (*entities.SignRequest, error)
}

type publisherDriverService interface {
	SendEvent(ctx context.Context, event proto.Message) error
}

type eventProcessorService interface {
	Process(ctx context.Context, event *pbApi.Event) error
}

type walletSessionStartedProcessorService interface {
	Process(ctx context.Context, event *pbApi.WalletSessionEvent) error
}

type walletSessionClosedProcessorService interface {
	Process(ctx context.Context, event *pbApi.WalletSessionEvent) error
}

type signRequestPreparedProcessorService interface {
	Process(ctx context.Context, event *pbApi.SignRequestEvent) error
}

type transactionalStatementManager interface {
	BeginContextualTxStatement(ctx context.Context) (context.Context, error)
	CommitContextualTxStatement(ctx context.Context) error
	RollbackContextualTxStatement(ctx context.Context) error
	BeginTxWithRollbackOnError(ctx context.Context,
		callback func(txStmtCtx context.Context) error,
	) error
}
