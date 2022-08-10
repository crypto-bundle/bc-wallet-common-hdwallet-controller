package wallet

import "errors"

var (
	ErrWrongMnemonicHash    = errors.New("wrong mnemonic hash")
	ErrPassedWalletNotFound = errors.New("passed wallet not found")
)
