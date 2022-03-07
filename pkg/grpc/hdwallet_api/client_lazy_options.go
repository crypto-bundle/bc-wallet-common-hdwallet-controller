package hdwallet_api

import (
	"context"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

const connectTimeout = 10 * time.Second

type LazyConnection struct {
	connection *grpc.ClientConn

	address string
	options []grpc.DialOption
}

func NewLazyConnection(address string, additionalOptions ...grpc.DialOption) *LazyConnection {
	options := DefaultDialOptions()
	options = append(
		options,
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),
		),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
	)
	options = append(options, additionalOptions...)

	return &LazyConnection{
		address: address,
		options: options,
	}
}

func (d *LazyConnection) connect(ctx context.Context) error {
	ctx, finish := context.WithTimeout(ctx, connectTimeout)
	defer finish()

	conn, err := grpc.DialContext(ctx, d.address, d.options...)
	if err != nil {
		return err
	}
	d.connection = conn
	return nil
}

func (d *LazyConnection) Invoke(ctx context.Context, method string, args, reply interface{},
	opts ...grpc.CallOption) error {
	if d.connection == nil {
		err := d.connect(ctx)
		if err != nil {
			return err
		}
	}
	return d.connection.Invoke(ctx, method, args, reply, opts...)
}

func (d *LazyConnection) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string,
	opts ...grpc.CallOption) (grpc.ClientStream, error) {
	if d.connection == nil {
		err := d.connect(ctx)
		if err != nil {
			return nil, err
		}
	}
	return d.connection.NewStream(ctx, desc, method, opts...)
}

func (d *LazyConnection) Close() {
	if d.connection != nil {
		_ = d.connection.Close()
	}
}
