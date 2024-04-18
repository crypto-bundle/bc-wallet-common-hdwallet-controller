package pg_store

import (
	"context"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/types"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var nopSessionCallback = func(item *entities.MnemonicWalletSession) error {
	return nil
}

func (s *pgRepository) AddNewWalletSession(ctx context.Context,
	sessionItem *entities.MnemonicWalletSession,
) (*entities.MnemonicWalletSession, error) {
	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		var sessionID uint32
		row := stmt.QueryRowx(`INSERT INTO "mnemonic_wallet_sessions" ("uuid", 
				"access_token_uuid", "mnemonic_wallet_uuid",
				"status",
			  	"started_at", "expired_at",           
       			"created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`,
			sessionItem.UUID,
			sessionItem.AccessTokenUUID, sessionItem.MnemonicWalletUUID,
			sessionItem.Status,
			sessionItem.StartedAt, sessionItem.ExpiredAt,
			sessionItem.CreatedAt, sessionItem.UpdatedAt)

		err := row.Scan(&sessionID)
		if err != nil {
			s.logger.Error("failed to insert in wallets table", zap.Error(err))

			return ErrUnableExecuteQuery
		}

		sessionItem.ID = sessionID

		return nil
	}); err != nil {
		return nil, err
	}

	return sessionItem, nil
}

func (s *pgRepository) UpdateWalletSessionStatusByWalletUUID(ctx context.Context,
	walletUUID string,
	newStatus types.MnemonicWalletSessionStatus,
) error {

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		_, callbackErr := stmt.Exec(`UPDATE "mnemonic_wallet_sessions" 
			SET "status" = $1,
				"updated_at" = now()
			WHERE "mnemonic_wallet_uuid" = $2`,
			newStatus, walletUUID)
		if callbackErr != nil {
			return callbackErr
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *pgRepository) UpdateWalletSessionStatusBySessionUUID(ctx context.Context,
	sessionUUID string,
	newStatus types.MnemonicWalletSessionStatus,
) (result *entities.MnemonicWalletSession, err error) {

	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`UPDATE "mnemonic_wallet_sessions" 
			SET "status" = $1,
				"updated_at" = now()
			WHERE "uuid" = $2
			RETURNING *`,
			newStatus, sessionUUID)

		sessionItem := &entities.MnemonicWalletSession{}

		clbErr := row.Scan(sessionItem)
		if clbErr != nil {
			s.logger.Error("unable to update wallet session status", zap.Error(clbErr))

			return fmt.Errorf("%s: %w", ErrUnableExecuteQuery.Error(), clbErr)
		}

		result = sessionItem

		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func (s *pgRepository) UpdateMultipleWalletSessionStatus(ctx context.Context,
	walletsUUIDs []string,
	newStatus types.MnemonicWalletSessionStatus,
) (count uint, list []string, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, clbErr := sqlx.In(`UPDATE "mnemonic_wallet_sessions"
	       SET "status" = ?
	       WHERE "mnemonic_wallet_uuid" IN (?)
	       RETURNING "uuid"`, newStatus, walletsUUIDs)

		bonded := stmt.Rebind(query)
		returnedRows, clbErr := stmt.Queryx(bonded, args...)
		if clbErr != nil {
			return clbErr
		}
		defer returnedRows.Close()

		sessionsUUIDList := make([]string, 0)

		for returnedRows.Next() {
			var sessionUUID string

			scanErr := returnedRows.Scan(&sessionUUID)
			if scanErr != nil {
				return scanErr
			}

			sessionsUUIDList = append(sessionsUUIDList, sessionUUID)

			count++
		}

		list = sessionsUUIDList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) GetWalletSessionByUUID(ctx context.Context,
	sessionUUID string,
) (*entities.MnemonicWalletSession, error) {
	var result *entities.MnemonicWalletSession = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, clbErr := sqlx.In(`SELECT *
			FROM "mnemonic_wallet_sessions"
			WHERE "uuid" = ? AND
					now() BETWEEN "started_at" AND "expired_at" 
				AND
			    	"status" IN (?)
			ORDER BY "expired_at";`, sessionUUID, []types.MnemonicWalletSessionStatus{
			types.MnemonicWalletSessionStatusPrepared,
		})
		if clbErr != nil {
			return clbErr
		}
		bonded := stmt.Rebind(query)
		row := stmt.QueryRowx(bonded, args...)
		if clbErr != nil {
			return clbErr
		}

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		session := &entities.MnemonicWalletSession{}
		callbackErr = row.StructScan(&session)
		if callbackErr != nil {
			return commonPostgres.EmptyOrError(callbackErr, "unable get session by UUID")
		}

		result = session

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pgRepository) GetActiveWalletSessionsByWalletUUID(ctx context.Context, walletUUID string) (
	count uint, list []*entities.MnemonicWalletSession, err error,
) {
	return s.GetWalletSessionsByWalletUUIDAndStatusClb(ctx, walletUUID,
		[]types.MnemonicWalletSessionStatus{
			types.MnemonicWalletSessionStatusPrepared,
		},
		nopSessionCallback)
}

func (s *pgRepository) GetWalletSessionsByWalletUUIDAndStatusClb(ctx context.Context,
	walletUUID string,
	sessionStatuses []types.MnemonicWalletSessionStatus,
	onItemCallBack func(item *entities.MnemonicWalletSession) error,
) (count uint, list []*entities.MnemonicWalletSession, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		query, args, clbErr := sqlx.In(`SELECT *
			FROM "mnemonic_wallet_sessions"
			WHERE "mnemonic_wallet_uuid" = ? AND
					now() BETWEEN "started_at" AND "expired_at" 
				AND
			      	"status" IN (?)
			ORDER BY "expired_at";`, walletUUID, sessionStatuses)
		if clbErr != nil {
			return clbErr
		}

		bonded := stmt.Rebind(query)
		rows, clbErr := stmt.Queryx(bonded, args...)
		if clbErr != nil {
			return clbErr
		}
		defer rows.Close()

		itemsList := make([]*entities.MnemonicWalletSession, 0)
		itemsCount := uint(0)

		for rows.Next() {
			data := &entities.MnemonicWalletSession{}

			scanErr := rows.StructScan(data)
			if scanErr != nil {
				return scanErr
			}

			scanErr = onItemCallBack(data)
			if scanErr != nil {
				return scanErr
			}

			itemsList = append(itemsList, data)
			itemsCount++
		}

		if itemsCount == 0 {
			itemsList = nil
		}

		count = itemsCount
		list = itemsList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}

func (s *pgRepository) GetActiveWalletSessions(ctx context.Context) (
	count uint, list []*entities.MnemonicWalletSession, err error,
) {
	return s.GetActiveWalletSessionsClb(ctx, nopSessionCallback)
}

func (s *pgRepository) GetActiveWalletSessionsClb(ctx context.Context,
	onItemCallBack func(item *entities.MnemonicWalletSession) error,
) (count uint, list []*entities.MnemonicWalletSession, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, queryErr := stmt.Queryx(`SELECT *
			FROM "mnemonic_wallet_sessions"
			WHERE now() BETWEEN "started_at" AND "expired_at" 
				AND "status" = $1
			ORDER BY "expired_at";`, types.MnemonicWalletSessionStatusPrepared)
		if queryErr != nil {
			return queryErr
		}

		defer rows.Close()

		itemsList := make([]*entities.MnemonicWalletSession, 0)
		itemsCount := uint(0)

		for rows.Next() {
			data := &entities.MnemonicWalletSession{}

			scanErr := rows.StructScan(data)
			if scanErr != nil {
				return scanErr
			}

			scanErr = onItemCallBack(data)
			if scanErr != nil {
				return scanErr
			}

			itemsList = append(itemsList, data)
			itemsCount++
		}

		if itemsCount == 0 {
			itemsList = nil
		}

		count = itemsCount
		list = itemsList

		return nil
	}); err != nil {
		return 0, nil, err
	}

	return
}
