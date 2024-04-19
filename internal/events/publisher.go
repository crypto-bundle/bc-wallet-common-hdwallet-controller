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
