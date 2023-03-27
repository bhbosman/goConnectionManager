package goConnectionManager

import (
	"github.com/bhbosman/gocommon/services/IFxService"
)

type IService interface {
	IConnectionManager
	IFxService.IFxServices
}
