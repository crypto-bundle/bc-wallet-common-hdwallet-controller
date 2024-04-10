package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/config"
	grpcHandlers "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/grpc"
	walletData "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/mnemonic_wallet_data/pg_store"
	walletRedisData "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/mnemonic_wallet_data/redis_store"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/sign_manager"
	signReqData "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/sign_request_data/postgres"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/internal/wallet_manager"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/hdwallet"

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
	Version = "DEVELOPMENT.VERSION"

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

func main() {
	var err error
	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	appCfg, _, err := config.Prepare(ctx, Version, ReleaseTag,
		CommitID, ShortCommitID,
		BuildNumber, BuildDateTS)
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	loggerSrv, err := commonLogger.NewService(appCfg)
	if err != nil {
		log.Fatal(err.Error(), err)
	}
	loggerEntry := loggerSrv.NewLoggerEntry("main").
		With(zap.String(app.BlockChainNameTag, appCfg.GetNetworkName()))

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

	mnemonicWalletDataSvc := walletData.NewPostgresStore(loggerEntry, pgConn)
	mnemonicWalletCacheDataSvc := walletRedisData.NewRedisStore(loggerEntry, appCfg, redisClient)

	signReqDataSvc := signReqData.NewPostgresStore(loggerEntry, pgConn)

	hdWalletClient := hdwallet.NewClient(appCfg)

	walletSvc := wallet_manager.NewService(loggerEntry, appCfg, mnemonicWalletDataSvc, mnemonicWalletCacheDataSvc, signReqDataSvc,
		hdWalletClient, pgConn)
	signReqSvc := sign_manager.NewService(loggerEntry, signReqDataSvc, hdWalletClient, pgConn)

	apiHandlers := grpcHandlers.New(loggerEntry, walletSvc, signReqSvc)

	GRPCSrv, err := grpcHandlers.NewServer(ctx, loggerEntry, appCfg, apiHandlers)
	if err != nil {
		loggerEntry.Fatal("unable to create grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	err = hdWalletClient.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init hd-wallet grpc client", zap.Error(err),
			zap.String("unix path", appCfg.HdWalletClientConfig.GetConnectionPath()))
	}

	err = GRPCSrv.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to listen init grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	err = hdWalletClient.Dial(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to dial hd-wallet grpc client", zap.Error(err),
			zap.String("unix path", appCfg.HdWalletClientConfig.GetConnectionPath()))
	}

	commonHealthcheck.NewHTTPHealthChecker(loggerEntry)

	go func() {
		err = GRPCSrv.ListenAndServe(ctx)
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

	syncErr := loggerEntry.Sync()
	if syncErr != nil {
		log.Print(syncErr.Error(), syncErr)
	}

	log.Print("stopped")
}
