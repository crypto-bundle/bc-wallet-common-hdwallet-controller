/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
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
 *
 */

package main

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/pow_validator"
	"log"
	"os"
	"os/signal"
	"syscall"

	accessToken "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/access_tokens"
	accessTtokenData "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/access_tokens/pg_store"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/config"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/events"
	grpcHandlers "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/grpc"
	walletData "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/mnemonic_wallet_data/pg_store"
	walletRedisData "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/mnemonic_wallet_data/redis_store"
	powProofData "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/pow_proofs_data/pg_store"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/sign_manager"
	signReqData "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/sign_request_data/postgres"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/wallet_manager"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	commonJWT "github.com/crypto-bundle/bc-wallet-common-lib-jwt/pkg/jwt"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	commonProfiler "github.com/crypto-bundle/bc-wallet-common-lib-profiler/pkg/profiler"
	commonRedis "github.com/crypto-bundle/bc-wallet-common-lib-redis/pkg/redis"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"

	_ "github.com/mailru/easyjson/gen"
	"go.uber.org/zap"
)

// DO NOT EDIT THESE VARIABLES DIRECTLY. These are build-time constants
// DO NOT USE THESE VARIABLES IN APPLICATION CODE. USE commonConfig.NewLdFlagsManager SERVICE-COMPONENT INSTEAD OF IT
var (
	// ReleaseTag - release tag in TAG.SHORT_COMMIT_ID.BUILD_NUMBER.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ReleaseTag = "v0.0.0-00000000-100500"

	// CommitID - latest commit id.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	CommitID = "0000000000000000000000000000000000000000"

	// ShortCommitID - first 12 characters from CommitID.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ShortCommitID = "00000000"

	// BuildNumber - ci/cd build number for BuildNumber
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildNumber string = "100500"

	// BuildDateTS - ci/cd build date in time stamp
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildDateTS string = "1713280105"
)

func main() {
	var err error
	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	appCfg, vaultSvc, err := config.Prepare(ctx, ReleaseTag,
		CommitID, ShortCommitID,
		BuildNumber, BuildDateTS)
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	transitSvc := commonVault.NewEncryptService(vaultSvc, appCfg.GetVaultCommonTransit())
	encryptorSvc := commonVault.NewEncryptService(vaultSvc, appCfg.GetVaultCommonTransit())

	loggerSrv, err := commonLogger.NewService(appCfg)
	if err != nil {
		log.Fatal(err.Error(), err)
	}
	loggerEntry := loggerSrv.NewLoggerEntry("main").
		With(zap.String(app.BlockChainNameTag, appCfg.GetNetworkName()))

	profiler := commonProfiler.NewHTTPServer(loggerEntry, appCfg.ProfilerConfig)

	pgConn := commonPostgres.NewConnection(ctx, appCfg, loggerEntry)
	_, err = pgConn.Connect()
	if err != nil {
		loggerEntry.Fatal("unable to connect to postgresql", zap.Error(err))
	}
	loggerEntry.Info("postgresql connected")

	natsConnSvc := commonNats.NewConnection(ctx, appCfg, loggerEntry)
	err = natsConnSvc.Connect()
	if err != nil {
		loggerEntry.Fatal("unable to connect to nats", zap.Error(err))
	}
	loggerEntry.Info("nats connected")

	redisSvc := commonRedis.NewConnection(ctx, appCfg, loggerEntry)
	if err != nil {
		loggerEntry.Fatal("unable create redis connection", zap.Error(err))
	}

	redisConn, err := redisSvc.Connect(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to connect to redis", zap.Error(err))
	}
	redisClient := redisConn.GetClient()
	loggerEntry.Info("redis connected")

	jwtSvc := commonJWT.NewJWTService(appCfg.JWTConfig.Key)

	mnemonicWalletDataSvc := walletData.NewPostgresStore(loggerEntry, pgConn)
	mnemonicWalletCacheDataSvc := walletRedisData.NewRedisStore(loggerEntry, appCfg, redisClient)

	accessTokenDataSvc := accessTtokenData.NewPostgresStore(loggerEntry, pgConn)
	jwtDecoder := accessToken.NewJWTDecoder(loggerEntry, jwtSvc)

	signReqDataSvc := signReqData.NewPostgresStore(loggerEntry, pgConn)
	powProofDataSvc := powProofData.NewPostgresStore(loggerEntry, pgConn)
	powValidatorSvc := pow_validator.NewValidatorHashCash(loggerEntry)

	hdWalletClient := hdwallet.NewClient(appCfg)

	eventPublisher := events.NewEventsBroadcaster(appCfg, natsConnSvc, redisConn)

	walletSvc := wallet_manager.NewService(loggerEntry, appCfg, transitSvc, encryptorSvc,
		accessTokenDataSvc, mnemonicWalletDataSvc, mnemonicWalletCacheDataSvc, signReqDataSvc,
		jwtDecoder, hdWalletClient, eventPublisher, pgConn)
	signReqSvc := sign_manager.NewService(loggerEntry, signReqDataSvc, hdWalletClient, eventPublisher, pgConn)

	apiHandlers := grpcHandlers.New(loggerEntry, walletSvc, signReqSvc)
	apiInterceptors := grpcHandlers.NewInterceptorsList(loggerEntry, powProofDataSvc, mnemonicWalletDataSvc,
		accessTokenDataSvc, powValidatorSvc, jwtSvc, pgConn)

	GRPCSrv, err := grpcHandlers.NewServer(loggerEntry, appCfg, apiHandlers, apiInterceptors)
	if err != nil {
		loggerEntry.Fatal("unable to create grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	eventWatcher := events.NewEventWatcher(loggerEntry, appCfg, redisConn, natsConnSvc,
		mnemonicWalletCacheDataSvc, mnemonicWalletDataSvc, signReqDataSvc, hdWalletClient, pgConn)

	err = profiler.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init profiler", zap.Error(err))
	}
	loggerEntry.Info("profiler successfully initiated")

	err = hdWalletClient.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init hd-wallet grpc client", zap.Error(err),
			zap.String("unix path", appCfg.HdWalletClientConfig.GetConnectionPath()))
	}
	loggerEntry.Info("hd-wallet client initiated")

	err = GRPCSrv.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to listen init grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}
	loggerEntry.Info("gRPC server initiated")

	err = hdWalletClient.Dial(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to dial hd-wallet grpc client", zap.Error(err),
			zap.String("unix path", appCfg.HdWalletClientConfig.GetConnectionPath()))
	}
	loggerEntry.Info("hd-wallet client successfully connected")

	err = eventWatcher.Run(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to start event watcher service", zap.Error(err))
	}
	loggerEntry.Info("event-watcher started successfully")

	// TODO: add healthcheck flow
	//checker := commonHealthcheck.NewHTTPHealthChecker(loggerEntry, appCfg)
	//checker.AddStartupProbeUnit(vaultSvc)
	//checker.AddStartupProbeUnit(redisConn)
	//checker.AddStartupProbeUnit(pgConn)
	//checker.AddStartupProbeUnit(natsConnSvc)

	err = GRPCSrv.ListenAndServe(ctx)
	if err != nil {
		loggerEntry.Error("unable to start grpc", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	err = profiler.ListenAndServe(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init profiler", zap.Error(err))
	}

	loggerEntry.Info("application started successfully", zap.String(app.GRPCBindPortTag, appCfg.GetBindPort()))

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	loggerEntry.Warn("shutdown application")
	cancelCtxFunc()

	syncErr := loggerEntry.Sync()
	if syncErr != nil {
		log.Print(syncErr.Error(), syncErr)
	}

	log.Print("stopped")
}
