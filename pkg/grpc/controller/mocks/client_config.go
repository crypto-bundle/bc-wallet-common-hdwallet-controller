/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
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
 *
 */

package mocks

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type controllerApiClientConfig struct {
	serverPort uint
	serverHost string

	network string
}

func (c controllerApiClientConfig) GetServerPort() uint {
	return c.serverPort
}

func (c controllerApiClientConfig) GetServerHost() string {
	return c.serverHost
}

func (c controllerApiClientConfig) GetServerBindAddress() string {
	return fmt.Sprintf("%s:%d",
		strings.TrimRight(c.serverHost, ":"), c.serverPort)
}

func (c controllerApiClientConfig) GetMaxReceiveMessageSize() int {
	return 1024 * 1024 * 9
}

func (c controllerApiClientConfig) GetMaxSendMessageSize() int {
	return 1024 * 1024 * 9
}

func (c controllerApiClientConfig) IsPowShieldEnabled() bool {
	return true
}

func (c controllerApiClientConfig) IsAccessTokenShieldEnabled() bool {
	return true
}

func (c controllerApiClientConfig) GetProviderName() string {
	return "crypto-bundle"
}
func (c controllerApiClientConfig) GetNetworkName() string {
	return c.network
}

type managerApiClientConfig struct {
	controllerApiClientConfig `json:"-"`

	RootToken string `json:"root_token"`
}

func (c managerApiClientConfig) GetRootToken() string {
	return c.RootToken
}

type walletApiClientConfig struct {
	controllerApiClientConfig
}

func NewManagerApiClientConfig(serverHost string,
	serverPort uint,
	network string,
	secretDataFilePath string,
) (*managerApiClientConfig, error) {
	fileRawData, err := os.ReadFile(secretDataFilePath)
	if err != nil {
		return nil, err
	}

	result := &managerApiClientConfig{
		controllerApiClientConfig: controllerApiClientConfig{
			serverPort: serverPort,
			serverHost: serverHost,
			network:    network,
		},
	}
	err = json.Unmarshal(fileRawData, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewWalletApiClientConfig(serverHost string,
	serverPort uint,
	network string,
) walletApiClientConfig {
	return walletApiClientConfig{
		controllerApiClientConfig: controllerApiClientConfig{
			serverPort: serverPort,
			serverHost: serverHost,
			network:    network,
		},
	}
}
