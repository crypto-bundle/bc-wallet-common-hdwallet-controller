package postgres

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
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

		var walletID uint32
		row := stmt.QueryRowx(`INSERT INTO "sign_requests" ("uuid", 
				"mnemonic_wallet_uuid", "session_uuid",
				"status",
				"created_at", "updated_at")
            VALUES($1, $2, $3, $4, $5, $6) RETURNING *;`,
			toSaveItem.UUID,
			toSaveItem.WalletUUID, toSaveItem.SessionUUID,
			toSaveItem.Status,
			date, nil)

		signReqItem := &entities.SignRequest{}
		err := row.Scan(&walletID)
		if err != nil {
			s.logger.Error("failed to insert in sign_requests table", zap.Error(err))

			return ErrUnableExecuteQuery
		}

		result = signReqItem

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}
