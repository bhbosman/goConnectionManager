package goConnectionManager

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/ISendMessage"
	model2 "github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

type Service struct {
	parentContext     context.Context
	ctx               context.Context
	cancelFunc        context.CancelFunc
	cmdChannel        chan interface{}
	onData            func() (IData, error)
	Logger            *zap.Logger
	state             IFxService.State
	ConnectionHelper  IHelper
	pubSub            *pubsub.PubSub
	goFunctionCounter GoFunctionCounter.IService
	subscribeChannel  *pubsub.ChannelSubscription
}

func (self *Service) ServiceName() string {
	return "ConnectionManager"
}

func (self *Service) State() IFxService.State {
	return self.state
}

func (self *Service) ConnectionInformationReceived(counters *model2.PublishRxHandlerCounters) error {
	result, err := CallIConnectionManagerConnectionInformationReceived(
		self.ctx, self.cmdChannel, true,
		counters)
	if err != nil {
		return nil
	}
	return result.Args0
}

func (self *Service) CloseAllConnections(ctx context.Context) error {
	result, err := CallIConnectionManagerCloseAllConnections(self.ctx, self.cmdChannel, true, ctx)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *Service) CloseConnection(id string) error {
	result, err := CallIConnectionManagerCloseConnection(self.ctx, self.cmdChannel, true, id)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *Service) GetConnections(ctx context.Context) ([]*model2.ConnectionInformation, error) {
	result, err := CallIConnectionManagerGetConnections(self.ctx, self.cmdChannel, true, ctx)
	if err != nil {
		return nil, err
	}
	return result.Args0, result.Args1
}

func (self *Service) NameConnection(id string, connectionName string) error {
	result, err := CallIConnectionManagerNameConnection(self.ctx, self.cmdChannel, true,
		id, connectionName)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *Service) OnStart(ctx context.Context) error {
	err := self.start(ctx)
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *Service) OnStop(ctx context.Context) error {
	err := self.shutdown(ctx)
	close(self.cmdChannel)
	self.state = IFxService.Stopped
	return err
}

func (self *Service) RegisterConnection(id string, function context.CancelFunc, CancelContext context.Context) error {
	result, err := CallIConnectionManagerRegisterConnection(self.ctx, self.cmdChannel, true,
		id, function, CancelContext)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *Service) DeregisterConnection(id string) error {
	result, err := CallIConnectionManagerDeregisterConnection(self.ctx, self.cmdChannel, true, id)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *Service) shutdown(_ context.Context) error {
	self.cancelFunc()
	return pubSub.Unsubscribe("Unsubscribe", self.pubSub, self.goFunctionCounter, self.subscribeChannel)
}

func (self *Service) start(_ context.Context) error {
	instanceData, err := self.onData()
	if err != nil {
		return err
	}

	return self.goFunctionCounter.GoRun(
		"Connection Manager Service",
		func() {
			self.goStart(instanceData)
		},
	)
}

func (self *Service) goStart(instanceData IData) {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.cmdChannel)
	self.subscribeChannel = pubsub.NewChannelSubscription(32)
	self.pubSub.AddSub(self.subscribeChannel, self.ConnectionHelper.RefreshChannelName())

	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		instanceData,
		[]ChannelHandler.ChannelHandler{
			{
				//BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(IConnectionManager); ok {
						return ChannelEventsForIConnectionManager(unk, message)
					}
					return false, nil
				},
			},
			{
				//BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ISendMessage.ISendMessage); ok {
						return true, unk.Send(message)
					}
					return false, nil
				},
			},
		},
		func() int {
			n := len(self.cmdChannel) + self.subscribeChannel.Count()
			return n
		},
		goCommsDefinitions.CreateTryNextFunc(self.cmdChannel),
		//func(i interface{}) {
		//	select {
		//	case self.cmdChannel <- i:
		//		break
		//	default:
		//		break
		//	}
		//},
	)
loop:
	for {
		select {
		case <-self.ctx.Done():
			err := instanceData.ShutDown()
			if err != nil {
				self.Logger.Error(
					"error on done",
					zap.Error(err))
			}
			break loop
		case event, ok := <-self.cmdChannel:
			if !ok {
				return
			}
			breakLoop, err := channelHandlerCallback(event)
			if err != nil || breakLoop {
				break loop
			}
		case event, ok := <-self.subscribeChannel.Data:
			if !ok {
				return
			}
			breakLoop, err := channelHandlerCallback(event)
			if err != nil || breakLoop {
				break loop
			}
		}
	}
}

func NewConnectionManagerService(
	parentContext context.Context,
	onData func() (IData, error),
	logger *zap.Logger,
	ConnectionHelper IHelper,
	pubSub *pubsub.PubSub,
	goFunctionCounter GoFunctionCounter.IService,
) IService {
	localCtx, localCancelFunc := context.WithCancel(parentContext)

	result := &Service{
		parentContext:     parentContext,
		ctx:               localCtx,
		cancelFunc:        localCancelFunc,
		cmdChannel:        make(chan interface{}, 32),
		onData:            onData,
		Logger:            logger,
		ConnectionHelper:  ConnectionHelper,
		pubSub:            pubSub,
		goFunctionCounter: goFunctionCounter,
	}

	return result
}
