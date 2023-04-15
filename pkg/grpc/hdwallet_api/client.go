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
	"errors"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"
	tronCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"

	commonGRPCClient "github.com/crypto-bundle/bc-wallet-common-lib-grpc/pkg/client"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnableDecodeGrpcErrorStatus = errors.New("unable to decode grpc error status")
)

type Client struct {
	cfg clientConfigService

	grpcClientOptions []originGRPC.DialOption

	grpcConn *originGRPC.ClientConn
	client   pbApi.HdWalletApiClient
}

// Init bc-wallet-tron-hdwallet GRPC-client service
// nolint:revive // fixme (autofix)
func (s *Client) Init(ctx context.Context) error {
	options := commonGRPCClient.DefaultDialOptions()
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
	grpcConn, err := originGRPC.Dial(s.cfg.GetHdWalletServerAddress(), s.grpcClientOptions...)
	if err != nil {
		return err
	}
	s.grpcConn = grpcConn

	s.client = pbApi.NewHdWalletApiClient(grpcConn)

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

// GetEnabledWallets is function for getting address from bc-wallet-tron-hdwallet
func (s *Client) GetEnabledWallets(ctx context.Context) (*pbApi.GetEnabledWalletsResponse, error) {
	enabledWallets, err := s.client.GetEnabledWallets(ctx, &pbApi.GetEnabledWalletsRequest{})
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
	tronCreatedTx *tronCore.Transaction,
) (*pbApi.SignTransactionResponse, error) {
	rawData, err := proto.Marshal(tronCreatedTx)
	if err != nil {
		return nil, err
	}

	signReq := &pbApi.SignTransactionRequest{
		WalletUUID:         walletUUID,
		MnemonicWalletUUID: mnemonicWalletUUID,
		AddressIdentity: &pbApi.DerivationAddressIdentity{
			AccountIndex:  accountIndex,
			InternalIndex: internalIndex,
			AddressIndex:  addressIndex,
		},
		CreatedTxData: rawData,
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
) (*pbApi.DerivationAddressResponse, error) {
	request := &pbApi.DerivationAddressRequest{
		WalletIdentity: &pbApi.WalletIdentity{
			WalletUUID: walletUUID,
		},
		MnemonicIdentity: &pbApi.MnemonicWalletIdentity{
			WalletUUID: mnemonicWalletUUID,
		},
		AddressIdentity: &pbApi.DerivationAddressIdentity{
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

// nolint:revive // fixme
func NewClient(ctx context.Context,
	cfg clientConfigService,
) (*Client, error) {
	srv := &Client{
		cfg: cfg,
	}

	return srv, nil
}
