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

package config

import (
	"fmt"
	"strings"

	"github.com/crypto-bundle/bc-wallet-eth-hdwallet/internal/app"

	"github.com/kelseyhightower/envconfig"
)

// GrpcConfig for application
type GrpcConfig struct {
	// BindRaw port string, default "8080"
	BindRaw string `envconfig:"API_GRPC_PORT" default:"8080"`

	// ----------------------------
	// Calculated config parameters
	Bind string
}

func (c *GrpcConfig) GetBindPort() string {
	return c.Bind
}

// Prepare variables to static configuration
func (c *GrpcConfig) Prepare() error {
	err := envconfig.Process(app.APIConfigPrefix, c)
	if err != nil {
		return err
	}

	c.Bind = fmt.Sprintf(":%s", strings.TrimLeft(c.BindRaw, ":"))

	return err
}
