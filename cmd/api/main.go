/*
 * MIT License
 *
 * Copyright (c) 2021-2023 Aleksei Kotelnikov
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"context"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/mnemonic"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/mnemonic_wallet_data"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/wallet_data"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/config"
	grpcHandlers "github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/grpc"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/wallet_manager"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api"

	"github.com/crypto-bundle/bc-wallet-common/pkg/crypter"
	"github.com/crypto-bundle/bc-wallet-common/pkg/logger"
	"github.com/crypto-bundle/bc-wallet-common/pkg/postgres"
	"github.com/crypto-bundle/bc-wallet-common/pkg/vault"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

func main() {
	var err error
	ctx := context.Background()
	baseCfg := &config.BaseConfig{}
	err = baseCfg.Prepare()
	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	cfg := &config.Config{}
	cfg.BaseConfig = baseCfg

	if baseCfg.IsDev() {
		loadErr := godotenv.Load(".env")
		if loadErr != nil {
			log.Fatal(loadErr)
		}
	}

	loggerSrv, err := logger.NewService(baseCfg)
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	loggerEntry, err := loggerSrv.NewLoggerEntry("main")
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	vcfg := &vault.Config{}
	err = envconfig.Process(app.VaultConfigPrefix, vcfg)
	if err != nil {
		loggerEntry.Fatal("vault config init error", zap.Error(err))
	}

	// vaultClient, err := vault.NewClient(ctx, vcfg)
	// if err != nil {
	//	loggerEntry.Fatal("vault client init error", zap.Error(err))
	// }
	//
	// cfg.VaultClient = vaultClient
	err = cfg.Prepare()
	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	pgConn := postgres.NewConnection(context.Background(), cfg, loggerEntry)
	_, err = pgConn.Connect()
	if err != nil {
		loggerEntry.Fatal(err.Error(), zap.Error(err))
	}

	listenConn, err := net.Listen("tcp", cfg.GetBindPort())
	if err != nil {
		loggerEntry.Fatal("unable to listen port", zap.Error(err),
			zap.String("port", cfg.GetBindPort()))
	}

	cryptoService := crypter.New(cfg.HDWalletConfig)
	if err != nil {
		loggerEntry.Fatal("unable to create crypto service instance", zap.Error(err))
	}

	walletDataSrv := wallet_data.NewService(loggerEntry, pgConn)
	mnemonicWalletDataSrv := mnemonic_wallet_data.NewService(loggerEntry, pgConn)
	mnemonicGenerator := mnemonic.NewMnemonicGenerator(loggerEntry)

	walletService, err := wallet_manager.NewService(loggerEntry, cfg, cryptoService,
		walletDataSrv, mnemonicWalletDataSrv,
		pgConn, mnemonicGenerator)
	if err != nil {
		loggerEntry.Fatal("unable to create wallet service instance", zap.Error(err))
	}

	apiHandlers, err := grpcHandlers.New(ctx, loggerEntry, walletService)
	if err != nil {
		loggerEntry.Fatal("unable to init grpc handlers", zap.Error(err))
	}

	srv, err := hdwallet_api.NewServer(ctx, loggerEntry, cfg, listenConn)
	if err != nil {
		loggerEntry.Fatal("unable to create grpc server instance", zap.Error(err),
			zap.String("port", cfg.Bind))
	}

	err = walletService.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init wallet service", zap.Error(err))
	}

	err = srv.Init(ctx, loggerEntry, apiHandlers)
	if err != nil {
		loggerEntry.Fatal("unable to listen init grpc server instance", zap.Error(err),
			zap.String("port", cfg.Bind))
	}

	go func() {
		err = srv.ListenAndServe(ctx)
		if err != nil {
			loggerEntry.Fatal("unable to start grpc handlers", zap.Error(err),
				zap.String("port", cfg.Bind))
		}
	}()

	loggerEntry.Info("application started successfully", zap.String(app.GRPCBindPortTag, cfg.GetBindPort()))

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	loggerEntry.Warn("shutdown application")
	srv.Shutdown(ctx)

	walletShutdownErr := walletService.Shutdown(ctx)
	if walletShutdownErr != nil {
		log.Fatal(walletShutdownErr.Error(), walletShutdownErr)
	}

	syncErr := loggerEntry.Sync()
	if syncErr != nil {
		log.Fatal(syncErr.Error(), syncErr)
	}

	log.Print("stopped")
}
