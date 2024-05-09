package hdwallet

type hdWalletClientConfig interface {
	GetConnectionPath() string
	GetUnixFileNameTemplate() string
}
