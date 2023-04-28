package grpc_hdwallet

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"

	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	"go.uber.org/zap"
)

const (
	MethodNameUnLoadMnemonic = "UnLoadMnemonic"
)

type unLoadMnemonicHandler struct {
	l *zap.Logger
}

// nolint:funlen // fixme
func (h *unLoadMnemonicHandler) Handle(ctx context.Context,
	req *pbApi.UnLoadMnemonicRequest,
) (*pbApi.UnLoadMnemonicResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	return nil, nil
}

func MakeUnLoadMnemonicHandler(loggerEntry *zap.Logger,
) *unLoadMnemonicHandler {
	return &unLoadMnemonicHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameUnLoadMnemonic)),
	}
}
