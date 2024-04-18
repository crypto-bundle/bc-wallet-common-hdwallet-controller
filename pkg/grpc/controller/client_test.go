package controller

import (
	"context"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	commonGRPCClient "github.com/crypto-bundle/bc-wallet-common-lib-grpc/pkg/client"
	originGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestImportWallet(t *testing.T) {
	options := []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
		originGRPC.WithChainUnaryInterceptor(commonGRPCClient.DefaultInterceptorsOptions()...),
	}
	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
	if err != nil {
		t.Error(err)
	}

	client := NewHdWalletControllerApiClient(grpcConn)
	ctx := context.Background()

	resp, err := client.ImportWallet(ctx, &ImportWalletRequest{
		MnemonicPhrase: []byte("vault:v1:6J6tXJH615/2yugCijTvWsxkk7CkN3S3KGgLgM+h1eYVtnWbzTs1SVC08P9ou0FLvQsW/qMT1yJyE184t68VIK0opB9y8nQQ0+hpgmUgZMXhMP8GpNf9FGZ6cCg2d0py2fo59gYS2yeSEkz5TnCKdUUHDaUIIbUTn70xPgqQlnGZLGtak1ap6Eji4KGUix2EkJTtO/ZicjdkPOhk3an7sM29sFw8VZYtRNY50SfajPcBrsBCkAhmLO0j"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fail()
	}
}

func TestAddNewWalletRequest(t *testing.T) {
	options := []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
		originGRPC.WithChainUnaryInterceptor(commonGRPCClient.DefaultInterceptorsOptions()...),
	}
	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
	if err != nil {
		t.Error(err)
	}

	client := NewHdWalletControllerApiClient(grpcConn)
	ctx := context.Background()

	resp, err := client.AddNewWallet(ctx, &AddNewWalletRequest{})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fail()
	}
}

func TestEnableWalletRequest(t *testing.T) {
	options := []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
		originGRPC.WithChainUnaryInterceptor(commonGRPCClient.DefaultInterceptorsOptions()...),
	}
	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
	if err != nil {
		t.Error(err)
	}

	client := NewHdWalletControllerApiClient(grpcConn)
	ctx := context.Background()

	resp, err := client.EnableWallet(ctx, &EnableWalletRequest{
		WalletIdentity: &common.MnemonicWalletIdentity{
			WalletUUID: "8b10e0c1-4281-4c5f-aecd-d5ec3cf0e85e",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fail()
	}
}

func TestGetEnabledWalletsRequest(t *testing.T) {
	options := []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
		originGRPC.WithChainUnaryInterceptor(commonGRPCClient.DefaultInterceptorsOptions()...),
	}
	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
	if err != nil {
		t.Error(err)
	}

	client := NewHdWalletControllerApiClient(grpcConn)
	ctx := context.Background()

	resp, err := client.GetEnabledWallets(ctx, &GetEnabledWalletsRequest{})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fail()
	}
}

func TestGetWalletInfoRequest(t *testing.T) {
	options := []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
		originGRPC.WithChainUnaryInterceptor(commonGRPCClient.DefaultInterceptorsOptions()...),
	}
	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
	if err != nil {
		t.Error(err)
	}

	client := NewHdWalletControllerApiClient(grpcConn)
	ctx := context.Background()

	resp, err := client.GetWalletInfo(ctx, &GetWalletInfoRequest{
		WalletIdentity: &common.MnemonicWalletIdentity{
			WalletUUID: "8b10e0c1-4281-4c5f-aecd-d5ec3cf0e85e",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fail()
	}
}

func TestDisableWalletRequest(t *testing.T) {
	options := []originGRPC.DialOption{
		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
		originGRPC.WithBlock(),
		originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
		originGRPC.WithChainUnaryInterceptor(commonGRPCClient.DefaultInterceptorsOptions()...),
	}
	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
	if err != nil {
		t.Error(err)
	}

	client := NewHdWalletControllerApiClient(grpcConn)
	ctx := context.Background()

	resp, err := client.DisableWallet(ctx, &DisableWalletRequest{
		WalletIdentity: &common.MnemonicWalletIdentity{
			WalletUUID: "8b10e0c1-4281-4c5f-aecd-d5ec3cf0e85e",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fail()
	}
}
