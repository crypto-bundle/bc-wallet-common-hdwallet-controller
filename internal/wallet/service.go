package wallet

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/entities"
	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/hdwallet"
	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/mnemonic"
	"github.com/cryptowize-tech/bc-wallet-eth-hdwallet/internal/wallet/repository"

	"github.com/cryptowize-tech/bc-wallet-common/pkg/postgres"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger         *zap.Logger
	cfg            config
	repo           repo
	mnemoGenerator mnemonicGenerator

	// hotWallets - Ethereum hot wallets
	hotWallets map[string]walleter
	// hotHdWallets - hot hdwallets
	hotHdWallets map[string]hdWalleter
}

func (s *Service) Init(ctx context.Context) error {
	wallets, err := s.repo.GetAllEnabledMnemonicWallets(ctx)
	if err != nil {
		return err
	}

	for i, _ := range wallets {
		hdWallet, creatErr := hdwallet.NewFromString(string(wallets[i].EncryptedData))
		if creatErr != nil {
			return creatErr
		}

		ethWallet, walletErr := hdWallet.NewEthWallet(uint32(0), uint32(0), uint32(0))
		if walletErr != nil {
			return walletErr
		}

		s.hotWallets[wallets[i].UUID.String()] = ethWallet
		s.hotHdWallets[wallets[i].UUID.String()] = hdWallet
	}

	return nil
}

func (s *Service) Shutdown(ctx context.Context) error {

	return nil
}

func (s *Service) GetMnemonicWallet(ctx context.Context,
	walletConfig mnemonicWalletConfig,
) (*entities.MnemonicWallet, error) {
	wallet, err := s.repo.GetMnemonicWalletByHash(ctx, walletConfig.GetMnemonicWalletHash())
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *Service) GetEnabledWalletsUUID(ctx context.Context) ([]string, error) {
	wallets, err := s.repo.GetAllEnabledNonHotMnemonicWallets(ctx)
	if err != nil {
		return nil, err
	}

	uuids := make([]string, len(wallets))

	for i, _ := range wallets {
		uuids[i] = wallets[i].UUID.String()
	}

	return uuids, nil
}

func (s *Service) GetAddressByPath(ctx context.Context,
	walletUUID string,
	account, change, index uint32,
) (string, error) {
	ethWallet, err := s.hotHdWallets[walletUUID].NewEthWallet(account, change, index)
	if err != nil {
		return "", err
	}

	ethAddress, err := ethWallet.GetAddress()
	if err != nil {
		return "", err
	}

	return ethAddress, nil
}

func (s *Service) CreateNewMnemonicWallet(ctx context.Context,
	title string,
	purpose string,
	isHot bool,
) (*entities.MnemonicWallet, error) {
	newWalletMnemonic, err := s.mnemoGenerator.Generate(ctx)
	if err != nil {
		return nil, err
	}

	walletEntity := &entities.MnemonicWallet{
		Title:         title,
		UUID:          uuid.New(),
		Hash:          fmt.Sprintf("%x", sha256.Sum256([]byte(newWalletMnemonic))),
		IsHotWallet:   isHot,
		Purpose:       purpose,
		EncryptedData: []byte(newWalletMnemonic),
	}

	walletEntity, err = s.repo.AddNewMnemonicWallet(ctx, walletEntity)
	if err != nil {
		return nil, err
	}

	if isHot {
		hdWallet, creatErr := hdwallet.NewFromString(string(walletEntity.EncryptedData))
		if creatErr != nil {
			return nil, creatErr
		}

		ethWallet, walletErr := hdWallet.NewEthWallet(uint32(0), uint32(0), uint32(0))
		if walletErr != nil {
			return nil, walletErr
		}

		uuidStr := walletEntity.UUID.String()

		s.hotWallets[uuidStr] = ethWallet
		s.hotHdWallets[uuidStr] = hdWallet
	}

	return walletEntity, nil
}

func New(logger *zap.Logger,
	cfg config,
	pgConn *postgres.Connection,
) (*Service, error) {
	pgRepos, err := repository.NewPostgresStore(logger, pgConn)
	if err != nil {
		return nil, err
	}

	mnemoGenerator, err := mnemonic.NewMnemonicGenerator(logger)
	if err != nil {
		return nil, err
	}

	return &Service{
		logger:         logger,
		cfg:            cfg,
		repo:           pgRepos,
		mnemoGenerator: mnemoGenerator,

		hotWallets:   make(map[string]walleter, 0),
		hotHdWallets: make(map[string]hdWalleter, 0),
	}, nil
}
