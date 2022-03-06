package config

import (
	"fmt"

	"bc-wallet-eth-hdwallet/internal/app"

	"github.com/kelseyhightower/envconfig"
)

type vaultDbCredentials struct {
	DbUsername *string `json:"DB_USERNAME,omitempty"`
	DbPassword *string `json:"DB_PASSWORD,omitempty"`
}

type DBConfig struct {
	DbHost         string `envconfig:"DB_HOST" json:"-"`
	DbPort         uint16 `envconfig:"DB_PORT" json:"-"`
	DbName         string `envconfig:"DB_DATABASE" json:"-"`
	DbUsername     string `envconfig:"DB_USERNAME" json:"DB_USERNAME"`
	DbPassword     string `envconfig:"DB_PASSWORD" json:"DB_PASSWORD"`
	DbSSLMode      string `envconfig:"DB_SSL_MODE" default:"prefer" json:"-"`
	DbMaxOpenConns uint8  `envconfig:"DB_MAX_OPEN_CONNECTIONS" default:"8" json:"-"`
	DbMaxIdleConns uint8  `envconfig:"DB_MAX_IDLE_CONNECTIONS" default:"8" json:"-"`
	// DbConnectRetryCount is the maximum number of reconnection tries. If 0 - infinite loop
	DbConnectRetryCount uint8 `envconfig:"DB_RETRY_COUNT" default:"0" json:"-"`
	// DbConnectTimeOut is the timeout in millisecond to connect between connection tries
	DbConnectTimeOut uint16 `envconfig:"DB_RETRY_TIMEOUT" default:"5000" json:"-"`

	// --- CALCULATED ---
	vaultData []byte
}

func (c *DBConfig) Prepare() error {
	err := envconfig.Process(app.DatabaseConfigPrefix, c)
	if err != nil {
		return err
	}

	//credentials := vaultDbCredentials{}
	//err = json.Unmarshal(c.vaultData, &credentials)
	//if err != nil {
	//	return err
	//}
	//
	//if credentials.DbUsername != nil {
	//	c.DbUsername = *credentials.DbUsername
	//}
	//
	//if credentials.DbPassword != nil {
	//	c.DbPassword = *credentials.DbPassword
	//}

	return nil
}

func (c *DBConfig) GetDatabaseDSN() string {
	return fmt.Sprintf("mysql://%s:%s@%s/%s?sslmode=%t",
		c.DbUsername, c.DbPassword, c.DbHost, c.DbName, c.DbSSLMode)
}

func (c *DBConfig) GetHost() string {
	return c.DbHost
}

func (c *DBConfig) GetPort() uint16 {
	return c.DbPort
}

func (c *DBConfig) GetDbName() string {
	return c.DbName
}

func (c *DBConfig) GetUser() string {
	return c.DbUsername
}

func (c *DBConfig) GetPassword() string {
	return c.DbPassword
}

func (c *DBConfig) GetTLSMode() string {
	return c.DbSSLMode
}

func (c *DBConfig) GetRetryCount() uint8 {
	return c.DbConnectRetryCount
}

func (c *DBConfig) GetConnectTimeOut() uint16 {
	return c.DbConnectTimeOut
}

func (c *DBConfig) GetMaxOpenConns() uint8 {
	return c.DbMaxOpenConns
}

func (c *DBConfig) GetMaxIdleConns() uint8 {
	return c.DbMaxIdleConns
}
