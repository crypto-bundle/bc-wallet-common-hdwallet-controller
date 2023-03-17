package grpc

import (
	"go.uber.org/zap"
)

// grpcMarshaller is struct for marshall data from types.PublicData* to protobuf
type grpcMarshaller struct {
	logger *zap.Logger
}

func newGRPCMarshaller(logger *zap.Logger) *grpcMarshaller {
	return &grpcMarshaller{logger: logger}
}
