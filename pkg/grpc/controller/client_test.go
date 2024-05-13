package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"math"
	"math/big"
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

func TestFindSha256(t *testing.T) {
	req := &GetWalletInfoRequest{
		WalletIdentifier: &common.MnemonicWalletIdentity{
			WalletUUID: uuid.NewString(),
		},
	}

	reqRaw, err := proto.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	target := big.NewInt(1)
	target.Lsh(target, uint(256-8))

	target2 := big.NewInt(1)
	target2.Lsh(target2, uint(256-24))

	nonce := int64(0)

	//t.Logf("target: %d", target)
	//t.Logf("target2: %d", target2)

	var reqInt *big.Int = big.NewInt(0)

	for nonce != math.MaxInt64 {
		concatRaw := append(reqRaw[:], byte(nonce))

		hash := sha256.Sum256(concatRaw)
		reqInt.SetBytes(hash[:])

		t.Logf("\r%x: calc: %d, target: %d", hash, reqInt, target)

		if reqInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}

	t.Logf("addr raw str: %x", reqInt.Bytes())
}

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
				AddressIndexTo:   150,
			},
			{
				AccountIndex:     5,
				InternalIndex:    4,
				AddressIndexFrom: 100,
				AddressIndexTo:   150,
			},
			{
				AccountIndex:     5,
				InternalIndex:    3,
				AddressIndexFrom: 100,
				AddressIndexTo:   150,
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
