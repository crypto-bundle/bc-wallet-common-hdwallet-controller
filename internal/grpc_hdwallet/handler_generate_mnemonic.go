package grpc_hdwallet

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-tron-hdwallet/internal/app"

	tracer "github.com/crypto-bundle/bc-wallet-common-lib-tracer/pkg/tracer/opentracing"
	pbApi "github.com/crypto-bundle/bc-wallet-tron-hdwallet/pkg/grpc/hdwallet_api/proto"

	"go.uber.org/zap"
)

const (
	MethodNameGenerateMnemonic = "GenerateMnemonic"
)

type generateMnemonicHandler struct {
	l *zap.Logger
}

// nolint:funlen // fixme
func (h *generateMnemonicHandler) Handle(ctx context.Context,
	req *pbApi.GenerateMnemonicRequest,
) (*pbApi.GenerateMnemonicResponse, error) {
	var err error
	_, span, finish := tracer.Trace(ctx)

	defer func() { finish(err) }()

	span.SetTag(app.BlockChainNameTag, app.BlockChainName)

	return nil, nil
}

func MakeGenerateMnemonicHandler(loggerEntry *zap.Logger,
) *generateMnemonicHandler {
	return &generateMnemonicHandler{
		l: loggerEntry.With(zap.String(MethodNameTag, MethodNameGenerateMnemonic)),
	}
}
