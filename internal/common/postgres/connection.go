package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type connectionParams struct {
	host     string
	port     uint16
	user     string
	password string
	database string

	retryTimeOut time.Duration
	retryCount   uint8

	readTimeOut  int
	writeTimeOut int

	maxOpenConn uint8
	maxIdleConn uint8

	sslMode string

	debug bool
}

// Connection struct to store and manipulate postgres database connection
type Connection struct {
	Dbx    *sqlx.DB
	params *connectionParams
	l      *zap.Logger
}

func (c *Connection) Close() error {
	err := c.Dbx.Close()
	if err != nil {
		return err
	}

	return nil
}

// Connect to Clickhouse database
func (c *Connection) Connect() (*Connection, error) {
	retryDecValue := uint8(1)
	retryCount := c.params.retryCount
	if retryCount == 0 {
		retryDecValue = 0
		retryCount = 1
	}
	try := 0

	for i := retryCount; i != 0; i -= retryDecValue {
		dbx, err := connectToPostgres(context.Background(), c.params)
		if err != nil {
			c.l.Error("unable connect to postgres. reconnecting...",
				zap.Error(err), zap.Int("iteration", try))
			try++
			time.Sleep(c.params.retryTimeOut)

			continue
		}

		err = dbx.Ping()
		if err != nil {
			c.l.Error("unable ping postgres. reconnecting...",
				zap.Error(err), zap.Int("iteration", try),
				zap.Any("params", c.params))
			try++
			time.Sleep(c.params.retryTimeOut)

			continue
		}

		rows, err := dbx.Query("SELECT 1")
		if err != nil {
			c.l.Error("unable make sql request. reconnecting...",
				zap.Error(err), zap.Int("iteration", try),
				zap.Any("params", c.params))
			try++
			time.Sleep(c.params.retryTimeOut)

			continue
		}
		err = rows.Close()
		if err != nil {
			c.l.Error("unable to close rows statement. reconnecting...",
				zap.Error(err), zap.Int("iteration", try),
				zap.Any("params", c.params))
			try++
			time.Sleep(c.params.retryTimeOut)

			continue
		}

		dbx.SetMaxOpenConns(int(c.params.maxOpenConn))
		dbx.SetMaxIdleConns(int(c.params.maxIdleConn))

		c.Dbx = dbx
		return c, nil
	}

	return c, nil
}

// NewConnection to postgres db
func NewConnection(ctx context.Context, cfg DbConfig, logger *zap.Logger) *Connection {
	conn := &Connection{
		params: &connectionParams{
			host:     cfg.GetHost(),
			port:     cfg.GetPort(),
			user:     cfg.GetUser(),
			password: cfg.GetPassword(),
			database: cfg.GetDbName(),

			retryCount:   cfg.GetRetryCount(),
			retryTimeOut: time.Duration(cfg.GetConnectTimeOut()) * time.Millisecond,

			maxOpenConn: cfg.GetMaxOpenConns(),
			maxIdleConn: cfg.GetMaxIdleConns(),

			debug: cfg.IsDebug(),

			sslMode: cfg.GetTLSMode(),
		},
		l: logger,
	}

	return conn
}

func connectToPostgres(ctx context.Context, params *connectionParams) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		params.user, params.password, params.host, params.port, params.database, params.sslMode))
}
