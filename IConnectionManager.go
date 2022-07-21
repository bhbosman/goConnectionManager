package goConnectionManager

import (
	"github.com/bhbosman/gocommon/model"
)

type IPublishConnectionInformation interface {
	ConnectionInformationReceived(counters *model.PublishRxHandlerCounters) error
}

type IConnectionManager interface {
	IPublishConnectionInformation
	model.IRegisterToConnectionManager
	IObtainConnectionManagerInformation
	ICommandsToConnectionManager
}
