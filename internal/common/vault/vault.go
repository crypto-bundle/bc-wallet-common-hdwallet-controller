package vault

import (
	"context"
	"encoding/json"

	vaultApi "github.com/hashicorp/vault/api"
	userPassAuth "github.com/hashicorp/vault/api/auth/userpass"
)

const githubAuthPath = "auth/github/login"

type Vaulter interface {
	Encrypt(toEncrypt []byte) ([]byte, error)
	Decrypt(cipherBytes []byte) ([]byte, error)

	GetCredentialsBytes() (b []byte, err error)
	GetCredentialsByPathAndKey(path, field string) (string, error)
	GetCredentialsByPathAndKeys(path string, fields ...string) (map[string]string, error)
}

type service struct {
	client *vaultApi.Client
	cfg    *Config
}

// GetCredentialsBytes returns all secrets bytes from default path.
func (s *service) GetCredentialsBytes() (b []byte, err error) {
	secret, err := s.client.Logical().Read(s.cfg.DataPath)
	if err != nil {
		return nil, NewInternalError(ErrReadSecret, err)
	}
	if secret == nil {
		return nil, ErrEmptySecret
	}

	return json.Marshal(secret.Data["data"])
}

// GetCredentialsByPathAndKey returns secret by path and field.
func (s *service) GetCredentialsByPathAndKey(path, key string) (string, error) {
	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return "", NewInternalError(ErrReadSecret, err)
	}
	if secret == nil {
		return "", ErrEmptySecret
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", ErrCastSecret
	}

	keyVal, ok := data[key]
	if !ok {
		return "", ErrNotExistingKey.WithMsg(key)
	}

	keyString, ok := keyVal.(string)
	if !ok {
		return "", ErrKeyType.WithMsg(key)
	}

	return keyString, nil
}

// GetCredentialsByPathAndKeys returns fields sets from path.
func (s *service) GetCredentialsByPathAndKeys(path string, keys ...string) (map[string]string, error) {
	res := make(map[string]string, len(keys))

	secret, err := s.client.Logical().Read(path)
	if err != nil {
		return nil, NewInternalError(ErrReadSecret, err)
	}
	if secret == nil {
		return nil, ErrEmptySecret
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, ErrCastSecret
	}

	for _, k := range keys {
		keyVal, ok := data[k]
		if !ok {
			return res, ErrNotExistingKey.WithMsg(k)
		}

		keyString, ok := keyVal.(string)
		if !ok {
			return res, ErrKeyType.WithMsg(k)
		}

		res[k] = keyString
	}

	return res, nil
}

// NewClientByGithubToken initialize vault client with github_token authorization.
func NewClientByGithubToken(cfg *Config, ghToken string) (Vaulter, error) {
	clientOpts := vaultApi.DefaultConfig()
	clientOpts.Address = cfg.Address

	client, err := vaultApi.NewClient(clientOpts)
	if err != nil {
		return nil, NewInternalError(ErrVaultAPIClientInit, err)
	}

	secret, err := client.Logical().Write(githubAuthPath, map[string]interface{}{"token": ghToken})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, ErrEmptySecret
	}

	client.SetToken(secret.Auth.ClientToken)

	return &service{
		client: client,
		cfg:    cfg,
	}, nil
}

// NewClient initialize vault client with service account token authorization.
func NewClient(ctx context.Context, cfg *Config) (Vaulter, error) {
	// to have the fallback way for application config initialization only from envs.
	if cfg.IsEmpty() {
		return nil, nil
	}

	clientOpts := vaultApi.DefaultConfig()
	clientOpts.Address = cfg.Address

	client, err := vaultApi.NewClient(clientOpts)
	if err != nil {
		return nil, NewInternalError(ErrVaultAPIClientInit, err)
	}

	password := &userPassAuth.Password{
		FromFile:   "",
		FromEnv:    "",
		FromString: cfg.AuthPassword,
	}

	auth, err := userPassAuth.NewUserpassAuth(cfg.AuthUsername, password)
	if err != nil {
		return nil, NewInternalError(ErrUserNamePathAuthInit, err)
	}

	authInfo, err := client.Auth().Login(ctx, auth)
	if err != nil {
		return nil, NewInternalError(ErrUserNamePathLogin, err)
	}
	if authInfo == nil {
		return nil, NewInternalError(ErrNotExistingAuthInfo, err)
	}

	client.SetToken(authInfo.Auth.ClientToken)

	return &service{
		client: client,
		cfg:    cfg,
	}, nil
}
