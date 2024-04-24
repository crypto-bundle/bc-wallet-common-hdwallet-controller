package pg_store

import (
	"context"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"

	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

	"github.com/google/uuid"
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
				"is_hot", "unload_interval",
				"rsa_encrypted", "rsa_encrypted_hash", "vault_encrypted", "vault_encrypted_hash", 
				"created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;`,
			wallet.UUID.String(), wallet.WalletUUID.String(),
			wallet.MnemonicHash,
			wallet.IsHotWallet, wallet.UnloadInterval,
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
		row := stmt.QueryRowx(`SELECT  "id", "uuid", "wallet_uuid", "mnemonic_hash", 
       			"rsa_encrypted", "rsa_encrypted_hash",
       			"vault_encrypted", "vault_encrypted_hash",
       			"is_hot", "unload_interval", 
       			"created_at", "updated_at"
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
		row := stmt.QueryRowx(`SELECT "id", "uuid", "wallet_uuid", "mnemonic_hash", 
       			"rsa_encrypted", "rsa_encrypted_hash",
       			"vault_encrypted", "vault_encrypted_hash",
       			"is_hot", "unload_interval",  
       			"created_at", "updated_at"
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

func (s *pgRepository) GetAllHotMnemonicWallets(ctx context.Context) ([]*entities.MnemonicWallet, error) {
	var wallets []*entities.MnemonicWallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, err := stmt.Queryx(`SELECT "id", "uuid", "wallet_uuid", "mnemonic_hash", 
       			"rsa_encrypted", "rsa_encrypted_hash",
       			"vault_encrypted", "vault_encrypted_hash",
       			"is_hot", "unload_interval",  
       			"created_at", "updated_at"
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
		query, args, err := sqlx.In(`SELECT "id", "uuid", "wallet_uuid", "mnemonic_hash", 
       			"rsa_encrypted", "rsa_encrypted_hash",
       			"vault_encrypted", "vault_encrypted_hash",
       			"is_hot", "unload_interval", 
       			"created_at", "updated_at"
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
		rows, err := stmt.Queryx(`SELECT  "id", "uuid", "wallet_uuid", "mnemonic_hash", 
       			"rsa_encrypted", "rsa_encrypted_hash",
       			"vault_encrypted", "vault_encrypted_hash",
       			"is_hot", "unload_interval", 
       			"created_at", "updated_at"
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
