/*
 *
 *
 * MIT-License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package events

import (
	"context"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	pbHdwallet "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"
	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonRedis "github.com/crypto-bundle/bc-wallet-common-lib-redis/pkg/redis"

	originRedis "github.com/go-redis/redis/v8"
	originNats "github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type watcher struct {
	logger *zap.Logger

	natsConn  *originNats.Conn
	redisConn *originRedis.Client

	natsSubQueueName  string
	redisSubQueueName string

	appInstanceIdentifier string
	workersCount          int
	bufferSize            int

	sessionStartedHandler  walletSessionStartedProcessorService
	sessionClosedHandler   walletSessionClosedProcessorService
	signReqPreparedHandler signRequestPreparedProcessorService
}

func (w *watcher) Run(ctx context.Context) error {
	redisSub := w.redisConn.PSubscribe(ctx, w.redisSubQueueName)
	redisDataChan := redisSub.Channel(originRedis.WithChannelSize(w.bufferSize))

	natsDataChan := make(chan *originNats.Msg, w.bufferSize)
	natsSub, err := w.natsConn.ChanSubscribe(w.redisSubQueueName, natsDataChan)
	if err != nil {
		return err
	}

	for i := 0; i != w.workersCount; i++ {
		go w.runRedisWatcher(ctx, redisSub, redisDataChan)

		go w.runNatsWatcher(ctx, natsSub, natsDataChan)
	}

	return nil
}

func (w *watcher) runRedisWatcher(ctx context.Context,
	sub *originRedis.PubSub,
	dataChan <-chan *originRedis.Message,
) {
	for {
		select {
		case <-ctx.Done():
			closeErr := sub.Close()
			if closeErr != nil {
				w.logger.Error("unable to close redis watcher channel", zap.Error(closeErr))
			}

			return

		case msg := <-dataChan:
			procErr := w.processMessage(context.Background(), []byte(msg.Payload))
			if procErr != nil {
				w.logger.Error("unable to process incoming currency event message",
					zap.Error(procErr))
			}
		}
	}
}

func (w *watcher) runNatsWatcher(ctx context.Context,
	sub *originNats.Subscription,
	dataChan chan *originNats.Msg,
) {
	for {
		select {
		case <-ctx.Done():
			closeErr := sub.Unsubscribe()
			if closeErr != nil {
				w.logger.Error("unable to unsubscribe from nats watcher channel", zap.Error(closeErr))
			}

			return

		case msg := <-dataChan:
			procErr := w.processMessage(context.Background(), msg.Data)
			if procErr != nil {
				w.logger.Error("unable to process incoming currency event message",
					zap.Error(procErr))
			}
		}
	}
}

func (w *watcher) processMessage(ctx context.Context, rawMsg []byte) error {
	pbEvent := &pbApi.Event{}
	err := proto.Unmarshal(rawMsg, pbEvent)
	if err != nil {
		return err
	}

	if pbEvent.AppInstanceIdentifier.UUID == w.appInstanceIdentifier {
		return nil
	}

	switch pbEvent.EventType {
	case pbApi.Event_EVENT_TYPE_SESSION:
		pbSessionEvent := &pbApi.WalletSessionEvent{}
		innerErr := proto.Unmarshal(pbEvent.Data, pbSessionEvent)
		if innerErr != nil {
			return innerErr
		}

		return w.processWalletSessionEvent(ctx, pbSessionEvent)

	case pbApi.Event_EVENT_TYPE_SIGN_REQUEST:
		pbSignRequestEvent := &pbApi.SignRequestEvent{}
		innerErr := proto.Unmarshal(pbEvent.Data, pbSignRequestEvent)
		if innerErr != nil {
			return innerErr
		}

		return w.processSignEvent(ctx, pbSignRequestEvent)
	default:

		return nil
	}
}

func (w *watcher) processWalletSessionEvent(ctx context.Context, event *pbApi.WalletSessionEvent) error {
	switch event.EventType {
	case pbApi.WalletSessionEvent_STARTED:
		return w.sessionStartedHandler.Process(ctx, event)
	case pbApi.WalletSessionEvent_CLOSED:
		return w.sessionClosedHandler.Process(ctx, event)
	default:
		return nil
	}
}

func (w *watcher) processSignEvent(ctx context.Context, event *pbApi.SignRequestEvent) error {
	switch event.EventType {
	case pbApi.SignRequestEvent_PREPARED:
		return w.signReqPreparedHandler.Process(ctx, event)
	default:
		return nil
	}
}

func NewEventWatcher(logger *zap.Logger,
	cfgSvc configService,
	redisConn *commonRedis.Connection,
	natsConn *commonNats.Connection,
	walletCacheDataSvc mnemonicWalletsCacheStoreService,
	walletDataSvc mnemonicWalletsDataService,
	signReqDataSvc signRequestDataService,
	hdWalletSvc pbHdwallet.HdWalletApiClient,
	txStmtManager transactionalStatementManager,
) *watcher {
	return &watcher{
		logger: logger,

		natsConn:  natsConn.GetConnection(),
		redisConn: redisConn.GetClient(),

		bufferSize:   cfgSvc.GetEventChannelBufferSize(),
		workersCount: cfgSvc.GetEventChannelWorkersCount(),

		redisSubQueueName: cfgSvc.GetEventChannelName(),
		natsSubQueueName:  cfgSvc.GetEventChannelName(),

		appInstanceIdentifier: cfgSvc.GetInstanceIdentifier().String(),

		sessionStartedHandler: MakeEventSessionStartedHandler(walletCacheDataSvc, walletDataSvc,
			hdWalletSvc, txStmtManager),
		sessionClosedHandler: MakeEventSessionClosedHandler(walletCacheDataSvc, walletDataSvc,
			hdWalletSvc, txStmtManager),
		signReqPreparedHandler: MakeEventSignReqPreparedHandler(walletCacheDataSvc, walletDataSvc,
			signReqDataSvc,
			hdWalletSvc, txStmtManager),
	}
}
