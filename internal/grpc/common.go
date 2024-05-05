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
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"google.golang.org/protobuf/types/known/anypb"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	"github.com/google/uuid"
)

type configService interface {
	IsDev() bool
	IsDebug() bool
	IsLocal() bool

	GetBindPort() string

	GetProviderName() string
	GetNetworkName() string
}

type mnemonicWalletsDataService interface {
	AddNewMnemonicWallet(ctx context.Context, wallet *entities.MnemonicWallet) (*entities.MnemonicWallet, error)
	AddNewMnemonicWallets(ctx context.Context,
		walletsList []*entities.MnemonicWallet,
	) (uint32, []*entities.MnemonicWallet, error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletUUID(ctx context.Context, walletUUID uuid.UUID) (*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDList(ctx context.Context,
		UUIDList []string,
	) ([]*entities.MnemonicWallet, error)
	GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetAllNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
}

type walletSessionDataService interface {
	AddNewWalletSession(ctx context.Context,
		sessionItem *entities.MnemonicWalletSession,
	) (*entities.MnemonicWalletSession, error)
	GetAllActiveSessions(ctx context.Context,
	) (uint, []*entities.MnemonicWalletSession, error)
	GetAllActiveSessionsClb(ctx context.Context,
		onItemCallBack func(item *entities.MnemonicWalletSession) error,
	) (uint, []*entities.MnemonicWalletSession, error)
	GetWalletActiveSessions(ctx context.Context,
		walletUUID string,
	) (uint, []*entities.MnemonicWalletSession, error)
	GetWalletActiveSessionsClb(ctx context.Context,
		walletUUID string,
		onItemCallBack func(item *entities.MnemonicWalletSession) error,
	) (uint, []*entities.MnemonicWalletSession, error)
	GetWalletSessionByUUID(ctx context.Context,
		sessionUUID string,
	) (*entities.MnemonicWalletSession, error)
}

type walletManagerService interface {
	AddNewWallet(ctx context.Context) (*entities.MnemonicWallet, error)
	ImportWallet(ctx context.Context, mnemonicData []byte) (*entities.MnemonicWallet, error)
	EnableWalletByUUID(ctx context.Context,
		walletUUID string,
	) (*entities.MnemonicWallet, error)
	DisableWalletByUUID(ctx context.Context,
		walletUUID string,
	) (*entities.MnemonicWallet, error)
	DisableWalletsByUUIDList(ctx context.Context,
		walletUUIDs []string,
	) (count uint, list []*entities.MnemonicWallet, err error)
	EnableWalletsByUUIDList(ctx context.Context,
		walletUUIDs []string,
	) (count uint, list []*entities.MnemonicWallet, err error)
	GetEnabledWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetWalletByUUID(ctx context.Context, walletUUID string) (*entities.MnemonicWallet, error)
	GetAccount(ctx context.Context,
		mnemonicUUID string,
		accountParameters *anypb.Any,
	) (address *string, err error)
	GetAccounts(ctx context.Context,
		mnemonicUUID string,
		accountsParameters *anypb.Any,
	) (count uint64, list []*pbCommon.AccountIdentity, err error)
	StartWalletSession(ctx context.Context,
		walletUUID string,
	) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error)
	StartSessionForWallet(ctx context.Context,
		wallet *entities.MnemonicWallet,
	) (*entities.MnemonicWalletSession, error)
	CloseWalletSession(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error)
	GetWalletSessionInfo(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) (*entities.MnemonicWallet, *entities.MnemonicWalletSession, error)
	GetWalletSessionsByWalletUUID(ctx context.Context,
		walletUUID string,
	) (wallet *entities.MnemonicWallet, list []*entities.MnemonicWalletSession, err error)
}

type signManagerService interface {
	GetActiveSignRequest(ctx context.Context,
		signUUID string,
	) (signReqItem *entities.SignRequest, err error)
	PrepareSignRequest(ctx context.Context,
		mnemonicUUID string,
		sessionUUID string,
		purposeUUID string,
		accountParameters *anypb.Any,
	) (addr *pbCommon.AccountIdentity, signReqItem *entities.SignRequest, err error)
	ExecuteSignRequest(ctx context.Context,
		signReqItem *entities.SignRequest,
		transactionData []byte,
	) (signerAddr *pbCommon.AccountIdentity, signedData []byte, err error)
	CloseSignRequest(ctx context.Context,
		signReqUUID string,
	) error
	CloseSignRequestBySession(ctx context.Context,
		sessionUUID string,
	) (count uint, list []*entities.SignRequest, err error)
	CloseSignRequestByWallet(ctx context.Context,
		walletUUID string,
	) (count uint, list []*entities.SignRequest, err error)
	CloseSignRequestByMultipleWallets(ctx context.Context,
		walletUUIDs []string,
	) (count uint, list []*entities.SignRequest, err error)
}

type marshallerService interface {
	MarshallCreateWalletData(wallet *entities.MnemonicWallet) *pbApi.AddNewWalletResponse
	MarshallGetEnabledWallets([]*entities.MnemonicWallet) *pbApi.GetEnabledWalletsResponse
	MarshallWalletSessions(
		sessionsList []*entities.MnemonicWalletSession,
	) []*pbApi.SessionInfo
}

type addWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.AddNewWalletRequest) (*pbApi.AddNewWalletResponse, error)
}

type importWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.ImportWalletRequest) (*pbApi.ImportWalletResponse, error)
}

type disableWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.DisableWalletRequest) (*pbApi.DisableWalletResponse, error)
}

type disableWalletsHandlerService interface {
	Handle(ctx context.Context, request *pbApi.DisableWalletsRequest) (*pbApi.DisableWalletsResponse, error)
}

type enableWalletsHandlerService interface {
	Handle(ctx context.Context, request *pbApi.EnableWalletsRequest) (*pbApi.EnableWalletsResponse, error)
}

type enableWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.EnableWalletRequest) (*pbApi.EnableWalletResponse, error)
}

type getWalletHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetWalletInfoRequest) (*pbApi.GetWalletInfoResponse, error)
}

type getAccountHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetAccountRequest) (*pbApi.GetAccountResponse, error)
}

type getAccountsHandlerService interface {
	Handle(ctx context.Context,
		request *pbApi.GetMultipleAccountRequest,
	) (*pbApi.GetMultipleAccountResponse, error)
}

type prepareSignRequestHandlerService interface {
	Handle(ctx context.Context, request *pbApi.PrepareSignRequestReq) (*pbApi.PrepareSignRequestResponse, error)
}

type executeSignRequestHandlerService interface {
	Handle(ctx context.Context, request *pbApi.ExecuteSignRequestReq) (*pbApi.ExecuteSignRequestResponse, error)
}

type getEnabledWalletsHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetEnabledWalletsRequest) (*pbApi.GetEnabledWalletsResponse, error)
}

type startWalletSessionHandlerService interface {
	Handle(ctx context.Context, request *pbApi.StartWalletSessionRequest) (*pbApi.StartWalletSessionResponse, error)
}

type closeWalletSessionHandlerService interface {
	Handle(ctx context.Context, request *pbApi.CloseWalletSessionsRequest) (*pbApi.CloseWalletSessionsResponse, error)
}

type getWalletSessionHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetWalletSessionRequest) (*pbApi.GetWalletSessionResponse, error)
}

type getWalletSessionsHandlerService interface {
	Handle(ctx context.Context, request *pbApi.GetWalletSessionsRequest) (*pbApi.GetWalletSessionsResponse, error)
}
