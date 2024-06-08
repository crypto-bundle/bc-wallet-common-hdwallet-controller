/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package grpc

import (
	"context"
	"errors"
	"net"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	commonGRPCServer "github.com/crypto-bundle/bc-wallet-common-lib-grpc/pkg/server"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type WalletApiServer struct {
	logger            *zap.Logger
	grpcServer        *grpc.Server
	grpcServerOptions []grpc.ServerOption

	handlers         pbApi.HdWalletControllerWalletApiServer
	interceptorsList []grpc.UnaryServerInterceptor

	config configService

	listener net.Listener
}

func (s *WalletApiServer) Init(_ context.Context) error {
	options := commonGRPCServer.DefaultServeOptions()
	msgSizeOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(commonGRPCServer.DefaultServerMaxReceiveMessageSize),
		grpc.MaxSendMsgSize(commonGRPCServer.DefaultServerMaxSendMessageSize),
		grpc.ChainUnaryInterceptor(s.interceptorsList...),
	}
	options = append(options, msgSizeOptions...)
	options = append(options, grpc.StatsHandler(otelgrpc.NewServerHandler()))

	s.grpcServerOptions = options

	return nil
}

func (s *WalletApiServer) shutdown() error {
	s.logger.Info("start close instances")

	s.grpcServer.GracefulStop()
	//err := s.listener.Close()
	//if err != nil {
	//	return err
	//}

	s.logger.Info("grpc server shutdown completed")

	return nil
}

func (s *WalletApiServer) ListenAndServe(ctx context.Context) (err error) {
	listenConn, err := net.Listen("tcp", s.config.GetWalletApiBindAddress())
	if err != nil {
		s.logger.Error("unable to listen port", zap.Error(err),
			zap.String("bind address", s.config.GetWalletApiBindAddress()))

		return err
	}
	s.listener = listenConn

	s.grpcServer = grpc.NewServer(s.grpcServerOptions...)
	if (s.config.IsDev() || s.config.IsLocal()) && s.config.IsDebug() {
		reflection.Register(s.grpcServer)
	}

	go s.serve(ctx)

	return nil
}

func (s *WalletApiServer) serve(ctx context.Context) {
	newCtx, causeFunc := context.WithCancelCause(ctx)
	pbApi.RegisterHdWalletControllerWalletApiServer(s.grpcServer, s.handlers)

	go func() {
		err := s.grpcServer.Serve(s.listener)
		if err != nil {
			s.logger.Error("unable to start serving", zap.Error(err),
				zap.String("bind address", s.config.GetWalletApiBindAddress()))
			causeFunc(err)
		}
	}()

	<-newCtx.Done()
	intErr := newCtx.Err()
	if !errors.Is(intErr, context.Canceled) {
		s.logger.Error("ctx cause errors", zap.Error(intErr))
	}

	intErr = s.shutdown()
	if intErr != nil {
		s.logger.Error("unable to graceful shutdown http server", zap.Error(intErr))
	}

	return
}

// nolint:revive // fixme
func NewWalletApiServer(loggerSrv *zap.Logger,
	cfg configService,
	handlers pbApi.HdWalletControllerWalletApiServer,
	interceptorsList []grpc.UnaryServerInterceptor,
) (*WalletApiServer, error) {
	l := loggerSrv.Named("grpc.server")

	srv := &WalletApiServer{
		logger: l,

		config: cfg,

		handlers:         handlers,
		interceptorsList: interceptorsList,
	}

	return srv, nil
}
