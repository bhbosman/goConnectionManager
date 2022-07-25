package goConnectionManager

import (
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type IData interface {
	IConnectionManager
	IDataShutDown.IDataShutDown
	ISendMessage.ISendMessage
}
