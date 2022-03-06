package config

type vaulter interface {
	GetCredentialsBytes() (b []byte, err error)
}
