package hdwallet_api

import (
	"context"
	"math"
	"net"
	"time"

	"github.com/cryptowize-tech/bc-wallet-common/pkg/dns"

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
