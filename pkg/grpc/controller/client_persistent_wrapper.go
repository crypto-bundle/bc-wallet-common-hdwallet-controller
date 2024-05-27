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
	"fmt"
	"os"

	"go.uber.org/zap"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverAddressTemplate = "BC_WALLET_%s_HDWALLET_SERVICE_HOST"
	serverPortTemplate    = "BC_WALLET_%s_HDWALLET_SERVICE_PORT"
)

type ClientWrapperPersistentFlow struct {
	logger *zap.Logger

	originGRPCClient HdWalletControllerApiClient

	obscurityDataSvc   obscurityDataProvider
	accessTokenDataSvc accessTokensDataService
}

// Init bc-connector-tron service
func (s *ClientWrapperPersistentFlow) Init(_ context.Context) error {
	host := os.Getenv(fmt.Sprintf(serverAddressTemplate, s.processingNetwork))
	port := os.Getenv(fmt.Sprintf(serverPortTemplate, s.processingNetwork))
	s.serverAddress = fmt.Sprintf("%s:%s", host, port)

	options := []originGRPC.DialOption{
		originGRPC.WithReturnConnectionError(),
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		originGRPC.WithBlock(),
	}
	msgSizeOptions := originGRPC.WithDefaultCallOptions(
		originGRPC.MaxCallRecvMsgSize(1024*1024*3),
		originGRPC.MaxCallSendMsgSize(1024*1024*3),
	)
	options = append(options, msgSizeOptions)

	s.dialOptions = options

	return nil
}

func (s *ClientWrapperPersistentFlow) Dial(ctx context.Context) error {
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

// NewControllerClient is new instance of
func NewControllerClient(ctx context.Context,
	logger *zap.Logger,
	cfgSvc baseConfigService,
	obscurityDataSvc obscurityDataProvider,
	accessTokenDataSvc accessTokensDataService,
) (*ClientWrapperPersistentFlow, error) {
	return &ClientWrapperPersistentFlow{
		logger: logger,

		processingNetwork:  cfgSvc.GetNetworkName(),
		processingProvider: cfgSvc.GetProviderName(),

		obscurityDataSvc:   obscurityDataSvc,
		accessTokenDataSvc: accessTokenDataSvc,
	}, nil
}
