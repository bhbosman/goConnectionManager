package goConnectionManager

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideCreateConnectionManager() fx.Option {
	return fx.Options(
		fx.Provide(
			func(
				params struct {
					fx.In
					Context           context.Context `name:"Application"`
					Logger            *zap.Logger
					ConnectionHelper  IHelper
					GoFunctionCounter GoFunctionCounter.IService
				},
			) (func() (IData, error), error) {
				return func() (IData, error) {
					return newData(
						params.Context,
						//params.PubSub,
						params.Logger.Named("ConnectionManagerData"),
						params.ConnectionHelper,
						params.GoFunctionCounter,
					)
				}, nil
			},
		),

		fx.Provide(
			func(
				params struct {
					fx.In
					Context           context.Context `name:"Application"`
					OnData            func() (IData, error)
					PubSub            *pubsub.PubSub `name:"Application"`
					Logger            *zap.Logger
					ConnectionHelper  IHelper
					GoFunctionCounter GoFunctionCounter.IService
				},
			) (IService, error) {
				cm := NewConnectionManagerService(
					params.Context,
					params.OnData,
					params.Logger.Named("ConnectionManagerData"),
					params.ConnectionHelper,
					params.PubSub,
					params.GoFunctionCounter,
				)
				return cm, nil
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						Service IService
					},
				) (IPublishConnectionInformation, error) {
					return params.Service, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						Service IService
					},
				) (model.IRegisterToConnectionManager, error) {
					return params.Service, nil
				},
			},
		),

		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						Service IService
					},
				) (IObtainConnectionManagerInformation, error) {
					return params.Service, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						Service IService
					},
				) (ICommandsToConnectionManager, error) {
					return params.Service, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						Lifecycle fx.Lifecycle
						PubSub    *pubsub.PubSub `name:"Application"`
					},
				) (IHelper, error) {
					result, err := NewHelper(
						params.PubSub,
					)
					if err != nil {
						return nil, err
					}
					params.Lifecycle.Append(
						fx.Hook{
							OnStart: nil,
							OnStop: func(ctx context.Context) error {
								return result.Close()
							},
						})
					return result, nil
				},
			},
		),
	)
}
