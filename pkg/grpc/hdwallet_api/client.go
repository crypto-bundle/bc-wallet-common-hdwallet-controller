package hdwallet_api

import (
	pbApi "bc-wallet-eth-hdwallet/pkg/grpc/hdwallet_api/proto"
	"context"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	cfg clientConfig

	client pbApi.HdWalletApiClient
}

// Init bcexplorer service
// nolint:revive // fixme (autofix)
func (s *Client) Init(ctx context.Context) error {
	options := DefaultDialOptions()
	msgSizeOptions := originGRPC.WithDefaultCallOptions(
		originGRPC.MaxCallRecvMsgSize(DefaultClientMaxReceiveMessageSize),
		originGRPC.MaxCallSendMsgSize(DefaultClientMaxSendMessageSize),
	)
	options = append(options, msgSizeOptions,
		originGRPC.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)

	grpcConn, err := originGRPC.Dial(s.cfg.GetHdWalletServerAddress(), options...)
	if err != nil {
		return err
	}

	s.client = pbApi.NewHdWalletApiClient(grpcConn)

	return nil
}

// Shutdown bcexplorer service
// nolint:revive // fixme (autofix)
func (s *Client) Shutdown(ctx context.Context) error {

	return nil
}

// GetEnabledWallets is function for getting address from bcexplorer
func (s *Client) GetEnabledWallets(ctx context.Context) (*pbApi.GetEnabledWalletsResponse, error) {
	request := &pbApi.GetEnabledWalletsRequest{}

	enabledWallets, err := s.client.GetEnabledWallets(ctx, request)
	if err != nil {
		return nil, err
	}

	return enabledWallets, nil
}

// GetDerivationAddress is function for getting address from bcexplorer
func (s *Client) GetDerivationAddress(ctx context.Context,
	walletUUID string,
	accountIndex uint32,
	internalIndex uint32,
	addressIndex uint32,
) (*pbApi.DerivationAddressResponse, error) {
	request := &pbApi.DerivationAddressRequest{
		AddressIdentity: &pbApi.DerivationAddressIdentity{
			AccountIndex:  accountIndex,
			InternalIndex: internalIndex,
			AddressIndex:  addressIndex,
		},
	}

	md := metadata.New(map[string]string{
		"wallet_uuid": walletUUID,
	})

	requestCtx := metadata.NewOutgoingContext(ctx, md)

	address, err := s.client.GetDerivationAddress(requestCtx, request)
	if err != nil {
		return nil, err
	}

	return address, nil
}
