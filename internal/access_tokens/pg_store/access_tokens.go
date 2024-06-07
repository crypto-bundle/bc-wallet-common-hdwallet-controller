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
	"errors"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/jmoiron/sqlx"
	"time"

	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

	"go.uber.org/zap"
)

var (
	ErrUnablePrepareQuery    = errors.New("unable to prepare query")
	ErrUnableExecuteQuery    = errors.New("unable to execute query")
	ErrUnableGetLastInsertID = errors.New("unable get last insert id")
)

type pgRepository struct {
	pgConn *commonPostgres.Connection
	logger *zap.Logger
}

func (s *pgRepository) AddNewAccessToken(ctx context.Context,
	toSaveItem *entities.AccessToken,
) (result *entities.AccessToken, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		date := time.Now()

		row := stmt.QueryRowx(`INSERT INTO "access_tokens" ("uuid", "wallet_uuid",
			 	"role", "raw_data", "hash",
        		"expired_at", "created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING *;`,
			toSaveItem.UUID, toSaveItem.WalletUUID,
			toSaveItem.RawData, toSaveItem.Hash,
			toSaveItem.ExpiredAt,
			date, date)

		item := &entities.AccessToken{}
		clbErr := row.StructScan(item)
		if clbErr != nil {
			s.logger.Error("failed to insert new access_token item", zap.Error(clbErr))

			return fmt.Errorf("%s: %w", ErrUnableExecuteQuery, clbErr)
		}

		result = item

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pgRepository) AddMultipleAccessTokens(ctx context.Context,
	bcTxItems []*entities.AccessToken,
) (count uint, list []*entities.AccessToken, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		bondedSql, args, queryErr := stmt.BindNamed(`
			WITH "inserted" AS (
				INSERT INTO "access_tokens" ("uuid", "wallet_uuid", "role", "raw_data", "hash",
        		"expired_at", "created_at", "updated_at")
				VALUES(:uuid, :wallet_uuid, :role, :raw_data, :hash, :expired_at, :created_at, :updated_at) 
				ON CONFLICT DO NOTHING
				RETURNING *
			)
			SELECT * FROM "inserted"
			ORDER BY "id";`, bcTxItems)
		if queryErr != nil {
			return queryErr
		}

		rows, queryErr := stmt.Queryx(bondedSql, args...)
		if queryErr != nil {
			return queryErr
		}

		itemsList := make([]*entities.AccessToken, len(bcTxItems))
		itemsCount := uint(0)

		defer func() {
			rows.Close()
			rows = nil
			itemsList = nil
			args = nil
		}()

		for rows.Next() {
			itemData := &entities.AccessToken{}

			iterErr := rows.StructScan(itemData)
			if iterErr != nil {
				return iterErr
			}

			itemsList[itemsCount] = itemData
			itemsCount++
		}

		if itemsCount == 0 {
			return nil // returning count = 0, list = nil
		}

		count = itemsCount
		list = itemsList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) GetAccessTokenInfoByUUID(ctx context.Context,
	tokenUUID string,
) (*entities.AccessToken, error) {
	var result *entities.AccessToken = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT  *
	       FROM "access_tokens"
	       WHERE "uuid" = $1`, tokenUUID)

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		accessToken := &entities.AccessToken{}
		callbackErr = row.StructScan(accessToken)
		if callbackErr != nil {
			return commonPostgres.EmptyOrError(callbackErr, "unable get access token by uuid")
		}

		result = accessToken

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func NewPostgresStore(logger *zap.Logger,
	pgConn *commonPostgres.Connection,
) *pgRepository {
	return &pgRepository{
		pgConn: pgConn,
		logger: logger,
	}
}
