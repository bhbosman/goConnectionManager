package goConnectionManager

import (
	"github.com/bhbosman/gocommon/model"
	"go.uber.org/fx"
)

func ProvideRegisterToConnectionManager() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					ConnectionManager IService
				},
			) (model.IRegisterToConnectionManager, error) {
				return params.ConnectionManager, nil
			},
		},
	)
}
