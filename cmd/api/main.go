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
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/config"
	grpcHandlers "github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/grpc"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/mnemonic"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/mnemonic_wallet_data"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/wallet_data"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/wallet_manager"
	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api"

	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	"github.com/crypto-bundle/bc-wallet-common/pkg/crypter"

	"go.uber.org/zap"
)

// DO NOT EDIT THIS VARIABLES DIRECTLY. These are build-time constants
// DO NOT USE THESE VARIABLES IN APPLICATION CODE. USE commonConfig.NewLdFlagsManager SERVICE-COMPONENT INSTEAD OF IT
var (
	// Version - version time.RFC3339.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	Version = "DEVELOPMENT.VESION"

	// ReleaseTag - release tag in TAG.%Y-%m-%dT%H-%M-%SZ.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	ReleaseTag = "DEVELOPMENT.RELEASE_TAG"

	// CommitID - latest commit id.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	CommitID = "DEVELOPMENT.COMMIT_HASH"

	// ShortCommitID - first 12 characters from CommitID.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	ShortCommitID = "DEVELOPMENT.SHORT_COMMIT_HASH"

	// BuildNumber - ci/cd build number for BuildNumber
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	BuildNumber uint64 = 0

	// BuildDateTS - ci/cd build date in time stamp
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	BuildDateTS uint64 = 0
)

func main() {
	var err error
	ctx := context.Background()

	appCfg, err := config.Prepare(ctx, Version, ReleaseTag,
		CommitID, ShortCommitID,
		BuildNumber, BuildDateTS)
	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	loggerSrv, err := commonLogger.NewService(appCfg)
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	loggerEntry, err := loggerSrv.NewLoggerEntry("main")
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	pgConn := commonPostgres.NewConnection(context.Background(), appCfg, loggerEntry)
	_, err = pgConn.Connect()
	if err != nil {
		loggerEntry.Fatal(err.Error(), zap.Error(err))
	}

	listenConn, err := net.Listen("tcp", appCfg.GetBindPort())
	if err != nil {
		loggerEntry.Fatal("unable to listen port", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	cryptoService := crypter.New(appCfg.HDWalletConfig)
	if err != nil {
		loggerEntry.Fatal("unable to create crypto service instance", zap.Error(err))
	}

	walletDataSrv := wallet_data.NewService(loggerEntry, pgConn)
	mnemonicWalletDataSrv := mnemonic_wallet_data.NewService(loggerEntry, pgConn)
	mnemonicGenerator := mnemonic.NewMnemonicGenerator(loggerEntry)

	walletService, err := wallet_manager.NewService(loggerEntry, appCfg, cryptoService,
		walletDataSrv, mnemonicWalletDataSrv,
		pgConn, mnemonicGenerator)
	if err != nil {
		loggerEntry.Fatal("unable to create wallet service instance", zap.Error(err))
	}

	apiHandlers, err := grpcHandlers.New(ctx, loggerEntry, walletService)
	if err != nil {
		loggerEntry.Fatal("unable to init grpc handlers", zap.Error(err))
	}

	srv, err := hdwallet_api.NewServer(ctx, loggerEntry, appCfg, listenConn)
	if err != nil {
		loggerEntry.Fatal("unable to create grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	err = walletService.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init wallet service", zap.Error(err))
	}

	err = srv.Init(ctx, loggerEntry, apiHandlers)
	if err != nil {
		loggerEntry.Fatal("unable to listen init grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	go func() {
		err = srv.ListenAndServe(ctx)
		if err != nil {
			loggerEntry.Fatal("unable to start grpc handlers", zap.Error(err),
				zap.String("port", appCfg.GetBindPort()))
		}
	}()

	loggerEntry.Info("application started successfully", zap.String(app.GRPCBindPortTag, appCfg.GetBindPort()))

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
