package goConnectionManager

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type IPublishConnectionInformation interface {
	ConnectionInformationReceived(counters *model.PublishRxHandlerCounters) error
}

type IConnectionManager interface {
	ISendMessage.ISendMessage
	ISendMessage.IMultiSendMessage
	IPublishConnectionInformation
	model.IRegisterToConnectionManager
	IObtainConnectionManagerInformation
	ICommandsToConnectionManager
}
