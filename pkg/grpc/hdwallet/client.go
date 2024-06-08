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

package hdwallet

import (
	"context"
	"fmt"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	cfg hdWalletClientConfig

	grpcClientOptions []originGRPC.DialOption

	grpcConn *originGRPC.ClientConn
	HdWalletApiClient
}

// Init bc-wallet-tron-hdwallet GRPC-client service
// nolint:revive // fixme (autofix)
func (s *Client) Init(ctx context.Context) error {
	diallerSvc := newSocketDialer(s.cfg.GetConnectionPath(), s.cfg.GetUnixFileNameTemplate())

	dialFunc, err := diallerSvc.Prepare()
	if err != nil {
		return err
	}

	options := []originGRPC.DialOption{
		originGRPC.WithContextDialer(dialFunc),
		originGRPC.WithReturnConnectionError(),
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		originGRPC.WithBlock(),
	}
	msgSizeOptions := originGRPC.WithDefaultCallOptions(
		originGRPC.MaxCallRecvMsgSize(1024*1024*3),
		originGRPC.MaxCallSendMsgSize(1024*1024*3),
	)
	options = append(options, msgSizeOptions)

	s.grpcClientOptions = options

	return nil
}

func (s *Client) Dial(ctx context.Context) error {
	grpcConn, dialErr := originGRPC.Dial(s.cfg.GetConnectionPath(), s.grpcClientOptions...)
	if dialErr != nil {
		return dialErr
	}

	if grpcConn == nil {
		return fmt.Errorf("%w: %s", ErrUnableToFindActiveFileSocket, s.cfg.GetConnectionPath())
	}

	s.grpcConn = grpcConn
	s.HdWalletApiClient = NewHdWalletApiClient(s.grpcConn)

	go func() {
		<-ctx.Done()

		_ = s.shutdown()
	}()

	return nil
}

// Shutdown grpc hdwallet-client service
func (s *Client) shutdown() error {
	err := s.grpcConn.Close()
	if err != nil {
		return err
	}

	return nil
}

// nolint:revive // fixme
func NewClient(cfg hdWalletClientConfig) *Client {
	srv := &Client{
		cfg: cfg,
	}

	return srv
}
