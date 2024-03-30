package pg_store

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
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
			  	"expired_at",           
       			"created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id;`,
			sessionItem.UUID,
			sessionItem.AccessTokenUUID, sessionItem.MnemonicWalletUUID,
			sessionItem.ExpiredAt,
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
) error {

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		res, callbackErr := stmt.Exec(`UPDATE "mnemonic_wallet_sessions" 
			SET "status" = $1,
				"updated_at" = now()
			WHERE "mnemonic_wallet_uuid" = $2`,
			walletUUID)
		if callbackErr != nil {
			return callbackErr
		}

		affectedRows, err := res.RowsAffected()
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

func (s *pgRepository) GetWalletSessionByUUID(ctx context.Context,
	sessionUUID string,
) (*entities.MnemonicWalletSession, error) {
	var result *entities.MnemonicWalletSession = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT *
	       FROM "mnemonic_wallet_sessions"
	       WHERE "uuid" = $1`, sessionUUID)

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		accessToken := &entities.MnemonicWalletSession{}
		callbackErr = row.StructScan(&accessToken)
		if callbackErr != nil {
			return commonPostgres.EmptyOrError(callbackErr, "unable get wallet session by uuid")
		}

		result = accessToken

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pgRepository) GetActiveWalletSessionsByWalletUUID(ctx context.Context) (
	count uint, list []*entities.MnemonicWalletSession, err error,
) {
	return s.GetActiveWalletSessionsClb(ctx, nopSessionCallback)
}

func (s *pgRepository) GetActiveWalletSessionsByWalletUUIDClb(ctx context.Context,
	walletUUID string,
	onItemCallBack func(item *entities.MnemonicWalletSession) error,
) (count uint, list []*entities.MnemonicWalletSession, err error) {
	if err = s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		rows, queryErr := stmt.Queryx(`SELECT *
			FROM "mnemonic_wallet_sessions"
			WHERE "mnemonic_wallet_uuid" = $1 AND
			      "now()" < "expired_at"
			ORDER BY "expired_at";`, walletUUID)
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
			WHERE "now()" < "expired_at"
			ORDER BY "expired_at";`)
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
