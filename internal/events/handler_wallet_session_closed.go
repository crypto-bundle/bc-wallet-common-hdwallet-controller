package events

import (
	"context"

	"github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/internal/entities"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	pbHdwallet "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
)

type sessionClosedHandler struct {
	walletCacheDataSvc mnemonicWalletsCacheStoreService
	walletDataSvc      mnemonicWalletsDataService

	txStmtManager transactionalStatementManager

	hdWalletSvc pbHdwallet.HdWalletApiClient
}

func (h *sessionClosedHandler) Process(ctx context.Context, event *pbApi.WalletSessionEvent) error {
	walletItem, sessionItem, err := h.walletCacheDataSvc.GetMnemonicWalletSessionInfoByUUID(ctx,
		event.WalletIdentifier.WalletUUID, event.SessionIdentifier.SessionUUID)
	if err != nil {
		return err
	}

	if walletItem != nil && sessionItem != nil {
		return h.process(ctx, walletItem, sessionItem)
	}

	// in case of missing cache data
	if err = h.txStmtManager.BeginTxWithRollbackOnError(ctx, func(txStmtCtx context.Context) error {
		wallet, clbErr := h.walletDataSvc.GetMnemonicWalletByUUID(ctx, event.WalletIdentifier.WalletUUID)
		if clbErr != nil {
			return clbErr
		}

		session, clbErr := h.walletDataSvc.GetWalletSessionByUUID(ctx, event.SessionIdentifier.SessionUUID)
		if clbErr != nil {
			return clbErr
		}

		walletItem = wallet
		sessionItem = session

		return nil
	}); err != nil {
		return err
	}

	if sessionItem == nil || walletItem == nil {
		return nil
	}

	return h.process(ctx, walletItem, sessionItem)
}

func (h *sessionClosedHandler) process(ctx context.Context,
	wallet *entities.MnemonicWallet,
	session *entities.MnemonicWalletSession,
) error {
	if session.IsSessionActive() {
		return nil
	}

	_, err := h.hdWalletSvc.UnLoadMnemonic(ctx, &pbHdwallet.UnLoadMnemonicRequest{
		MnemonicIdentity: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: wallet.UUID.String(),
			WalletHash: wallet.MnemonicHash,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func MakeEventSessionClosedHandler(walletCacheDataSvc mnemonicWalletsCacheStoreService,
	walletDataSvc mnemonicWalletsDataService,
	hdWalletSvc pbHdwallet.HdWalletApiClient,
	txStmtManager transactionalStatementManager,
) *sessionClosedHandler {
	return &sessionClosedHandler{
		walletCacheDataSvc: walletCacheDataSvc,
		walletDataSvc:      walletDataSvc,
		txStmtManager:      txStmtManager,
		hdWalletSvc:        hdWalletSvc,
	}
}
