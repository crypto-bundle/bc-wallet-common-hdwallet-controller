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

package hdwallet_api

import (
	"context"
	"net"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

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
	handlers          pbApi.HdWalletApiServer
	config            configService

	listener net.Listener
}

func (s *Server) Init(ctx context.Context) error {
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

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("start close instances")

	s.grpcServer.GracefulStop()
	err := s.listener.Close()
	if err != nil {
		return err
	}

	s.logger.Info("grpc server shutdown completed")

	return nil
}

func (s *Server) ListenAndServe(ctx context.Context) (err error) {
	listenConn, err := net.Listen("tcp", s.config.GetBindPort())
	if err != nil {
		s.logger.Error("unable to listen port", zap.Error(err),
			zap.String("port", s.config.GetBindPort()))
	}
	s.listener = listenConn

	s.grpcServer = grpc.NewServer(s.grpcServerOptions...)
	if (s.config.IsDev() || s.config.IsLocal()) && s.config.IsDebug() {
		reflection.Register(s.grpcServer)
	}

	pbApi.RegisterHdWalletApiServer(s.grpcServer, s.handlers)

	s.logger.Info("grpc serve success")

	return s.grpcServer.Serve(s.listener)
}

// nolint:revive // fixme
func NewServer(ctx context.Context,
	loggerSrv *zap.Logger,
	cfg configService,
	handlers pbApi.HdWalletApiServer,
) (*Server, error) {
	l := loggerSrv.Named("grpc.server").With(
		zap.String(app.ApplicationNameTag, app.ApplicationName),
		zap.String(app.BlockChainNameTag, app.BlockChainName))

	srv := &Server{
		logger:   l,
		config:   cfg,
		handlers: handlers,
	}

	return srv, nil
}
