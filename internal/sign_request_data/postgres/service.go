package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	"github.com/jmoiron/sqlx"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

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

func (s *pgRepository) AddSignRequestItem(ctx context.Context,
	toSaveItem *entities.SignRequest,
) (*entities.SignRequest, error) {
	var result *entities.SignRequest = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		date := time.Now()

		row := stmt.QueryRowx(`INSERT INTO "sign_requests" ("uuid", 
				"mnemonic_wallet_uuid", "session_uuid", "purpose_uuid", 
				"status",
				"created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING *;`,
			toSaveItem.UUID,
			toSaveItem.WalletUUID, toSaveItem.SessionUUID, toSaveItem.PurposeUUID,
			toSaveItem.Status,
			date, nil)

		signReqItem := &entities.SignRequest{}
		clbErr := row.StructScan(signReqItem)
		if clbErr != nil {
			s.logger.Error("failed to insert in sign_requests table", zap.Error(clbErr))

			return fmt.Errorf("%w: %s", clbErr, ErrUnableExecuteQuery.Error())
		}

		result = signReqItem

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pgRepository) UpdateSignRequestItemStatus(ctx context.Context,
	signReqUUID string,
	newStatus types.SignRequestStatus,
) error {
	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		_, clbErr := stmt.Exec(`UPDATE "sign_requests" 
			SET "status" = $1,
			    "updated_at" = now()
			WHERE "uuid" = $2`, newStatus,
			signReqUUID)
		if clbErr != nil {
			return commonPostgres.EmptyOrError(clbErr, "unable to update sign_requests item status")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *pgRepository) UpdateSignRequestItemStatusBySessionUUID(ctx context.Context,
	sessionUUID string,
	newStatus types.SignRequestStatus,
) (count uint, list []*entities.SignRequest, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, clbErr := stmt.Queryx(`UPDATE "sign_requests" 
			SET "status" = $1,
				"updated_at" = now()
			WHERE "session_uuid" = $2
			RETURNING *`,
			newStatus, sessionUUID)
		if clbErr != nil {
			return commonPostgres.EmptyOrError(clbErr,
				"unable to update request items by session uuid")
		}

		defer rows.Close()

		signRequestsList := make([]*entities.SignRequest, 0)

		for rows.Next() {
			updatedReq := &entities.SignRequest{}

			scanErr := rows.StructScan(updatedReq)
			if scanErr != nil {
				return scanErr
			}

			signRequestsList = append(signRequestsList, updatedReq)
			count++
		}

		list = signRequestsList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) UpdateSignRequestItemStatusByWalletsUUIDList(ctx context.Context,
	walletUUIDs []string,
	newStatus types.SignRequestStatus,
) (count uint, list []*entities.SignRequest, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, clbErr := sqlx.In(`UPDATE "sign_requests"
	       	SET "status" = ?,
	    		"updated_at" = now()
	       	WHERE "mnemonic_wallet_uuid" IN (?)
	       	RETURNING *`, newStatus, walletUUIDs)

		bonded := stmt.Rebind(query)
		returnedRows, clbErr := stmt.Queryx(bonded, args...)
		if clbErr != nil {
			return commonPostgres.EmptyOrError(clbErr,
				"unable to update request items by wallets uuid list")
		}

		defer returnedRows.Close()

		signRequestsList := make([]*entities.SignRequest, 0)

		for returnedRows.Next() {
			updatedReq := &entities.SignRequest{}

			scanErr := returnedRows.StructScan(updatedReq)
			if scanErr != nil {
				return scanErr
			}

			signRequestsList = append(signRequestsList, updatedReq)

			count++
		}

		list = signRequestsList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) UpdateSignRequestItemStatusByWalletUUID(ctx context.Context,
	walletUUID string,
	newStatus types.SignRequestStatus,
) (count uint, list []*entities.SignRequest, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, clbErr := stmt.Queryx(`UPDATE "sign_requests" 
			SET "status" = $1,
				"updated_at" = now()
			WHERE "mnemonic_wallet_uuid" = $2
			RETURNING *`,
			newStatus, walletUUID)
		if clbErr != nil {
			return commonPostgres.EmptyOrError(clbErr,
				"unable to update request items by wallet uuid")
		}

		defer rows.Close()

		signRequestsList := make([]*entities.SignRequest, 0)

		for rows.Next() {
			updatedReq := &entities.SignRequest{}

			scanErr := rows.StructScan(updatedReq)
			if scanErr != nil {
				return scanErr
			}

			signRequestsList = append(signRequestsList, updatedReq)
			count++
		}

		list = signRequestsList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) GetSignRequestItemByUUIDAndStatus(ctx context.Context,
	signReqUUID string,
	status types.SignRequestStatus,
) (*entities.SignRequest, error) {
	var item *entities.SignRequest = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT *
	       FROM "sign_requests"
	       WHERE "uuid" = $1 AND "status" = $2`, signReqUUID, status)

		queryErr := row.Err()
		if queryErr != nil {
			return queryErr
		}

		item = &entities.SignRequest{}
		err := row.StructScan(item)
		if err != nil {
			return commonPostgres.EmptyOrError(err, "unable get sign request by uuid")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return item, nil
}

func NewPostgresStore(logger *zap.Logger,
	pgConn *commonPostgres.Connection,
) *pgRepository {
	return &pgRepository{
		pgConn: pgConn,
		logger: logger,
	}
}
