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
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"

	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrUnablePrepareQuery = errors.New("unable to prepare query")
	ErrUnableExecuteQuery = errors.New("unable to execute query")
	ErrHaveNoAffectedRows = errors.New("have ho affected row")
)

type pgRepository struct {
	pgConn *commonPostgres.Connection
	logger *zap.Logger
}

func (s *pgRepository) AddNewMnemonicWallet(ctx context.Context,
	toSaveItem *entities.MnemonicWallet,
) (*entities.MnemonicWallet, error) {
	var result *entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		date := time.Now()

		row := stmt.QueryRowx(`INSERT INTO "mnemonic_wallets" ("uuid", 
				"mnemonic_hash",
				"status", "unload_interval",
				"vault_encrypted", "vault_encrypted_hash", 
				"created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;`,
			toSaveItem.UUID.String(),
			toSaveItem.MnemonicHash,
			toSaveItem.Status, toSaveItem.UnloadInterval,
			toSaveItem.VaultEncrypted, toSaveItem.VaultEncryptedHash,
			date, date)

		wallet := &entities.MnemonicWallet{}
		err := row.StructScan(wallet)
		if err != nil {
			s.logger.Error("failed to insert in mnemonic_wallets", zap.Error(err))

			return ErrUnableExecuteQuery
		}

		result = wallet

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pgRepository) UpdateMultipleWalletsStatus(ctx context.Context,
	walletUUIDs []string,
	newStatus types.MnemonicWalletStatus,
) (count uint, list []string, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, clbErr := sqlx.In(`UPDATE "mnemonic_wallets"
	       SET "status" = ?
	       WHERE "uuid" IN (?)
	       RETURNING "uuid"`, newStatus, walletUUIDs)

		bonded := stmt.Rebind(query)
		returnedRows, clbErr := stmt.Queryx(bonded, args...)
		if clbErr != nil {
			return clbErr
		}
		defer returnedRows.Close()

		walletsUUIDList := make([]string, 0)

		for returnedRows.Next() {
			var walletUUID string

			scanErr := returnedRows.Scan(&walletUUID)
			if scanErr != nil {
				return scanErr
			}

			walletsUUIDList = append(walletsUUIDList, walletUUID)

			count++
		}

		list = walletsUUIDList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) UpdateWalletStatus(ctx context.Context,
	walletUUID string,
	newStatus types.MnemonicWalletStatus,
) (*entities.MnemonicWallet, error) {
	var wallet *entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`UPDATE "mnemonic_wallets" 
			SET "status" = $1,
			    "updated_at" = now()
			WHERE "uuid" = $2
			RETURNING *`, newStatus,
			walletUUID)

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		wallet = &entities.MnemonicWallet{}
		callbackErr = row.StructScan(wallet)
		if callbackErr != nil {
			return commonPostgres.EmptyOrError(callbackErr, "unable get mnemonic wallet by hash")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *pgRepository) UpdateMultipleWalletsStatusRetWallets(ctx context.Context,
	walletUUIDs []string,
	newStatus types.MnemonicWalletStatus,
) (count uint, list []*entities.MnemonicWallet, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, clbErr := sqlx.In(`UPDATE "mnemonic_wallets"
	       SET "status" = ?
	       WHERE "uuid" IN (?)
	       RETURNING *`, newStatus, walletUUIDs)

		bonded := stmt.Rebind(query)
		returnedRows, clbErr := stmt.Queryx(bonded, args...)
		if clbErr != nil {
			return clbErr
		}
		defer returnedRows.Close()

		walletsList := make([]*entities.MnemonicWallet, 0)

		for returnedRows.Next() {
			wallet := &entities.MnemonicWallet{}

			scanErr := returnedRows.StructScan(wallet)
			if scanErr != nil {
				return scanErr
			}

			walletsList = append(walletsList, wallet)

			count++
		}

		list = walletsList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error) {
	var wallet *entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT  *
	       FROM "mnemonic_wallets"
	       WHERE "mnemonic_hash" = $1`, hash)

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		wallet = &entities.MnemonicWallet{}
		callbackErr = row.StructScan(&wallet)
		if callbackErr != nil {
			return commonPostgres.EmptyOrError(callbackErr, "unable get mnemonic wallet by hash")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *pgRepository) GetMnemonicWalletByUUID(ctx context.Context,
	uuid string,
) (*entities.MnemonicWallet, error) {
	var wallet *entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT *
	       FROM "mnemonic_wallets"
	       WHERE "uuid" = $1`, uuid)

		queryErr := row.Err()
		if queryErr != nil {
			return queryErr
		}

		wallet = &entities.MnemonicWallet{}
		err := row.StructScan(wallet)
		if err != nil {
			return commonPostgres.EmptyOrError(err, "unable get mnemonic wallet by uuid")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *pgRepository) GetMnemonicWalletsByStatus(ctx context.Context,
	status types.MnemonicWalletStatus,
) ([]*entities.MnemonicWallet, error) {
	var wallets []*entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, err := stmt.Queryx(`SELECT *
	       FROM "mnemonic_wallets"
	       WHERE "status" = $1`, status)
		if err != nil {
			return err
		}

		defer rows.Close()

		wallets = make([]*entities.MnemonicWallet, 0)

		for rows.Next() {
			walletData := &entities.MnemonicWallet{}

			scanErr := rows.StructScan(walletData)
			if scanErr != nil {
				return scanErr
			}

			wallets = append(wallets, walletData)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (s *pgRepository) GetMnemonicWalletsByUUIDListAndStatus(ctx context.Context,
	UUIDList []string,
	statuses []types.MnemonicWalletStatus,
) ([]string, []*entities.MnemonicWallet, error) {
	var wallets []*entities.MnemonicWallet = nil
	var walletsUUIDs []string = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, err := sqlx.In(`SELECT *
	       FROM "mnemonic_wallets"
	       WHERE "uuid" IN (?) AND "status" IN (?)`, UUIDList, statuses)

		bonded := stmt.Rebind(query)
		returnedRows, err := stmt.Queryx(bonded, args...)
		if err != nil {
			return err
		}
		defer returnedRows.Close()
		wallets = make([]*entities.MnemonicWallet, 0)

		for returnedRows.Next() {
			walletData := &entities.MnemonicWallet{}

			scanErr := returnedRows.StructScan(walletData)
			if scanErr != nil {
				return err
			}

			wallets = append(wallets, walletData)
			walletsUUIDs = append(walletsUUIDs, walletData.UUID.String())
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return walletsUUIDs, wallets, nil
}

func (s *pgRepository) GetMnemonicWalletsByUUIDList(ctx context.Context,
	UUIDList []string,
) ([]*entities.MnemonicWallet, error) {
	var wallets []*entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, err := sqlx.In(`SELECT *
	       FROM "mnemonic_wallets"
	       WHERE "wallet_uuid" IN (?)`, UUIDList)

		bonded := stmt.Rebind(query)
		returnedRows, err := stmt.Queryx(bonded, args...)
		if err != nil {
			return err
		}
		defer returnedRows.Close()
		wallets = make([]*entities.MnemonicWallet, 0)

		for returnedRows.Next() {
			walletData := &entities.MnemonicWallet{}

			scanErr := returnedRows.StructScan(walletData)
			if scanErr != nil {
				return err
			}

			wallets = append(wallets, walletData)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (s *pgRepository) GetAllNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	var wallets []*entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, err := stmt.Queryx(`SELECT *
	       FROM "mnemonic_wallets"
	       WHERE "is_hot" = false`)
		if err != nil {
			return err
		}
		defer rows.Close()

		wallets = make([]*entities.MnemonicWallet, 0)

		for rows.Next() {
			walletData := &entities.MnemonicWallet{}

			scanErr := rows.StructScan(walletData)
			if scanErr != nil {
				return err
			}

			wallets = append(wallets, walletData)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return wallets, nil
}

func NewPostgresStore(logger *zap.Logger,
	pgConn *commonPostgres.Connection,
) *pgRepository {
	return &pgRepository{
		pgConn: pgConn,
		logger: logger,
	}
}
