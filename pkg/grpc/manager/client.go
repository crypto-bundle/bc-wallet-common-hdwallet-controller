package manager

import (
	"context"
	"errors"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-manager/pkg/grpc/common"

	commonGRPCClient "github.com/crypto-bundle/bc-wallet-common-lib-grpc/pkg/client"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	ErrUnableDecodeGrpcErrorStatus = errors.New("unable to decode grpc error status")
)

type Client struct {
	cfg HdWalletGRPCClientConfig

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
	grpcConn, err := originGRPC.Dial(s.cfg.GetHdWalletApiHost(), s.grpcClientOptions...)
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

// AddNewWallet is function for add new wallet
func (s *Client) AddNewWallet(ctx context.Context,
	title, purpose string,
	strategy uint32,
) (*AddNewWalletResponse, error) {
	request := &AddNewWalletRequest{
		Title:    title,
		Purpose:  purpose,
		Strategy: pbCommon.WalletMakerStrategy(strategy),
	}

	walletResp, err := s.client.AddNewWallet(ctx, request)
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			return nil, ErrUnableDecodeGrpcErrorStatus
		}

		switch grpcStatus.Code() {
		case codes.NotFound:
			return nil, nil
		default:
			return nil, err
		}
	}

	return walletResp, nil
}

// GetEnabledWallets is function for getting address from bc-wallet-tron-hdwallet
func (s *Client) GetEnabledWallets(ctx context.Context) (*GetEnabledWalletsResponse, error) {
	enabledWallets, err := s.client.GetEnabledWallets(ctx, &GetEnabledWalletsRequest{})
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			return nil, ErrUnableDecodeGrpcErrorStatus
		}

		switch grpcStatus.Code() {
		case codes.NotFound:
			return nil, nil
		default:
			return nil, err
		}
	}

	return enabledWallets, nil
}

// SignTransaction is function for sign tron preparedTransaction by private key via hdwalelt srv
func (s *Client) SignTransaction(ctx context.Context,
	walletUUID string,
	mnemonicWalletUUID string,
	accountIndex, internalIndex, addressIndex uint32,
	tronCreatedTxData []byte,
) (*SignTransactionResponse, error) {
	signReq := &SignTransactionRequest{
		WalletUUID:         walletUUID,
		MnemonicWalletUUID: mnemonicWalletUUID,
		AddressIdentity: &pbCommon.DerivationAddressIdentity{
			AccountIndex:  accountIndex,
			InternalIndex: internalIndex,
			AddressIndex:  addressIndex,
		},
		CreatedTxData: tronCreatedTxData,
	}

	signResp, err := s.client.SignTransaction(ctx, signReq)
	if err != nil {
		return nil, err
	}

	return signResp, nil
}

// GetDerivationAddress is function for getting address from bc-wallet-tron-hdwallet
func (s *Client) GetDerivationAddress(ctx context.Context,
	walletUUID string,
	mnemonicWalletUUID string,
	accountIndex uint32,
	internalIndex uint32,
	addressIndex uint32,
) (*DerivationAddressResponse, error) {
	request := &DerivationAddressRequest{
		WalletIdentity: &pbCommon.WalletIdentity{
			WalletUUID: walletUUID,
		},
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: mnemonicWalletUUID,
		},
		AddressIdentity: &pbCommon.DerivationAddressIdentity{
			AccountIndex:  accountIndex,
			InternalIndex: internalIndex,
			AddressIndex:  addressIndex,
		},
	}

	address, err := s.client.GetDerivationAddress(ctx, request)
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			return nil, ErrUnableDecodeGrpcErrorStatus
		}

		switch grpcStatus.Code() {
		case codes.NotFound:
			return nil, nil
		default:
			return nil, err
		}
	}

	return address, nil
}

// GetDerivationAddressByRange is function for getting address from bc-wallet-tron-hdwallet
func (s *Client) GetDerivationAddressByRange(ctx context.Context,
	walletUUID string,
	mnemonicWalletUUID string,
	ranges []*pbCommon.RangeRequestUnit,
) (*DerivationAddressByRangeResponse, error) {
	request := &DerivationAddressByRangeRequest{
		WalletIdentity: &pbCommon.WalletIdentity{
			WalletUUID: walletUUID,
		},
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: mnemonicWalletUUID,
		},
		Ranges: ranges,
	}

	address, err := s.client.GetDerivationAddressByRange(ctx, request)
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			return nil, ErrUnableDecodeGrpcErrorStatus
		}

		switch grpcStatus.Code() {
		case codes.NotFound:
			return nil, nil
		default:
			return nil, err
		}
	}

	return address, nil
}

// GetWalletInfo is function for getting full wallet info from bc-wallet-tron-hdwallet
func (s *Client) GetWalletInfo(ctx context.Context,
	walletUUID string,
) (*GetWalletInfoResponse, error) {
	request := &GetWalletInfoRequest{
		WalletIdentity: &pbCommon.WalletIdentity{
			WalletUUID: walletUUID,
		},
	}

	walletInfo, err := s.client.GetWalletInfo(ctx, request)
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			return nil, ErrUnableDecodeGrpcErrorStatus
		}

		switch grpcStatus.Code() {
		case codes.NotFound:
			return nil, nil
		default:
			return nil, err
		}
	}

	return walletInfo, nil
}

// nolint:revive // fixme
func NewClient(ctx context.Context,
	cfg HdWalletGRPCClientConfig,
) (*Client, error) {
	srv := &Client{
		cfg: cfg,
	}

	return srv, nil
}
