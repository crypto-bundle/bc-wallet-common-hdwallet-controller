package grpc

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MethodSignTransaction = "SignTransaction"
)

type SignTransactionHandler struct {
	l             *zap.Logger
	walletSrv     walletManagerService
	marshallerSrv marshallerService
}

// nolint:funlen // fixme
func (h *SignTransactionHandler) Handle(ctx context.Context,
	req *pbApi.SignTransactionRequest,
) (*pbApi.SignTransactionResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	walletUUID, err := uuid.Parse(req.WalletUUID)
	if err != nil {
		return nil, err
	}

	mnemonicWalletUUID, err := uuid.Parse(req.MnemonicWalletUUID)
	if err != nil {
		return nil, err
	}

	walletPubData, err := h.walletSrv.GetWalletByUUID(ctx, walletUUID)
	if err != nil {
		h.l.Error("unable get wallet", zap.Error(err))

		return nil, status.Error(codes.Internal, "something went wrong")
	}
	if walletPubData == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	mnemoWalletData, isExists := walletPubData.MnemonicWalletsByUUID[mnemonicWalletUUID]
	if !isExists {
		return nil, status.Error(codes.NotFound, "mnemonic wallet not found")
	}

	signedTxData, err := h.walletSrv.SignTransaction(ctx, walletPubData.UUID, mnemoWalletData.UUID,
		req.AddressIdentity.AccountIndex, req.AddressIdentity.InternalIndex, req.AddressIdentity.AddressIndex,
		req.CreatedTxData)
	if err != nil {
		h.l.Error("unable to sign transaction", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	response, err := h.marshallerSrv.MarshallSignTransaction(signedTxData)
	if err != nil {
		h.l.Error("unable to marshall signed transaction data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

func MakeSignTransactionsHandler(loggerEntry *zap.Logger,
	walletSrv walletManagerService,
	marshallerSrv marshallerService,
) *SignTransactionHandler {
	return &SignTransactionHandler{
		l:             loggerEntry.With(zap.String(MethodNameTag, MethodSignTransaction)),
		walletSrv:     walletSrv,
		marshallerSrv: marshallerSrv,
	}
}
