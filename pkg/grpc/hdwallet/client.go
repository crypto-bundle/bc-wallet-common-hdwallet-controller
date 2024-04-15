package hdwallet

import (
	"context"
	"net"

	commonGRPCClient "github.com/crypto-bundle/bc-wallet-common-lib-grpc/pkg/client"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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
	options := []originGRPC.DialOption{
		originGRPC.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return net.Dial("unix", addr)
		}),
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
		originGRPC.WithChainUnaryInterceptor(commonGRPCClient.DefaultInterceptorsOptions()...),
	}
	msgSizeOptions := originGRPC.WithDefaultCallOptions(
		originGRPC.MaxCallRecvMsgSize(commonGRPCClient.DefaultClientMaxReceiveMessageSize),
		originGRPC.MaxCallSendMsgSize(commonGRPCClient.DefaultClientMaxSendMessageSize),
	)
	options = append(options, msgSizeOptions,
		originGRPC.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	s.grpcClientOptions = options

	return nil
}

func (s *Client) Dial(ctx context.Context) error {
	grpcConn, err := originGRPC.Dial(s.cfg.GetConnectionPath(), s.grpcClientOptions...)
	if err != nil {
		return err
	}
	s.grpcConn = grpcConn

	s.HdWalletApiClient = NewHdWalletApiClient(grpcConn)

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
