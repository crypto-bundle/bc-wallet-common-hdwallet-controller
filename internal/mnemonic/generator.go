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

package mnemonic

import (
	"context"
	"errors"

	"github.com/tyler-smith/go-bip39"
	"go.uber.org/zap"
)

var (
	ErrMnemonicIsInvalid = errors.New("mnemonic is invalid")
)

type Generator struct {
	l                            *zap.Logger
	defaultMnemonicSentenceCount uint8
	defaultMnemonicBitSize       int
}

func (o *Generator) Generate(ctx context.Context) (string, error) {
	return o.generate(ctx)
}

func (o *Generator) generate(_ context.Context) (string, error) {
	entropy, err := bip39.NewEntropy(o.defaultMnemonicBitSize)
	if err != nil {
		o.l.Error("error create entropy", zap.Error(err))
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		o.l.Error("error create new mnemonic", zap.Error(err))

		return "", err
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		o.l.Error("created mnemonic is invalid", zap.Error(ErrMnemonicIsInvalid))

		return "", ErrMnemonicIsInvalid
	}

	return mnemonic, nil
}

func NewMnemonicGenerator(logger *zap.Logger,
	defaultMnemonicWordCount uint8,
) *Generator {
	var defaultMnemonicBitSize int
	switch defaultMnemonicWordCount {
	case 18: // 18 word mnemo for 192 bit entropy
		defaultMnemonicBitSize = 192
	case 21: // 21 word mnemo for 224 bit entropy
		defaultMnemonicBitSize = 224
	case 24: // 24 word mnemo for 256 bit entropy
	default:
		defaultMnemonicBitSize = 256
	}

	return &Generator{
		l:                            logger,
		defaultMnemonicSentenceCount: defaultMnemonicWordCount,
		defaultMnemonicBitSize:       defaultMnemonicBitSize,
	}
}
