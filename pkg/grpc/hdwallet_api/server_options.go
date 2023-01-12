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
	"math"
	"time"

	originGRPC "google.golang.org/grpc"
	grpcKeepalive "google.golang.org/grpc/keepalive"
)

const (
	DefaultServerMaxReceiveMessageSize = math.MaxInt32
	DefaultServerMaxSendMessageSize    = 1024 * 1024 * 24
)

func DefaultEnforcementServerOptions() grpcKeepalive.EnforcementPolicy {
	return grpcKeepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}
}

func DefaultKeepaliveServerOptions() grpcKeepalive.ServerParameters {
	return grpcKeepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second,
		MaxConnectionAge:      30 * time.Second,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  5 * time.Second,
		Timeout:               1 * time.Second,
	}
}

func DefaultServeOptions() []originGRPC.ServerOption {
	return []originGRPC.ServerOption{
		originGRPC.KeepaliveEnforcementPolicy(DefaultEnforcementServerOptions()),
		originGRPC.KeepaliveParams(DefaultKeepaliveServerOptions()),
	}
}
