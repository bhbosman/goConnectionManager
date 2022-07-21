package goConnectionManager

import (
	"fmt"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"go.uber.org/fx"
)

func ProvideConnectionManager(ConnectionManager IService) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() (IService, error) {
				if ConnectionManager == nil {
					return nil, fmt.Errorf("ConnectionManager = nil, please resolve")
				}
				if ConnectionManager.State() != IFxService.Started {
					return nil, IFxService.NewServiceStateError(
						ConnectionManager.ServiceName(),
						"Service in incorrect state", IFxService.Started,
						ConnectionManager.State())
				}
				return ConnectionManager, nil
			},
		},
	)
}
