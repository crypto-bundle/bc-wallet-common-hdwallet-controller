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
	l *zap.Logger
}

func (o *Generator) Generate(ctx context.Context) (string, error) {
	return o.generate(ctx)
}

func (o *Generator) generate(ctx context.Context) (string, error) {
	// 24 word mnemo for 256 bit entropy
	entropy, err := bip39.NewEntropy(256)
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

func NewMnemonicGenerator(logger *zap.Logger) (*Generator, error) {
	return &Generator{
		l: logger,
	}, nil
}
