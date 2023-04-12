package pgstore

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/entities"

	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrUnablePrepareQuery    = errors.New("unable to prepare query")
	ErrUnableExecuteQuery    = errors.New("unable to execute query")
	ErrUnableGetLastInsertID = errors.New("unable get last insert id")
	ErrHaveNoAffectedRows    = errors.New("have ho affected row")
)

type pgRepository struct {
	pgConn *commonPostgres.Connection
	logger *zap.Logger

	defaultOnScanMutator func(ctx context.Context, wallet *entities.Wallet) error
}

func (s *pgRepository) AddNewWallet(ctx context.Context, wallet *entities.Wallet) (*entities.Wallet, error) {
	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		var walletID uint32
		row := stmt.QueryRowx(`INSERT INTO "wallets" ("uuid", "title", "purpose", "is_enabled", "strategy",
       			"created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id;`,
			wallet.UUID.String(), wallet.Title, wallet.Purpose, wallet.IsEnabled,
			wallet.Strategy, wallet.CreatedAt, wallet.UpdatedAt)

		err := row.Scan(&walletID)
		if err != nil {
			s.logger.Error("failed to insert in wallets table", zap.Error(err))

			return ErrUnableExecuteQuery
		}

		wallet.ID = walletID

		return nil
	}); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *pgRepository) UpdateIsEnabledWalletByUUID(ctx context.Context, uuid string, isEnabled bool) error {
	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		date := time.Now()

		result, err := stmt.Exec(`UPDATE "wallets" 
			SET "is_enabled" = $1, "updated_at" = $2
			WHERE "uuid" = $3`, isEnabled, date, uuid)
		if err != nil {
			return err
		}

		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if affectedRows == 0 {
			return ErrHaveNoAffectedRows
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *pgRepository) GetWalletByUUID(ctx context.Context, walletUUID uuid.UUID) (*entities.Wallet, error) {
	var wallet *entities.Wallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT "id", "title", "uuid", "purpose", "is_enabled", "strategy",
       			"created_at", "updated_at"
	       FROM "wallets"
	       WHERE "uuid" = $1`, walletUUID.String())

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		wallet = &entities.Wallet{}
		callbackErr = row.StructScan(&wallet)
		if callbackErr != nil {
			return commonPostgres.EmptyOrError(callbackErr, "unable get wallet by uuid")
		}

		return nil

	}); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *pgRepository) GetAllEnabledWallets(ctx context.Context) ([]*entities.Wallet, error) {
	var wallets []*entities.Wallet = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, err := stmt.Queryx(`SELECT "id", "title", "uuid", "purpose", "is_enabled",  "strategy",
       			"created_at", "updated_at"
	       FROM "wallets"
	       WHERE "is_enabled" = true`)
		if err != nil {
			return err
		}

		defer rows.Close()

		err = rows.Err()
		if err != nil {
			return commonPostgres.EmptyOrError(err, "unable get all enabled wallets")
		}

		walletsList := make([]*entities.Wallet, 0)
		count := 0

		for rows.Next() {
			walletData := &entities.Wallet{}

			scanErr := rows.StructScan(walletData)
			if scanErr != nil {
				return err
			}

			walletsList = append(walletsList, walletData)
			count++
		}

		if count > 0 {
			wallets = walletsList
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (s *pgRepository) GetAllEnabledWalletUUIDList(ctx context.Context) ([]string, error) {
	var walletsUUIDList []string = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, err := stmt.Queryx(`SELECT "uuid"
	       FROM "wallets"
	       WHERE "is_enabled" = true`)

		if err != nil {
			return err
		}
		defer rows.Close()

		walletsUUIDList = make([]string, 0)

		for rows.Next() {
			var walletUUID string

			scanErr := rows.Scan(walletUUID)
			if scanErr != nil {
				return err
			}

			walletsUUIDList = append(walletsUUIDList, walletUUID)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return walletsUUIDList, nil
}

func NewPostgresStore(logger *zap.Logger,
	pgConn *commonPostgres.Connection,
) *pgRepository {
	return &pgRepository{
		pgConn: pgConn,
		logger: logger,
		defaultOnScanMutator: func(ctx context.Context, wallet *entities.Wallet) error {
			return nil
		},
	}
}
