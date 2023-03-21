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

import "errors"

var (
	ErrWrongMnemonicWordsCount = errors.New("wrong mnemonic words count.allowed: 15, 18, 21, 24")
	ErrBaseAppConfigIsEmpty    = errors.New("base app config is missing or empty")
	ErrMinWordsCount           = errors.New("minimal words count for environment is 21")
)

type MnemonicConfig struct {
	MnemonicWordCount uint8 `envconfig:"HDWALLET_WORDS_COUNT" default:"24"`

	baseAppCfgSrv baseConfigService
}

func (c *MnemonicConfig) GetDefaultMnemonicWordsCount() uint8 {
	return c.MnemonicWordCount
}

//nolint:funlen // its ok
func (c *MnemonicConfig) Prepare() error {
	if c.baseAppCfgSrv == nil {
		return ErrBaseAppConfigIsEmpty
	}

	if c.baseAppCfgSrv.IsProd() && c.MnemonicWordCount <= 18 {
		return ErrMinWordsCount
	}

	switch c.MnemonicWordCount {
	case 15, 18, 21, 24:
		return nil
	default:
		return ErrWrongMnemonicWordsCount
	}
}

//nolint:funlen // its ok
func (c *MnemonicConfig) PrepareWith(dependentCfgList ...interface{}) error {
	for _, cfgSrv := range dependentCfgList {
		switch castedCfg := cfgSrv.(type) {
		case baseConfigService:
			c.baseAppCfgSrv = castedCfg
		default:
			continue
		}
	}

	return nil
}
