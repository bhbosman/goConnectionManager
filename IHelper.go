package goConnectionManager

import (
	"io"
)

type IHelper interface {
	PublishChannelName() string
	RefreshChannelName() string
	Pub(msg interface{}, topics ...string)
	RefreshData(subscriptionName string)
	io.Closer
}
