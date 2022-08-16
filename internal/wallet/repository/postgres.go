package repository

import (
	"context"
	"time"

	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/entities"

	"github.com/cryptowize-tech/bc-wallet-common/pkg/postgres"

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
	pgConn *postgres.Connection
	logger *zap.Logger
}

func (s *pgRepository) AddNewMnemonicWallet(ctx context.Context, wallet *entities.MnemonicWallet) (*entities.MnemonicWallet, error) {
	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		date := time.Now()

		var walletID uint32
		row := stmt.QueryRowx(`INSERT INTO "mnemonic_wallets" ("wallet_uuid", "hash", "purpose", "rsa_encrypted", 
				"rsa_encrypted_hash", "is_hot", "is_enabled", "created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`,
			wallet.UUID.String(), wallet.Hash, wallet.Purpose, wallet.RsaEncrypted, wallet.RsaEncryptedHash,
			wallet.IsHotWallet, true, date, date)

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
	wallet := &entities.MnemonicWallet{}

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT "id", "wallet_uuid", "hash", "purpose", "rsa_encrypted", 
				"rsa_encrypted_hash", "is_hot", "is_enabled", "created_at", "updated_at"
	       FROM "mnemonic_wallets"
	       WHERE "hash" = $1`, hash)

		err := row.StructScan(&wallet)
		if err != nil {
			return postgres.EmptyOrError(err, "unable get site by domain")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *pgRepository) GetMnemonicWalletUUID(ctx context.Context, uuid string) (*entities.MnemonicWallet, error) {
	wallet := &entities.MnemonicWallet{}

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT "id", "wallet_uuid", "hash", "purpose", "rsa_encrypted",
       			"rsa_encrypted_hash", "is_hot", "is_enabled", "created_at", "updated_at"
	       FROM "mnemonic_wallets"
	       WHERE "wallet_uuid" = $1`, uuid)

		err := row.StructScan(&wallet)
		if err != nil {
			return postgres.EmptyOrError(err, "unable get site by domain")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *pgRepository) GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	wallets := make([]*entities.MnemonicWallet, 0)

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, err := stmt.Queryx(`SELECT "id", "wallet_uuid", "hash", "purpose", "rsa_encrypted",
       			"rsa_encrypted_hash", "is_enabled", "created_at", "updated_at"
	       FROM "mnemonic_wallets"
	       WHERE "is_hot" = true`)

		if err != nil {
			return err
		}
		defer rows.Close()

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

func (s *pgRepository) GetAllEnabledMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	wallets := make([]*entities.MnemonicWallet, 0)

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, err := stmt.Queryx(`SELECT "id", "wallet_uuid", "hash", "purpose", "rsa_encrypted",
       			"rsa_encrypted_hash", "is_enabled", "created_at", "updated_at"
	       FROM "mnemonic_wallets"
	       WHERE "is_enabled" = true`)

		if err != nil {
			return err
		}
		defer rows.Close()

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

func (s *pgRepository) GetAllEnabledHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	wallets := make([]*entities.MnemonicWallet, 0)

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, err := stmt.Queryx(`SELECT "id", "wallet_uuid", "hash", "purpose", "rsa_encrypted",
       			"rsa_encrypted_hash", "is_enabled", "created_at", "updated_at"
	       FROM "mnemonic_wallets"
	       WHERE "is_hot" = true AND "is_enabled" = true`)

		if err != nil {
			return err
		}
		defer rows.Close()

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

func (s *pgRepository) GetAllEnabledNonHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	wallets := make([]*entities.MnemonicWallet, 0)

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, err := stmt.Queryx(`SELECT "id", "wallet_uuid", "hash", "purpose", "rsa_encrypted",
       			"rsa_encrypted_hash", "is_enabled", "created_at", "updated_at"
	       FROM "mnemonic_wallets"
	       WHERE "is_hot" = false AND "is_enabled" = true`)

		if err != nil {
			return err
		}
		defer rows.Close()

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
	pgConn *postgres.Connection,
) (*pgRepository, error) {
	return &pgRepository{
		pgConn: pgConn,
		logger: logger,
	}, nil
}
