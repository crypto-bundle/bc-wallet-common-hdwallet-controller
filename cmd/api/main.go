package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"bc-wallet-eth-hdwallet/internal/app"
	"bc-wallet-eth-hdwallet/internal/wallet"
	"bc-wallet-eth-hdwallet/pkg/grpc/hdwallet_api"

	"bc-wallet-eth-hdwallet/internal/config"
	grpcHandlers "bc-wallet-eth-hdwallet/internal/grpc"

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

	//vaultClient, err := vault.NewClient(ctx, vcfg)
	//if err != nil {
	//	loggerEntry.Fatal("vault client init error", zap.Error(err))
	//}
	//
	//cfg.VaultClient = vaultClient
	err = cfg.Prepare()
	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}

	conn := postgres.NewConnection(context.Background(), cfg, loggerEntry)
	_, err = conn.Connect()
	if err != nil {
		loggerEntry.Fatal(err.Error(), zap.Error(err))
	}

	listenConn, err := net.Listen("tcp", cfg.GetBindPort())
	if err != nil {
		loggerEntry.Fatal("unable to listen port", zap.Error(err),
			zap.String("port", cfg.GetBindPort()))
	}

	walletService, err := wallet.New(loggerEntry, cfg, conn)
	if err != nil {
		loggerEntry.Fatal("unable to create wallet service instance", zap.Error(err))
	}

	err = walletService.Init(ctx)
	if err != nil {
		loggerEntry.Fatal("unable to init wallet service", zap.Error(err))
	}

	apiHandlers, err := grpcHandlers.New(ctx, cfg, loggerEntry, walletService)
	if err != nil {
		loggerEntry.Fatal("unable to init grpc handlers", zap.Error(err))
	}

	srv, err := hdwallet_api.NewServer(ctx, loggerEntry, cfg, listenConn)
	if err != nil {
		loggerEntry.Fatal("unable to create grpc server instance", zap.Error(err),
			zap.String("port", cfg.Bind))
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

	loggerEntry.Info("application started successfully")

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
