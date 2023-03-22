/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
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
 */

package pgstore

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"

	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

func (s *pgRepository) AddNewMnemonicWallet(ctx context.Context,
	wallet *entities.MnemonicWallet,
) (*entities.MnemonicWallet, error) {
	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		date := time.Now()

		var walletID uint32
		row := stmt.QueryRowx(`INSERT INTO "mnemonic_wallets" ("uuid", "wallet_uuid", 
				"mnemonic_hash",
				"is_hot", 
				"rsa_encrypted", "rsa_encrypted_hash", "vault_encrypted", "vault_encrypted_hash", 
				"created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;`,
			wallet.UUID.String(), wallet.WalletUUID.String(),
			wallet.MnemonicHash,
			wallet.IsHotWallet,
			wallet.RsaEncrypted, wallet.RsaEncryptedHash, wallet.VaultEncrypted, wallet.VaultEncryptedHash,
			date, date)

		err := row.Scan(&walletID)
		if err != nil {
			s.logger.Error("failed to insert in mnemonic_wallets", zap.Error(err))

			return ErrUnableExecuteQuery
		}

		wallet.ID = walletID

		return nil
	}); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *pgRepository) GetMnemonicWalletByHash(ctx context.Context, hash string) (*entities.MnemonicWallet, error) {
	var wallet *entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT  "id", "wallet_uuid", "mnemonic_hash", 
       			"rsa_encrypted", "rsa_encrypted_hash",
       			"vault_encrypted", "vault_encrypted_hash",
       			"is_hot", "created_at", "updated_at"
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

func (s *pgRepository) GetMnemonicWalletByUUID(ctx context.Context, uuid string) (*entities.MnemonicWallet, error) {
	var wallet *entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT  "id", "wallet_uuid", "mnemonic_hash", 
       			"rsa_encrypted", "rsa_encrypted_hash",
       			"vault_encrypted", "vault_encrypted_hash",
       			"is_hot", "created_at", "updated_at"
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

func (s *pgRepository) GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	var wallets []*entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, err := stmt.Queryx(`SELECT "id", "wallet_uuid", "mnemonic_hash", 
       			"rsa_encrypted", "rsa_encrypted_hash",
       			"vault_encrypted", "vault_encrypted_hash",
       			"is_hot", "created_at", "updated_at"
	       FROM "mnemonic_wallets"
	       WHERE "is_hot" = true`)
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

func (s *pgRepository) GetMnemonicWalletsByUUIDList(ctx context.Context,
	UUIDList []string,
) ([]*entities.MnemonicWallet, error) {
	var wallets []*entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, err := sqlx.In(`SELECT "id", "wallet_uuid", "mnemonic_hash", 
       			"rsa_encrypted", "rsa_encrypted_hash",
       			"vault_encrypted", "vault_encrypted_hash",
       			"is_hot", "created_at", "updated_at"
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
		rows, err := stmt.Queryx(`SELECT  "id", "wallet_uuid", "mnemonic_hash", 
       			"rsa_encrypted", "rsa_encrypted_hash",
       			"vault_encrypted", "vault_encrypted_hash",
       			"is_hot", "created_at", "updated_at"
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
