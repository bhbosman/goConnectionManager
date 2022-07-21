package goConnectionManager

import "github.com/cskr/pubsub"

type Helper struct {
	pubSub *pubsub.PubSub
}

type RefreshDataStart struct {
	SubscriptionName string
}

type RefreshDataStop struct {
	SubscriptionName string
}

type RefreshDataTo struct {
	SubscriptionName string
}

func (self *Helper) RefreshChannelName() string {
	return "ConnectionManagerRefresh"
}

func (self *Helper) Close() error {
	return nil
}

func (self *Helper) RefreshData(subscriptionName string) {
	self.pubSub.Pub(
		&RefreshDataTo{
			SubscriptionName: subscriptionName,
		},
		self.RefreshChannelName(),
	)

}

func (self *Helper) Pub(msg interface{}, topics ...string) {
	self.pubSub.Pub(msg, topics...)
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
