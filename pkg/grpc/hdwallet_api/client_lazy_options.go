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
