package hdwallet_api

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/app"
	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/config"
	pbApi "github.com/crypto-bundle/bc-wallet-eth-hdwallet/pkg/grpc/hdwallet_api/proto"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	logger     *zap.Logger
	grpcServer *grpc.Server
	handlers   pbApi.HdWalletApiServer
	config     *config.Config

	listener net.Listener
}

func (s *Server) Init(ctx context.Context,
	loggerEntry *zap.Logger,
	handlers pbApi.HdWalletApiServer,
) error {
	s.handlers = handlers

	loggerEntry.Info("init success")

	return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	s.logger.Info("start close instances")

	s.grpcServer.Stop()

	s.logger.Info("grpc server shutdown completed")
}

func (s *Server) ListenAndServe(ctx context.Context) (err error) {
	// todo: move to go-base
	options := DefaultServeOptions()
	msgSizeOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(DefaultServerMaxReceiveMessageSize),
		grpc.MaxSendMsgSize(DefaultServerMaxSendMessageSize),
	}
	options = append(options, msgSizeOptions...)
	options = append(options, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())))

	s.grpcServer = grpc.NewServer(options...)
	reflection.Register(s.grpcServer)
	pbApi.RegisterHdWalletApiServer(s.grpcServer, s.handlers)

	s.logger.Info("grpc serve success")

	return s.grpcServer.Serve(s.listener)
}

// nolint:revive // fixme
func NewServer(ctx context.Context,
	loggerSrv *zap.Logger,
	cfg *config.Config,
	listener net.Listener,
) (*Server, error) {
	l := loggerSrv.Named("grpc.server").With(
		zap.String(app.ApplicationNameTag, app.ApplicationName),
		zap.String(app.BlockChainNameTag, app.BlockChainName))

	srv := &Server{
		logger:   l,
		config:   cfg,
		listener: listener,
	}

	return srv, nil
}
