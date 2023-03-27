package goConnectionManager

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	model2 "github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/bhbosman/gocommon/services/IFxService"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

type service struct {
	parentContext     context.Context
	ctx               context.Context
	cancelFunc        context.CancelFunc
	cmdChannel        chan interface{}
	onData            func() (IData, error)
	logger            *zap.Logger
	state             IFxService.State
	ConnectionHelper  IHelper
	pubSub            *pubsub.PubSub
	goFunctionCounter GoFunctionCounter.IService
	subscribeChannel  *pubsub.NextFuncSubscription
}

func (self *service) Send(message interface{}) error {
	result, err := CallIConnectionManagerSend(self.ctx, self.cmdChannel, false, message)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) MultiSend(messages ...interface{}) {
	_, err := CallIConnectionManagerMultiSend(self.ctx, self.cmdChannel, false, messages...)
	if err != nil {
		return
	}
}

func (self *service) ServiceName() string {
	return "ConnectionManager"
}

func (self *service) State() IFxService.State {
	return self.state
}

func (self *service) ConnectionInformationReceived(counters *model2.PublishRxHandlerCounters) error {
	result, err := CallIConnectionManagerConnectionInformationReceived(
		self.ctx, self.cmdChannel, true,
		counters)
	if err != nil {
		return nil
	}
	return result.Args0
}

func (self *service) CloseAllConnections(ctx context.Context) error {
	result, err := CallIConnectionManagerCloseAllConnections(self.ctx, self.cmdChannel, true, ctx)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) CloseConnection(id string) error {
	result, err := CallIConnectionManagerCloseConnection(self.ctx, self.cmdChannel, true, id)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) GetConnections(ctx context.Context) ([]*model2.ConnectionInformation, error) {
	result, err := CallIConnectionManagerGetConnections(self.ctx, self.cmdChannel, true, ctx)
	if err != nil {
		return nil, err
	}
	return result.Args0, result.Args1
}

func (self *service) NameConnection(id string, connectionName string) error {
	result, err := CallIConnectionManagerNameConnection(self.ctx, self.cmdChannel, true,
		id, connectionName)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) OnStart(ctx context.Context) error {
	err := self.start(ctx)
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *service) OnStop(ctx context.Context) error {
	err := self.shutdown(ctx)
	self.cmdChannel <- 222
	self.cmdChannel <- 222
	self.cmdChannel <- 222
	close(self.cmdChannel)
	for _ = range self.cmdChannel {
		fmt.Sprintf("dddd")
	}
	self.state = IFxService.Stopped
	return err
}

func (self *service) RegisterConnection(
	id string,
	function context.CancelFunc,
	CancelContext context.Context,
	nextFuncOutBoundChannel rxgo.NextFunc,
	nextFuncInBoundChannel rxgo.NextFunc,
) error {
	result, err := CallIConnectionManagerRegisterConnection(
		self.ctx,
		self.cmdChannel,
		true,
		id,
		function,
		CancelContext,
		nextFuncOutBoundChannel,
		nextFuncInBoundChannel,
	)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) DeregisterConnection(id string) error {
	result, err := CallIConnectionManagerDeregisterConnection(self.ctx, self.cmdChannel, true, id)
	if err != nil {
		return err
	}
	return result.Args0
}

func (self *service) shutdown(_ context.Context) error {
	self.cancelFunc()
	return pubSub.Unsubscribe("Unsubscribe", self.pubSub, self.goFunctionCounter, self.subscribeChannel)
}

func (self *service) start(_ context.Context) error {
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

func (self *service) goStart(instanceData IData) {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.cmdChannel)
	self.subscribeChannel = pubsub.NewNextFuncSubscription(goCommsDefinitions.CreateNextFunc(self.cmdChannel))
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
	)
loop:
	for {
		select {
		case <-self.ctx.Done():
			err := instanceData.ShutDown()
			if err != nil {
				self.logger.Error(
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

	result := &service{
		parentContext:     parentContext,
		ctx:               localCtx,
		cancelFunc:        localCancelFunc,
		cmdChannel:        make(chan interface{}, 32),
		onData:            onData,
		logger:            logger,
		ConnectionHelper:  ConnectionHelper,
		pubSub:            pubSub,
		goFunctionCounter: goFunctionCounter,
	}

	return result
}
