package grpc

import (
	"context"
	"net"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/manager"

	commonGRPCServer "github.com/crypto-bundle/bc-wallet-common-lib-grpc/pkg/server"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	logger            *zap.Logger
	grpcServer        *grpc.Server
	grpcServerOptions []grpc.ServerOption
	handlers          pbApi.HdWalletManagerApiServer
	config            configService

	listener net.Listener
}

func (s *Server) Init(_ context.Context) error {
	options := commonGRPCServer.DefaultServeOptions()
	msgSizeOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(commonGRPCServer.DefaultServerMaxReceiveMessageSize),
		grpc.MaxSendMsgSize(commonGRPCServer.DefaultServerMaxSendMessageSize),
	}
	options = append(options, msgSizeOptions...)
	options = append(options, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())))

	s.grpcServerOptions = options

	return nil
}

func (s *Server) shutdown() error {
	s.logger.Info("start close instances")

	s.grpcServer.GracefulStop()
	//err := s.listener.Close()
	//if err != nil {
	//	return err
	//}

	s.logger.Info("grpc server shutdown completed")

	return nil
}

func (s *Server) ListenAndServe(ctx context.Context) (err error) {
	listenConn, err := net.Listen("tcp", s.config.GetBindPort())
	if err != nil {
		s.logger.Error("unable to listen port", zap.Error(err),
			zap.String("port", s.config.GetBindPort()))

		return err
	}
	s.listener = listenConn

	s.grpcServer = grpc.NewServer(s.grpcServerOptions...)
	if (s.config.IsDev() || s.config.IsLocal()) && s.config.IsDebug() {
		reflection.Register(s.grpcServer)
	}

	pbApi.RegisterHdWalletManagerApiServer(s.grpcServer, s.handlers)

	s.logger.Info("grpc serve success")

	go func() {
		err = s.grpcServer.Serve(s.listener)
		if err != nil {
			s.logger.Error("unable to start serving", zap.Error(err),
				zap.String("port", s.config.GetBindPort()))
		}
	}()

	<-ctx.Done()

	return s.shutdown()
}

// nolint:revive // fixme
func NewServer(ctx context.Context,
	loggerSrv *zap.Logger,
	cfg configService,
	handlers pbApi.HdWalletManagerApiServer,
) (*Server, error) {
	l := loggerSrv.Named("grpc.server")

	srv := &Server{
		logger:   l,
		config:   cfg,
		handlers: handlers,
	}

	return srv, nil
}
