package main

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/events"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/app"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/config"
	grpcHandlers "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/grpc"
	walletData "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/mnemonic_wallet_data/pg_store"
	walletRedisData "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/mnemonic_wallet_data/redis_store"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/sign_manager"
	signReqData "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/sign_request_data/postgres"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/wallet_manager"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

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

	mnemonicWalletDataSvc := walletData.NewPostgresStore(loggerEntry, pgConn)
	mnemonicWalletCacheDataSvc := walletRedisData.NewRedisStore(loggerEntry, appCfg, redisClient)

	signReqDataSvc := signReqData.NewPostgresStore(loggerEntry, pgConn)

	hdWalletClient := hdwallet.NewClient(appCfg)

	eventPublisher := events.NewEventsBroadcaster(appCfg, natsConnSvc, redisConn)

	walletSvc := wallet_manager.NewService(loggerEntry, appCfg, transitSvc, encryptorSvc,
		mnemonicWalletDataSvc, mnemonicWalletCacheDataSvc, signReqDataSvc,
		hdWalletClient, eventPublisher, pgConn)
	signReqSvc := sign_manager.NewService(loggerEntry, signReqDataSvc, hdWalletClient, eventPublisher, pgConn)

	apiHandlers := grpcHandlers.New(loggerEntry, walletSvc, signReqSvc)

	GRPCSrv, err := grpcHandlers.NewServer(ctx, loggerEntry, appCfg, apiHandlers)
	if err != nil {
		loggerEntry.Fatal("unable to create grpc server instance", zap.Error(err),
			zap.String("port", appCfg.GetBindPort()))
	}

	eventWatcher := events.NewEventWatcher(loggerEntry, appCfg, redisConn, natsConnSvc,
		mnemonicWalletCacheDataSvc, mnemonicWalletDataSvc, signReqDataSvc, hdWalletClient, pgConn)

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
