package postgres

type DbConfig interface {
	GetHost() string
	GetPort() uint16
	GetDbName() string
	GetUser() string
	GetPassword() string
	GetTLSMode() string
	GetRetryCount() uint8
	GetConnectTimeOut() uint16

	GetMaxOpenConns() uint8
	GetMaxIdleConns() uint8

	IsDebug() bool
}
