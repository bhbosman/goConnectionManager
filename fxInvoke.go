package goConnectionManager

import "go.uber.org/fx"

func InvokeConnectionManager() fx.Option {
	return fx.Options(
		fx.Invoke(
			func(
				params struct {
					fx.In
					LifeCycle         fx.Lifecycle
					ConnectionManager IService
				},
			) {
				params.LifeCycle.Append(
					fx.Hook{
						OnStart: params.ConnectionManager.OnStart,
						OnStop:  params.ConnectionManager.OnStop,
					},
				)
			},
		),
	)
}
