package wallet_manager

import "errors"

var (
	ErrWrongMnemonicHash                 = errors.New("wrong mnemonic hash")
	ErrWalletPoolIsNotEmpty              = errors.New("wallet pool is not empty - unable to overwrite it")
	ErrPassedWalletPoolUnitIsEmpty       = errors.New("passed wallets pool units list is empty")
	ErrPassedWalletNotFound              = errors.New("passed wallet not found")
	ErrPassedMnemonicWalletAlreadyExists = errors.New("passed mnemonic wallet already exists")
	ErrMnemonicAlreadySet                = errors.New("mnemonic unit already set in unit")
	ErrPassedWalletAlreadyExists         = errors.New("passed wallet already exists")
	ErrPassedMnemonicWalletNotFound      = errors.New("passed mnemonic wallet not found")
	ErrMethodUnimplemented               = errors.New("called service method unimplemented")
)
