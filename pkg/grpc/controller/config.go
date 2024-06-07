/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package controller

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ConfigHdWalletControllerApiClient struct {
	MaxReceiveMessageSize int `envconfig:"HDWALLET_CONTROLLER_CLIENT_RECEIVE_MESSAGE_SIZE" default:"262144"`
	MaxSendMessageSize    int `envconfig:"HDWALLET_CONTROLLER_CLIENT_SEND_MESSAGE_SIZE" default:"262144"`

	PowControlEnabled              bool `envconfig:"HDWALLET_CONTROLLER_CLIENT_POW_ENABLED" default:"true"`
	AccessTokenTokenControlEnabled bool `envconfig:"HDWALLET_CONTROLLER_CLIENT_AT_ENABLED" default:"true"`

	baseCfg baseConfigService

	serverPort uint
	serverHost string
	serverAddr string
}

func (c ConfigHdWalletControllerApiClient) GetServerPort() uint {
	return c.serverPort
}

func (c ConfigHdWalletControllerApiClient) GetServerHost() string {
	return c.serverHost
}

func (c ConfigHdWalletControllerApiClient) GetServerBindAddress() string {
	return c.serverAddr
}

func (c ConfigHdWalletControllerApiClient) GetMaxReceiveMessageSize() int {
	return c.MaxReceiveMessageSize
}

func (c ConfigHdWalletControllerApiClient) GetMaxSendMessageSize() int {
	return c.MaxSendMessageSize
}

func (c ConfigHdWalletControllerApiClient) IsPowShieldEnabled() bool {
	return c.PowControlEnabled
}

func (c ConfigHdWalletControllerApiClient) IsAccessTokenShieldEnabled() bool {
	return c.AccessTokenTokenControlEnabled
}

// Prepare variables to static configuration
func (c ConfigHdWalletControllerApiClient) Prepare() error {
	host := os.Getenv(fmt.Sprintf(serverAddressTemplate, c.baseCfg.GetNetworkName()))
	port := os.Getenv(fmt.Sprintf(serverPortTemplate, c.baseCfg.GetNetworkName()))

	srvPort, err := strconv.ParseUint(port, 10, 0)
	if err != nil {
		return nil
	}

	c.serverPort = uint(srvPort)
	c.serverHost = host

	c.serverAddr = fmt.Sprintf("%s:%d",
		strings.TrimRight(c.serverHost, ":"), c.serverPort)

	return nil
}

func (c ConfigHdWalletControllerApiClient) PrepareWith(cfgSvcList ...interface{}) error {
	for _, cfgSrv := range cfgSvcList {
		switch castedCfg := cfgSrv.(type) {
		case baseConfigService:
			c.baseCfg = castedCfg
		default:
			continue
		}
	}

	return nil
}
