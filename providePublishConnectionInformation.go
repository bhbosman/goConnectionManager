package goConnectionManager

import (
	"go.uber.org/fx"
)

func ProvidePublishConnectionInformation() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					ConnectionManager IService
				},
			) (IPublishConnectionInformation, error) {
				return params.ConnectionManager, nil
			},
		},
	)
}
