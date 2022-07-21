package goConnectionManager

import (
	"go.uber.org/fx"
)

func ProvideCommandsToConnectionManager() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					ConnectionManager IService
				},
			) (ICommandsToConnectionManager, error) {
				return params.ConnectionManager, nil
			},
		},
	)
}
