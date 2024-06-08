/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package events

import (
	"context"
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/controller"
	commonNats "github.com/crypto-bundle/bc-wallet-common-lib-nats-queue/pkg/nats"
	commonRedis "github.com/crypto-bundle/bc-wallet-common-lib-redis/pkg/redis"
	originRedis "github.com/go-redis/redis/v8"
	originNats "github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type publisher struct {
	appInstanceIdentifier *pbApi.AppInstanceIdentity

	natsConn  *originNats.Conn
	redisConn *originRedis.Client

	redisPubQueueName string
	natsPubQueueName  string
}

func (p *publisher) SendSessionStartEvent(ctx context.Context,
	walletUUID string,
	sessionUUID string,
) error {
	event := &pbApi.WalletSessionEvent{
		EventType: pbApi.WalletSessionEvent_STARTED,
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletUUID,
			WalletHash: "",
		},
		SessionIdentifier: &pbApi.WalletSessionIdentity{SessionUUID: sessionUUID},
	}

	err := p.sendEvent(ctx, pbApi.Event_EVENT_TYPE_SESSION, event)
	if err != nil {
		return err
	}

	return nil
}

func (p *publisher) SendSessionClosedEvent(ctx context.Context,
	walletUUID string,
	sessionUUID string,
) error {
	event := &pbApi.WalletSessionEvent{
		EventType: pbApi.WalletSessionEvent_CLOSED,
		WalletIdentifier: &pbCommon.MnemonicWalletIdentity{
			WalletUUID: walletUUID,
			WalletHash: "",
		},
		SessionIdentifier: &pbApi.WalletSessionIdentity{SessionUUID: sessionUUID},
	}

	err := p.sendEvent(ctx, pbApi.Event_EVENT_TYPE_SESSION, event)
	if err != nil {
		return err
	}

	return nil
}

func (p *publisher) SendSignPreparedEvent(ctx context.Context,
	signReqUUID string,
) error {
	event := &pbApi.SignRequestEvent{
		EventType:             pbApi.SignRequestEvent_PREPARED,
		SignRequestIdentifier: &pbApi.SignRequestIdentity{UUID: signReqUUID},
	}

	err := p.sendEvent(ctx, pbApi.Event_EVENT_TYPE_SIGN_REQUEST, event)
	if err != nil {
		return err
	}

	return nil
}

func (p *publisher) sendEvent(ctx context.Context,
	eventType pbApi.Event_Type,
	eventData proto.Message,
) error {
	rawData, err := proto.Marshal(eventData)
	if err != nil {
		return err
	}

	msg := &pbApi.Event{
		EventType:             eventType,
		AppInstanceIdentifier: p.appInstanceIdentifier,
		Data:                  rawData,
	}
	eventRaw, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = p.redisConn.Publish(ctx, p.redisPubQueueName, eventRaw).Result()
	if err == nil {
		return nil
	}

	err = p.natsConn.Publish(p.natsPubQueueName, eventRaw)
	if err == nil {
		return nil
	}

	return err
}

func NewEventsBroadcaster(cfgSvc configService,
	natsConn *commonNats.Connection,
	redisConn *commonRedis.Connection,
) *publisher {
	return &publisher{
		natsConn:  natsConn.GetConnection(),
		redisConn: redisConn.GetClient(),

		appInstanceIdentifier: &pbApi.AppInstanceIdentity{UUID: cfgSvc.GetInstanceIdentifier().String()},

		redisPubQueueName: cfgSvc.GetEventChannelName(),
		natsPubQueueName:  cfgSvc.GetEventChannelName(),
	}
}
