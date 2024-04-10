package manager

import (
	"context"

	commonGRPCClient "github.com/crypto-bundle/bc-wallet-common-lib-grpc/pkg/client"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	cfg HdWalletControllerGRPCClientConfig

	grpcClientOptions []originGRPC.DialOption

	grpcConn *originGRPC.ClientConn
	client   HdWalletManagerApiClient
}

// Init bc-wallet-tron-hdwallet GRPC-client service
// nolint:revive // fixme (autofix)
func (s *Client) Init(ctx context.Context) error {
	options := []originGRPC.DialOption{
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
		originGRPC.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)

	s.grpcClientOptions = options

	return nil
}

func (s *Client) Dial(ctx context.Context) error {
	grpcConn, err := originGRPC.Dial("/tmp/grpc.sock", s.grpcClientOptions...)
	if err != nil {
		return err
	}
	s.grpcConn = grpcConn

	s.client = NewHdWalletManagerApiClient(grpcConn)

	return nil
}

// Shutdown grpc hdwallet-client service
func (s *Client) Shutdown(ctx context.Context) error {
	err := s.grpcConn.Close()
	if err != nil {
		return err
	}

	return nil
}

// nolint:revive // fixme
func NewClient(ctx context.Context,
	cfg HdWalletControllerGRPCClientConfig,
) (*Client, error) {
	srv := &Client{
		cfg: cfg,
	}

	return srv, nil
}
