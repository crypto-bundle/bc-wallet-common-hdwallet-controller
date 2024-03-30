package pg_store

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/types"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"

	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

		var walletID uint32
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
		err := row.Scan(&walletID)
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
	uuid uuid.UUID,
) (*entities.MnemonicWallet, error) {
	var wallet *entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT *
	       FROM "mnemonic_wallets"
	       WHERE "uuid" = $1`, uuid.String())

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
