package goConnectionManager

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/cskr/pubsub"
)

type Helper struct {
	pubSub *pubsub.PubSub
}

//type RefreshDataStart struct {
//	SubscriptionName string
//}

//type RefreshDataStop struct {
//	SubscriptionName string
//}

type RefreshDataTo struct {
	PubSubBag goCommsDefinitions.IPubSubBag
}

func (self *Helper) RefreshChannelName() string {
	return "ConnectionManagerRefresh"
}

func (self *Helper) Close() error {
	return nil
}

func (self *Helper) Pub(msg interface{}, topics ...string) bool {
	return self.pubSub.PubWithContext(msg, topics...)
}

func (self *Helper) PublishChannelName() string {
	return "ActiveConnectionStatus"
}

func NewHelper(
	pubSub *pubsub.PubSub,
) (IHelper, error) {
	return &Helper{
		pubSub: pubSub,
	}, nil
}
