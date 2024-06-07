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
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

func (s *pgRepository) AddNewWalletSessionAccessTokenItem(ctx context.Context,
	toSaveItem *entities.AccessTokenWalletSession,
) (result *entities.AccessTokenWalletSession, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		date := time.Now()

		row := stmt.QueryRowx(`INSERT INTO "wallet_sessions_access_tokens" ("serial_number", "token_uuid",
				 "wallet_session_uuid",
				 "created_at")
            VALUES($1, $2, $3, $4) RETURNING *;`,
			toSaveItem.SerialNumber,
			toSaveItem.AccessTokeUUID, toSaveItem.SessionUUID,
			date)

		item := &entities.AccessTokenWalletSession{}
		clbErr := row.StructScan(item)
		if clbErr != nil {
			s.logger.Error("failed to insert new access_token_wallet_session item", zap.Error(clbErr))

			return fmt.Errorf("%s: %w", ErrUnableExecuteQuery, clbErr)
		}

		result = item

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pgRepository) GetWalletSessionAccessTokenItemsByTokenUUID(ctx context.Context,
	tokenUUID string,
) (count uint, list []*entities.AccessTokenWalletSession, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, clbErr := stmt.Queryx(`SELECT *
	       FROM "wallet_sessions_access_tokens"
	       WHERE "token_uuid" = $1`, tokenUUID)
		if clbErr != nil {
			return clbErr
		}
		defer rows.Close()

		itemsList := make([]*entities.AccessTokenWalletSession, 0)
		itemsCount := uint(0)

		for rows.Next() {
			data := &entities.AccessTokenWalletSession{}

			scanErr := rows.StructScan(data)
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

	return 0, list, nil
}

func (s *pgRepository) GetLastWalletSessionNumberByAccessTokenUUID(ctx context.Context,
	accessTokenUUID string,
) (serialNumber uint64, err error) {
	if err = s.pgConn.MustWithTransaction(ctx, func(stmt *sqlx.Tx) error {
		row := stmt.QueryRowx(`SELECT coalesce(max("serial_number"), -1) 
			FROM (
				SELECT "serial_number"
				FROM "wallet_sessions_access_tokens"
				WHERE 
				    "token_uuid" = $1
				FOR UPDATE
			) AS "serial_number"`, accessTokenUUID)

		clbErr := row.Err()
		if clbErr != nil {
			return clbErr
		}

		var number int64
		clbErr = row.Scan(number)
		if clbErr != nil {
			return commonPostgres.EmptyOrError(clbErr,
				"unable get last wallet session number")
		}

		if number >= 0 {
			serialNumber = uint64(number)
		}

		serialNumber = 0

		return nil
	}); err != nil {
		return 0, err
	}

	return
}

func (s *pgRepository) GetNextWalletSessionNumberByAccessTokenUUID(ctx context.Context,
	accessTokenUUID string,
) (nextSerialNumber uint64, err error) {
	if err = s.pgConn.MustWithTransaction(ctx, func(stmt *sqlx.Tx) error {
		row := stmt.QueryRowx(`SELECT coalesce(max("serial_number")+1, -1) AS "next_number" 
			FROM (
				SELECT "serial_number"
				FROM "wallet_sessions_access_tokens"
				WHERE 
				    "token_uuid" = $1
				FOR UPDATE
			) AS "serial_number"`, accessTokenUUID)

		clbErr := row.Err()
		if clbErr != nil {
			return clbErr
		}

		var number int64
		clbErr = row.Scan(&number)
		if clbErr != nil {
			return commonPostgres.EmptyOrError(clbErr,
				"unable get last wallet session number")
		}

		if number >= 0 {
			nextSerialNumber = uint64(number)

			return nil
		}

		nextSerialNumber = 0

		return nil
	}); err != nil {
		return 0, err
	}

	return
}

func (s *pgRepository) GetLastWalletSessionIdentityByAccessTokenUUID(ctx context.Context,
	accessTokenUUID string,
) (resultItem *entities.AccessTokenWalletSession, err error) {
	if err = s.pgConn.MustWithTransaction(ctx, func(stmt *sqlx.Tx) error {
		row := stmt.QueryRowx(`SELECT * FROM "wallet_sessions_access_tokens" WHERE
			"serial_number" = (SELECT coalesce(max("serial_number"), -1)
			FROM (
				SELECT "serial_number"
				FROM "wallet_sessions_access_tokens"
				WHERE
					"token_uuid" = $1
				FOR UPDATE
			) AS "serial_number") AND
    		"token_uuid" = $1`, accessTokenUUID)

		clbErr := row.Err()
		if clbErr != nil {
			return clbErr
		}

		item := &entities.AccessTokenWalletSession{}
		clbErr = row.StructScan(item)
		if clbErr != nil {
			return commonPostgres.EmptyOrError(clbErr,
				"unable get last wallet session identity")
		}

		resultItem = item

		return nil
	}); err != nil {
		return nil, err
	}

	return
}
