package grpc

import (
	"go.uber.org/zap"
	"sync"
)

// grpcMarshaller is struct for marshall data from types.PublicData* to protobuf
type grpcMarshaller struct {
	logger     *zap.Logger
	pbAddrPool *sync.Pool
}

func newGRPCMarshaller(logger *zap.Logger, pbAddrPool *sync.Pool) *grpcMarshaller {
	return &grpcMarshaller{logger: logger, pbAddrPool: pbAddrPool}
}
