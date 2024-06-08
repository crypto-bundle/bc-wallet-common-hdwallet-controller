/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package pg_store

import (
	"context"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

var nopSessionCallback = func(item *entities.MnemonicWalletSession) error {
	return nil
}

func (s *pgRepository) AddNewWalletSession(ctx context.Context,
	sessionItem *entities.MnemonicWalletSession,
) (*entities.MnemonicWalletSession, error) {
	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		var sessionID uint32
		row := stmt.QueryRowx(`INSERT INTO "mnemonic_wallet_sessions" ("uuid", 
				"mnemonic_wallet_uuid",
				"status",
			  	"started_at", "expired_at",           
       			"created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id;`,
			sessionItem.UUID,
			sessionItem.MnemonicWalletUUID,
			sessionItem.Status,
			sessionItem.StartedAt, sessionItem.ExpiredAt,
			sessionItem.CreatedAt, sessionItem.UpdatedAt)

		err := row.Scan(&sessionID)
		if err != nil {
			s.logger.Error("failed to insert in wallets table", zap.Error(err))

			return ErrUnableExecuteQuery
		}

		sessionItem.ID = sessionID

		return nil
	}); err != nil {
		return nil, err
	}

	return sessionItem, nil
}

func (s *pgRepository) UpdateWalletSessionStatusByWalletUUID(ctx context.Context,
	walletUUID string,
	newStatus types.MnemonicWalletSessionStatus,
) error {

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		_, callbackErr := stmt.Exec(`UPDATE "mnemonic_wallet_sessions" 
			SET "status" = $1,
				"updated_at" = $2
			WHERE "mnemonic_wallet_uuid" = $3`,
			newStatus, time.Now(), walletUUID)
		if callbackErr != nil {
			return callbackErr
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *pgRepository) UpdateWalletSessionStatusBySessionUUID(ctx context.Context,
	sessionUUID string,
	newStatus types.MnemonicWalletSessionStatus,
) (result *entities.MnemonicWalletSession, err error) {

	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`UPDATE "mnemonic_wallet_sessions" 
			SET "status" = $1,
				"updated_at" = $2
			WHERE "uuid" = $3
			RETURNING *`,
			newStatus, time.Now(), sessionUUID)

		sessionItem := &entities.MnemonicWalletSession{}
		clbErr := row.StructScan(sessionItem)
		if clbErr != nil {
			s.logger.Error("unable to update wallet session status", zap.Error(clbErr))

			return fmt.Errorf("%s: %w", ErrUnableExecuteQuery.Error(), clbErr)
		}

		result = sessionItem

		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func (s *pgRepository) UpdateMultipleWalletSessionStatus(ctx context.Context,
	walletsUUIDs []string,
	newStatus types.MnemonicWalletSessionStatus,
) (count uint, list []string, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, clbErr := sqlx.In(`UPDATE "mnemonic_wallet_sessions"
	       SET "status" = ?
	       WHERE "mnemonic_wallet_uuid" IN (?)
	       RETURNING "uuid"`, newStatus, walletsUUIDs)

		bonded := stmt.Rebind(query)
		returnedRows, clbErr := stmt.Queryx(bonded, args...)
		if clbErr != nil {
			return clbErr
		}
		defer returnedRows.Close()

		sessionsUUIDList := make([]string, 0)

		for returnedRows.Next() {
			var sessionUUID string

			scanErr := returnedRows.Scan(&sessionUUID)
			if scanErr != nil {
				return scanErr
			}

			sessionsUUIDList = append(sessionsUUIDList, sessionUUID)

			count++
		}

		list = sessionsUUIDList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) UpdateMultipleWalletSessionStatusClb(ctx context.Context,
	walletsUUIDs []string,
	newStatus types.MnemonicWalletSessionStatus,
	oldStatus []types.MnemonicWalletSessionStatus,
	clbFunc func(*entities.MnemonicWalletSession) error,
) (count uint, list []*entities.MnemonicWalletSession, err error) {
	err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, clbErr := sqlx.In(`UPDATE "mnemonic_wallet_sessions"
         	SET "status" = ?
			WHERE "mnemonic_wallet_uuid" IN (?) AND
				? BETWEEN "started_at" AND "expired_at" 
			AND
				"status" IN (?)
			RETURNING *`, newStatus, walletsUUIDs, time.Now(), oldStatus)

		bonded := stmt.Rebind(query)
		returnedRows, clbErr := stmt.Queryx(bonded, args...)
		if clbErr != nil {
			return clbErr
		}
		defer returnedRows.Close()

		sessionsList := make([]*entities.MnemonicWalletSession, 0)

		for returnedRows.Next() {
			session := &entities.MnemonicWalletSession{}

			loopErr := returnedRows.StructScan(session)
			if loopErr != nil {
				return loopErr
			}

			loopErr = clbFunc(session)
			if loopErr != nil {
				return loopErr
			}

			sessionsList = append(sessionsList, session)

			count++
		}

		list = sessionsList

		return nil
	})
	if err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) GetWalletSessionByUUID(ctx context.Context,
	sessionUUID string,
) (*entities.MnemonicWalletSession, error) {
	var result *entities.MnemonicWalletSession = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, clbErr := sqlx.In(`SELECT *
			FROM "mnemonic_wallet_sessions"
			WHERE "uuid" = ? AND
					? BETWEEN "started_at" AND "expired_at" 
				AND
			    	"status" IN (?)
			ORDER BY "expired_at";`, sessionUUID, time.Now(), []types.MnemonicWalletSessionStatus{
			types.MnemonicWalletSessionStatusPrepared,
		})
		if clbErr != nil {
			return clbErr
		}
		bonded := stmt.Rebind(query)
		row := stmt.QueryRowx(bonded, args...)
		if clbErr != nil {
			return clbErr
		}

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		session := &entities.MnemonicWalletSession{}
		callbackErr = row.StructScan(session)
		if callbackErr != nil {
			return commonPostgres.EmptyOrError(callbackErr, "unable get session by UUID")
		}

		result = session

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pgRepository) GetActiveWalletSessionsByWalletUUID(ctx context.Context, walletUUID string) (
	count uint, list []*entities.MnemonicWalletSession, err error,
) {
	return s.GetWalletSessionsByWalletUUIDAndStatusClb(ctx, walletUUID,
		[]types.MnemonicWalletSessionStatus{
			types.MnemonicWalletSessionStatusPrepared,
		},
		nopSessionCallback)
}

func (s *pgRepository) GetWalletSessionsByWalletUUIDAndStatusClb(ctx context.Context,
	walletUUID string,
	sessionStatuses []types.MnemonicWalletSessionStatus,
	onItemCallBack func(item *entities.MnemonicWalletSession) error,
) (count uint, list []*entities.MnemonicWalletSession, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, clbErr := sqlx.In(`SELECT *
			FROM "mnemonic_wallet_sessions"
			WHERE "mnemonic_wallet_uuid" = ? AND
					? BETWEEN "started_at" AND "expired_at" 
				AND
			      	"status" IN (?)
			ORDER BY "expired_at";`, walletUUID, time.Now(), sessionStatuses)
		if clbErr != nil {
			return clbErr
		}

		bonded := stmt.Rebind(query)
		rows, clbErr := stmt.Queryx(bonded, args...)
		if clbErr != nil {
			return clbErr
		}
		defer rows.Close()

		itemsList := make([]*entities.MnemonicWalletSession, 0)
		itemsCount := uint(0)

		for rows.Next() {
			data := &entities.MnemonicWalletSession{}

			scanErr := rows.StructScan(data)
			if scanErr != nil {
				return scanErr
			}

			scanErr = onItemCallBack(data)
			if scanErr != nil {
				return scanErr
			}

			itemsList = append(itemsList, data)
			itemsCount++
		}

		if itemsCount == 0 {
			itemsList = nil
		}

		count = itemsCount
		list = itemsList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) GetActiveWalletSessions(ctx context.Context) (
	count uint, list []*entities.MnemonicWalletSession, err error,
) {
	return s.GetActiveWalletSessionsClb(ctx, nopSessionCallback)
}

func (s *pgRepository) GetActiveWalletSessionsClb(ctx context.Context,
	onItemCallBack func(item *entities.MnemonicWalletSession) error,
) (count uint, list []*entities.MnemonicWalletSession, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, queryErr := stmt.Queryx(`SELECT *
			FROM "mnemonic_wallet_sessions"
			WHERE $1 BETWEEN "started_at" AND "expired_at" 
				AND "status" = $2
			ORDER BY "expired_at";`, time.Now(), types.MnemonicWalletSessionStatusPrepared)
		if queryErr != nil {
			return queryErr
		}

		defer rows.Close()

		itemsList := make([]*entities.MnemonicWalletSession, 0)
		itemsCount := uint(0)

		for rows.Next() {
			data := &entities.MnemonicWalletSession{}

			scanErr := rows.StructScan(data)
			if scanErr != nil {
				return scanErr
			}

			scanErr = onItemCallBack(data)
			if scanErr != nil {
				return scanErr
			}

			itemsList = append(itemsList, data)
			itemsCount++
		}

		if itemsCount == 0 {
			itemsList = nil
		}

		count = itemsCount
		list = itemsList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}
