package goConnectionManager

import "github.com/bhbosman/gocommon/Services/IFxService"

type IService interface {
	IConnectionManager
	IFxService.IFxServices
}
