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
		fallthrough
	default:
		defaultMnemonicBitSize = 256
	}

	return &Generator{
		l:                            logger,
		defaultMnemonicSentenceCount: defaultMnemonicWordCount,
		defaultMnemonicBitSize:       defaultMnemonicBitSize,
	}
}
