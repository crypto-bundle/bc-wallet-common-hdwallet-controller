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
	"math"
	"net"
	"time"

	"github.com/crypto-bundle/bc-wallet-common/pkg/dns"

	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	originGRPC "google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	grpcKeepalive "google.golang.org/grpc/keepalive"
)

const (
	DefaultClientMaxReceiveMessageSize = 1024 * 1024 * 24
	DefaultClientMaxSendMessageSize    = math.MaxInt32
)

func DefaultKeepaliveClientOptions() grpcKeepalive.ClientParameters {
	return grpcKeepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}
}

func DefaultRetryOptions() []grpcRetry.CallOption {
	return []grpcRetry.CallOption{
		grpcRetry.WithMax(3),
		grpcRetry.WithBackoff(grpcRetry.BackoffLinear(1000 * time.Millisecond)),
		grpcRetry.WithCodes(grpcCodes.Aborted, grpcCodes.Unavailable),
	}
}

func DefaultInterceptorsOptions() []originGRPC.UnaryClientInterceptor {
	return []originGRPC.UnaryClientInterceptor{
		grpcRetry.UnaryClientInterceptor(DefaultRetryOptions()...),
	}
}

func DefaultDialOptions() []originGRPC.DialOption {
	return []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithKeepaliveParams(DefaultKeepaliveClientOptions()),
		originGRPC.WithChainUnaryInterceptor(DefaultInterceptorsOptions()...),
	}
}

// Dialer is a method to pass it in grpc.Dial options
func Dialer(_ context.Context, target string) (net.Conn, error) {
	addr, err := dns.Resolve("grpc", "tcp", target)
	if err != nil {
		return nil, err
	}

	cn, err := net.Dial("tcp", addr)

	return cn, err
}

func Dial(target string, opts []originGRPC.DialOption) (*originGRPC.ClientConn, error) {
	dialOptions := opts
	if len(dialOptions) == 0 {
		dialOptions = DefaultDialOptions()
	}

	return originGRPC.Dial(target, dialOptions...)
}
