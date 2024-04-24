package main

import (
	"context"
	"log"
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

	commonHealthcheck "github.com/crypto-bundle/bc-wallet-common-lib-healthcheck/pkg/healthcheck"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
	commonRedis "github.com/crypto-bundle/bc-wallet-common-lib-redis/pkg/redis"

	_ "github.com/mailru/easyjson/gen"
	"go.uber.org/zap"
)

// DO NOT EDIT THESE VARIABLES DIRECTLY. These are build-time constants
// DO NOT USE THESE VARIABLES IN APPLICATION CODE. USE commonConfig.NewLdFlagsManager SERVICE-COMPONENT INSTEAD OF IT
var (
	// Version - version time.RFC3339.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	Version = "DEVELOPMENT.VESION"

	// ReleaseTag - release tag in TAG.%Y-%m-%dT%H-%M-%SZ.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ReleaseTag = "DEVELOPMENT.RELEASE_TAG"

	// CommitID - latest commit id.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	CommitID = "DEVELOPMENT.COMMIT_HASH"

	// ShortCommitID - first 12 characters from CommitID.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ShortCommitID = "DEVELOPMENT.SHORT_COMMIT_HASH"

	// BuildNumber - ci/cd build number for BuildNumber
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildNumber uint64 = 0

	// BuildDateTS - ci/cd build date in time stamp
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildDateTS uint64 = 0
)

const ApplicationName = "bc-wallet-tron-hdwallet-api"

func main() {
	var err error
	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	appCfg, secretSrv, err := config.Prepare(ctx, Version, ReleaseTag,
		CommitID, ShortCommitID,
		BuildNumber, BuildDateTS, ApplicationName)
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	loggerSrv, err := commonLogger.NewService(appCfg)
	if err != nil {
		log.Fatal(err.Error(), err)
	}
	loggerEntry := loggerSrv.NewLoggerEntry("main")

	pgConn := commonPostgres.NewConnection(context.Background(), appCfg, loggerEntry)
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

	walletDataSrv := wallet_data.NewService(loggerEntry, pgConn)
	mnemonicWalletDataSrv, err := mnemonic_wallet_data.NewService(loggerEntry, appCfg,
		pgConn, redisClient, natsConnSvc)
	if err != nil {
		loggerEntry.Fatal("unable to create mnemonic wallet data service", zap.Error(err))
	}

	mnemonicGenerator := mnemonic.NewMnemonicGenerator(loggerEntry,
		appCfg.GetDefaultMnemonicWordsCount())

	walletService, err := wallet_manager.NewService(loggerEntry, appCfg, secretSrv,
		walletDataSrv, mnemonicWalletDataSrv,
		pgConn, mnemonicGenerator)
	if err != nil {
		loggerEntry.Fatal("unable to create wallet service instance", zap.Error(err))
	}

	apiHandlers := grpcHandlers.New(ctx, loggerEntry, walletService)

	srv, err := grpcHandlers.NewServer(ctx, loggerEntry, appCfg, apiHandlers)
	if err != nil {
		loggerEntry.Fatal("unable to create grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	err = walletService.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init wallet service", zap.Error(err))
	}

	err = walletService.Run(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to run wallet service", zap.Error(err))
	}

	err = srv.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to listen init grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	commonHealthcheck.NewHTTPHealthChecker(loggerEntry)

	go func() {
		err = srv.ListenAndServe(ctx)
		if err != nil {
			loggerEntry.Error("unable to start grpc", zap.Error(err),
				zap.String("port", appCfg.GetBindPort()))
		}
	}()

	loggerEntry.Info("application started successfully", zap.String(app.GRPCBindPortTag, appCfg.GetBindPort()))

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	loggerEntry.Warn("shutdown application")
	cancelCtxFunc()

	walletShutdownErr := walletService.Shutdown(ctx)
	if walletShutdownErr != nil {
		log.Printf("%s:%e", walletShutdownErr.Error(), walletShutdownErr)
	}

	syncErr := loggerEntry.Sync()
	if syncErr != nil {
		log.Print(syncErr.Error(), syncErr)
	}

	log.Print("stopped")
}
