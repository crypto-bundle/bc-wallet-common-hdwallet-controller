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

package wallet_manager

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"github.com/google/uuid"
	"time"
)

type accessTokenDataService interface {
	GetAccessTokenInfoByUUID(ctx context.Context,
		tokenUUID string,
	) (*entities.AccessToken, error)
	AddNewAccessToken(ctx context.Context,
		toSaveItem *entities.AccessToken,
	) (result *entities.AccessToken, err error)
	AddMultipleAccessTokens(ctx context.Context,
		bcTxItems []*entities.AccessToken,
	) (count uint, list []*entities.AccessToken, err error)
}

type accessTokenListIterator types.AccessTokenListIterator

type jwtService interface {
	ExtractFields(tokenData []byte) (*uuid.UUID, *time.Time, error)
}

type tokenDataAdapter interface {
	Adopt(walletUUID uuid.UUID,
		iterator accessTokenListIterator,
	) ([]*entities.AccessToken, error)
}

type configService interface {
	GetDefaultWalletSessionDelay() time.Duration
	GetDefaultWalletUnloadInterval() time.Duration
}

type mnemonicWalletsCacheStoreService interface {
	SetMnemonicWalletItem(ctx context.Context,
		walletItem *entities.MnemonicWallet,
	) error
	SetMultipleMnemonicWallets(ctx context.Context,
		walletItems []*entities.MnemonicWallet,
	) error
	GetAllWallets(ctx context.Context) ([]*entities.MnemonicWallet, error)
	GetMnemonicWalletByUUID(ctx context.Context,
		MnemonicWalletUUID string,
	) (*entities.MnemonicWallet, error)
	GetMnemonicWalletInfoByUUID(ctx context.Context,
		MnemonicWalletUUID string,
	) (wallet *entities.MnemonicWallet, sessions []*entities.MnemonicWalletSession, err error)
	GetMnemonicWalletSessionInfoByUUID(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) (wallet *entities.MnemonicWallet, session *entities.MnemonicWalletSession, err error)
	SetMnemonicWalletSessionItem(ctx context.Context,
		sessionItem *entities.MnemonicWalletSession,
	) error
	FullUnsetMnemonicWallet(ctx context.Context,
		mnemonicWalletUUID string,
	) error
	UnsetWalletSession(ctx context.Context,
		mnemonicWalletsUUID string,
		sessionsUUID string,
	) error
	UnsetMultipleWallets(ctx context.Context,
		walletIdentities []string,
	) error
	UnsetMultipleSessions(ctx context.Context,
		sessionByWalletIdentities map[string][]string,
	) error
}

type mnemonicWalletsDataService interface {
	AddNewMnemonicWallet(ctx context.Context,
		wallet *entities.MnemonicWallet,
	) (*entities.MnemonicWallet, error)
	UpdateWalletStatus(ctx context.Context,
		walletUUID string,
		newStatus types.MnemonicWalletStatus,
	) (*entities.MnemonicWallet, error)
	UpdateMultipleWalletsStatus(ctx context.Context,
		walletUUID []string,
		newStatus types.MnemonicWalletStatus,
	) (count uint, list []*entities.MnemonicWallet, err error)
	UpdateMultipleWalletsStatusClb(ctx context.Context,
		walletUUIDs []string,
		newStatus types.MnemonicWalletStatus,
		clbFunc func(idx uint, wallet *entities.MnemonicWallet) error,
	) (count uint, list []*entities.MnemonicWallet, err error)
	UpdateMultipleWalletsStatusRetWallets(ctx context.Context,
		walletUUIDs []string,
		newStatus types.MnemonicWalletStatus,
	) (count uint, list []*entities.MnemonicWallet, err error)
	GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletByUUID(ctx context.Context, uuid string) (*entities.MnemonicWallet, error)
	GetMnemonicWalletsByStatus(ctx context.Context,
		status types.MnemonicWalletStatus,
	) ([]*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDList(ctx context.Context,
		UUIDList []string,
	) ([]*entities.MnemonicWallet, error)
	GetMnemonicWalletsByUUIDListAndStatus(ctx context.Context,
		UUIDList []string,
		status []types.MnemonicWalletStatus,
	) ([]string, []*entities.MnemonicWallet, error)

	AddNewWalletSession(ctx context.Context,
		sessionItem *entities.MnemonicWalletSession,
	) (*entities.MnemonicWalletSession, error)
	UpdateWalletSessionStatusByWalletUUID(ctx context.Context,
		walletUUID string,
		status types.MnemonicWalletSessionStatus,
	) error
	UpdateWalletSessionStatusBySessionUUID(ctx context.Context,
		sessionUUID string,
		newStatus types.MnemonicWalletSessionStatus,
	) (result *entities.MnemonicWalletSession, err error)
	UpdateMultipleWalletSessionStatus(ctx context.Context,
		sessionsUUIDs []string,
		newStatus types.MnemonicWalletSessionStatus,
	) (count uint, sessions []string, err error)
	UpdateMultipleWalletSessionStatusClb(ctx context.Context,
		walletsUUIDs []string,
		newStatus types.MnemonicWalletSessionStatus,
		oldStatus []types.MnemonicWalletSessionStatus,
		clbFunc func(*entities.MnemonicWalletSession) error,
	) (count uint, list []*entities.MnemonicWalletSession, err error)
	GetWalletSessionByUUID(ctx context.Context,
		sessionUUID string,
	) (*entities.MnemonicWalletSession, error)
	GetActiveWalletSessionsByWalletUUID(ctx context.Context, walletUUID string) (
		count uint, list []*entities.MnemonicWalletSession, err error,
	)

	AddNewWalletSessionAccessTokenItem(ctx context.Context,
		toSaveItem *entities.AccessTokenWalletSession,
	) (result *entities.AccessTokenWalletSession, err error)
	GetWalletSessionAccessTokenItemsByTokenUUID(ctx context.Context,
		tokenUUID string,
	) (count uint, list []*entities.AccessTokenWalletSession, err error)
	GetLastWalletSessionNumberByAccessTokenUUID(ctx context.Context,
		accessTokenUUID string,
	) (serialNumber uint64, err error)
	GetNextWalletSessionNumberByAccessTokenUUID(ctx context.Context,
		accessTokenUUID string,
	) (nextSerialNumber uint64, err error)
}

type signRequestDataService interface {
	AddSignRequestItem(ctx context.Context,
		toSaveItem *entities.SignRequest,
	) (*entities.SignRequest, error)
	UpdateSignRequestItemStatus(ctx context.Context,
		signReqUUID string,
		newStatus types.SignRequestStatus,
	) error
}

type encryptService interface {
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(encMsg []byte) ([]byte, error)
}

type eventPublisherService interface {
	SendSessionStartEvent(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) error
	SendSessionClosedEvent(ctx context.Context,
		walletUUID string,
		sessionUUID string,
	) error
}

type transactionalStatementManager interface {
	BeginContextualTxStatement(ctx context.Context) (context.Context, error)
	CommitContextualTxStatement(ctx context.Context) error
	RollbackContextualTxStatement(ctx context.Context) error
	BeginTxWithRollbackOnError(ctx context.Context,
		callback func(txStmtCtx context.Context) error,
	) error
}
