package controller

import (
	"encoding/base64"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"google.golang.org/protobuf/types/known/anypb"
	"testing"
)

//func TestImportWallet(t *testing.T) {
//	options := []originGRPC.DialOption{
//		originGRPC.WithTransportCredentials(insecure.NewCredentials()),
//		// grpc.WithContextDialer(Dialer), // use it if u need load balancing via dns
//		originGRPC.WithBlock(),
//		originGRPC.WithKeepaliveParams(commonGRPCClient.DefaultKeepaliveClientOptions()),
//		originGRPC.WithChainUnaryInterceptor(commonGRPCClient.DefaultInterceptorsOptions()...),
//	}
//	grpcConn, err := originGRPC.Dial("localhost:8114", options...)
//	if err != nil {
//		t.Error(err)
//	}
//
//	client := NewHdWalletControllerApiClient(grpcConn)
//	ctx := context.Background()
//
//	resp, err := client.ImportWallet(ctx, &ImportWalletRequest{
//		MnemonicPhrase: []byte("vault:v1:6J6tXJH615/2yugCijTvWsxkk7CkN3S3KGgLgM+h1eYVtnWbzTs1SVC08P9ou0FLvQsW/qMT1yJyE184t68VIK0opB9y8nQQ0+hpgmUgZMXhMP8GpNf9FGZ6cCg2d0py2fo59gYS2yeSEkz5TnCKdUUHDaUIIbUTn70xPgqQlnGZLGtak1ap6Eji4KGUix2EkJTtO/ZicjdkPOhk3an7sM29sFw8VZYtRNY50SfajPcBrsBCkAhmLO0j"),
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if resp == nil {
//		t.Fail()
//	}
//}

func TestDisableWalletRequest(t *testing.T) {
	addr := &common.DerivationAddressIdentity{
		AccountIndex:  100,
		InternalIndex: 8,
		AddressIndex:  12,
	}

	anyData := anypb.Any{}
	err := anyData.MarshalFrom(addr)
	if err != nil {
		t.Fatal(err)
	}

	rawAddrStr := base64.StdEncoding.EncodeToString(anyData.Value)

	t.Logf("addr raw str: %s", rawAddrStr)
}

func TestDisableWalletRequest2(t *testing.T) {
	addr := &common.RangeUnitsList{
		RangeUnits: []*common.RangeRequestUnit{
			{
				AccountIndex:     5,
				InternalIndex:    5,
				AddressIndexFrom: 100,
				AddressIndexTo:   140,
			},
			{
				AccountIndex:     5,
				InternalIndex:    4,
				AddressIndexFrom: 1001,
				AddressIndexTo:   1041,
			},
			{
				AccountIndex:     5,
				InternalIndex:    3,
				AddressIndexFrom: 201,
				AddressIndexTo:   301,
			},
		},
	}

	anyData := anypb.Any{}
	err := anyData.MarshalFrom(addr)
	if err != nil {
		t.Fatal(err)
	}

	rawAddrStr := base64.StdEncoding.EncodeToString(anyData.Value)

	t.Logf("addr raw str: %s", rawAddrStr)
}
