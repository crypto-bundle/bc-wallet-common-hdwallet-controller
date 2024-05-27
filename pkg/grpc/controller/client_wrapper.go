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

package controller

import (
	"context"
	"go.uber.org/zap"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientWrapper struct {
	logger *zap.Logger

	processingNetwork  string
	processingProvider string

	serverAddress string
	serverPort    string

	txStmtManager      transactionalStatementManager
	obscurityDataSvc   obscurityDataProvider
	accessTokenDataSvc accessTokensDataService

	dialOptions      []originGRPC.DialOption
	grpcConn         *originGRPC.ClientConn
	originGRPCClient HdWalletControllerApiClient
}

// Init bc-connector-tron service
func (s *ClientWrapper) Init(_ context.Context, cfg hdWalletClientConfigService) error {
	s.serverAddress = cfg.GetServerBindAddress()

	options := []originGRPC.DialOption{
		originGRPC.WithReturnConnectionError(),
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		originGRPC.WithBlock(),
	}
	msgSizeOptions := originGRPC.WithDefaultCallOptions(
		originGRPC.MaxCallRecvMsgSize(cfg.GetMaxReceiveMessageSize()),
		originGRPC.MaxCallSendMsgSize(cfg.GetMaxSendMessageSize()),
	)
	options = append(options, msgSizeOptions)

	var interceptors []originGRPC.UnaryClientInterceptor

	if cfg.IsAccessTokenShieldEnabled() {
		interceptors = append(interceptors,
			newAccessTokenInterceptor(s.accessTokenDataSvc).Invoke)
	}

	if cfg.IsPowShieldEnabled() {
		interceptors = append(interceptors,
			newPowShieldInterceptor(s.obscurityDataSvc, s.accessTokenDataSvc, s.txStmtManager).Invoke)
	}

	if len(interceptors) > 0 {
		options = append(options, originGRPC.WithChainUnaryInterceptor(interceptors...))
	}

	s.dialOptions = options

	//s.powTarget = big.NewInt(1).
	//	Lsh(s.powTarget, uint(256-8))

	return nil
}

func (s *ClientWrapper) Dial(ctx context.Context) error {
	grpcConn, err := originGRPC.Dial(s.serverAddress, s.dialOptions...)
	if err != nil {
		return err
	}

	s.grpcConn = grpcConn
	s.originGRPCClient = NewHdWalletControllerApiClient(grpcConn)

	go func() {
		<-ctx.Done()

		closeErr := s.grpcConn.Close()
		if closeErr != nil {
			s.logger.Error("unable to close hd-wallet controller connection")
		}

	}()

	return nil
}

func NewClientWrapper(logger *zap.Logger,
	obscurityDataSvc obscurityDataProvider,
	accessTokenDataSvc accessTokensDataService,
	txStmtManager transactionalStatementManager,
) *ClientWrapper {
	return &ClientWrapper{
		logger: logger,

		txStmtManager:      txStmtManager,
		obscurityDataSvc:   obscurityDataSvc,
		accessTokenDataSvc: newAccessTokenDataWrapper(accessTokenDataSvc),
	}
}
