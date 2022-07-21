package goConnectionManager

import (
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
)

type IData interface {
	IConnectionManager
	IDataShutDown.IDataShutDown
	ISendMessage.ISendMessage
}
