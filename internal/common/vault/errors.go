package vault

import (
	"fmt"
)

const (
	ErrVaultAPIClientInit InternalError = "unable to initialize vault client"
	ErrK8sAuthInit        InternalError = "unable to initialize kubernetes auth method"
	ErrK8sLogin           InternalError = "unable to log in with kubernetes auth"

	ErrUserNamePathAuthInit InternalError = "unable to initialize username and pass auth method"
	ErrUserNamePathLogin    InternalError = "unable to log in with username and path auth"

	ErrNotExistingAuthInfo InternalError = "no auth info was returned after login"
	ErrReadSecret          InternalError = "unable to read secret"
	ErrEmptySecret         InternalError = "unable to get secret"
	ErrCastSecret          InternalError = "secret casting error"
	ErrNotExistingKey      InternalError = "missed key in secret"
	ErrKeyType             InternalError = "unexpected key type in secret"
	ErrTransitSecretFormat InternalError = "unexpected format for transit secret"
)

type InternalError string

func (e InternalError) Error() string {
	return string(e)
}

// WithMsg returns internal error with additional text.
func (e InternalError) WithMsg(msg string) error {
	return fmt.Errorf("%s: %s", e.Error(), msg)
}

// NewInternalError returns new error.
func NewInternalError(vErr InternalError, err error) error {
	return fmt.Errorf("%s: %q", vErr, err)
}
