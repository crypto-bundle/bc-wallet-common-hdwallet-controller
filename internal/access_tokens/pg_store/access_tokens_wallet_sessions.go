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

func (s *pgRepository) AddNewAccessTokenWalletSession(ctx context.Context,
	toSaveItem *entities.AccessTokenWalletSession,
) (result *entities.AccessTokenWalletSession, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		date := time.Now()

		row := stmt.QueryRowx(`INSERT INTO "access_tokens_wallet_sessions" ("token_uuid", "wallet_session_uuid",
				 "created_at")
            VALUES($1, $2, $3) RETURNING *;`,
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

func (s *pgRepository) GetWalletSessionByTokenUUID(ctx context.Context,
	tokenUUID string,
) (*entities.AccessTokenWalletSession, error) {
	var result *entities.AccessTokenWalletSession = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT *
	       FROM "access_tokens_wallet_sessions"
	       WHERE "token_uuid" = $1`, tokenUUID)

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		item := &entities.AccessTokenWalletSession{}
		callbackErr = row.StructScan(&item)
		if callbackErr != nil {
			return commonPostgres.EmptyOrError(callbackErr, "unable get record by token_uuid")
		}

		result = item

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}
