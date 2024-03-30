package pg_store

import (
	"context"
	"errors"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/entities"
	"github.com/jmoiron/sqlx"

	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

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

func (s *pgRepository) GetAccessTokenInfoByUUID(ctx context.Context,
	tokenUUID string,
) (*entities.AccessToken, error) {
	var result *entities.AccessToken = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT  "id", "uuid", "wallet_uuid",
        		"expired_at", "created_at", "updated_at"
	       FROM "access_tokens"
	       WHERE "uuid" = $1`, tokenUUID)

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		accessToken := &entities.AccessToken{}
		callbackErr = row.StructScan(&accessToken)
		if callbackErr != nil {
			return commonPostgres.EmptyOrError(callbackErr, "unable get access token by uuid")
		}

		result = accessToken

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *pgRepository) UpdateAccessTokenStatusByUUID(ctx context.Context,
	tokenUUID string,
) (*entities.AccessToken, error) {
	var result *entities.AccessToken = nil

	if err := s.pgConn.TryWithTransaction(ctx, func(stmt sqlx.Ext) error {
		row := stmt.QueryRowx(`SELECT  "id", "uuid", "wallet_uuid",
        		"expired_at", "created_at", "updated_at"
	       FROM "access_tokens"
	       WHERE "uuid" = $1`, tokenUUID)

		callbackErr := row.Err()
		if callbackErr != nil {
			return callbackErr
		}

		accessToken := &entities.AccessToken{}
		callbackErr = row.StructScan(&accessToken)
		if callbackErr != nil {
			return commonPostgres.EmptyOrError(callbackErr, "unable get access token by uuid")
		}

		result = accessToken

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func NewPostgresStore(logger *zap.Logger,
	pgConn *commonPostgres.Connection,
) *pgRepository {
	return &pgRepository{
		pgConn: pgConn,
		logger: logger,
	}
}
