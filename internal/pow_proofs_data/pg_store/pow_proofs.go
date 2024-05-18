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
	"fmt"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrUnablePrepareQuery = errors.New("unable to prepare query")
	ErrUnableExecuteQuery = errors.New("unable to execute query")
	ErrHaveNoAffectedRows = errors.New("have ho affected row")
)

var nopPowProofItemCallback = func(idx uint, item *entities.PowProof) error {
	return nil
}

type pgRepository struct {
	pgConn *commonPostgres.Connection
	logger *zap.Logger
}

func (s *pgRepository) AddNewPowProof(ctx context.Context,
	toSaveItem *entities.PowProof,
) (result *entities.PowProof, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		date := time.Now()

		row := stmt.QueryRowx(`INSERT INTO "pow_proofs" ("uuid", 
				"access_token_uuid",
			  	"message_check_nonce", "message_hash", "message_data",
				"created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING *;`,
			toSaveItem.UUID, toSaveItem.AccessTokenUUID,
			toSaveItem.MessageCheckNonce, toSaveItem.MessageHash, toSaveItem.MessageData,
			date, date)

		item := &entities.PowProof{}
		clbErr := row.StructScan(item)
		if clbErr != nil {
			s.logger.Error("failed to insert in pow-proof item", zap.Error(clbErr))

			return fmt.Errorf("%s: %w", ErrUnableExecuteQuery, clbErr)
		}

		result = item

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pgRepository) GetPowProofByUUID(ctx context.Context,
	uuid string,
) (powProof *entities.PowProof, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT *
	       FROM "pow_proofs"
	       WHERE "uuid" = $1`, uuid)

		clbErr := row.Err()
		if clbErr != nil {
			return clbErr
		}

		powProofItem := &entities.PowProof{}
		clbErr = row.StructScan(powProofItem)
		if clbErr != nil {
			return commonPostgres.EmptyOrError(clbErr, "unable get pow-proof item by uuid")
		}

		powProof = powProofItem

		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func (s *pgRepository) GetPowProofByMessageHash(ctx context.Context,
	messageHash string,
) (powProof *entities.PowProof, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT *
	       FROM "pow_proofs"
	       WHERE "message_hash" = $1`, messageHash)

		clbErr := row.Err()
		if clbErr != nil {
			return clbErr
		}

		powProofItem := &entities.PowProof{}
		clbErr = row.StructScan(powProofItem)
		if clbErr != nil {
			return commonPostgres.EmptyOrError(clbErr,
				"unable get pow-proof item by entity_uuid and entity_type")
		}

		powProof = powProofItem

		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func NewPostgresStore(logger *zap.Logger,
	pgConn *commonPostgres.Connection,
) *pgRepository {
	return &pgRepository{
		pgConn: pgConn,
		logger: logger,
	}
}
