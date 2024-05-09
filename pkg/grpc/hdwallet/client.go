package hdwallet

import (
	"context"
	"fmt"
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
	diallerSvc := newSocketDialer(s.cfg.GetConnectionPath(), s.cfg.GetUnixFileNameTemplate())

	err := diallerSvc.Prepare()
	if err != nil {
		return err
	}

	options := []originGRPC.DialOption{
		originGRPC.WithContextDialer(diallerSvc.DialCallback),
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithDisableRetry(),
		originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
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
