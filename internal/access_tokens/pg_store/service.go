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
	"errors"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/jmoiron/sqlx"

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

func (s *pgRepository) GetAccessTokenInfoByUUID(ctx context.Context,
	tokenUUID string,
) (*entities.AccessToken, error) {
	var result *entities.AccessToken = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT  "id", "uuid", "wallet_uuid",
        		"expired_at", "created_at", "updated_at"
	       FROM "access_tokens"
	       WHERE "uuid" = $1`, tokenUUID)

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		accessToken := &entities.AccessToken{}
		callbackErr = row.StructScan(&accessToken)
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

func (s *pgRepository) UpdateAccessTokenStatusByUUID(ctx context.Context,
	tokenUUID string,
) (*entities.AccessToken, error) {
	var result *entities.AccessToken = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT  "id", "uuid", "wallet_uuid",
        		"expired_at", "created_at", "updated_at"
	       FROM "access_tokens"
	       WHERE "uuid" = $1`, tokenUUID)

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		accessToken := &entities.AccessToken{}
		callbackErr = row.StructScan(&accessToken)
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
