package goConnectionManager

import (
	"go.uber.org/fx"
)

func ProvideObtainConnectionManagerInformation() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					ConnectionManager IService
				},
			) (IObtainConnectionManagerInformation, error) {
				return params.ConnectionManager, nil
			},
		},
	)
}
