package goConnectionManager

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"math"
)

type data struct {
	dirtyMap                map[string]bool
	appContext              context.Context
	connectionMap           map[string]*model.ConnectionInformation
	messageRouter           *messageRouter.MessageRouter
	logger                  *zap.Logger
	ConnectionManagerHelper IHelper
	goFunctionCounter       GoFunctionCounter.IService
}

func (self *data) Send(message interface{}) error {
	if self.appContext.Err() != nil {
		self.logger.Error(
			"App Context in Error",
			zap.String("Method", "Send"),
			zap.Error(self.appContext.Err()))
		return self.appContext.Err()
	}
	_, err := self.messageRouter.Route(message)
	return err
}

func (self *data) ShutDown() error {
	if self.appContext.Err() != nil {
		self.logger.Error(
			"App Context in Error",
			zap.String("Method", "ShutDown"),
			zap.Error(self.appContext.Err()))
		return self.appContext.Err()
	}
	self.logger.Info("", zap.String("Method", "ShutDown"))
	return self.appContext.Err()
}

func (self *data) CloseAllConnections(_ context.Context) error {
	if self.appContext.Err() != nil {
		self.logger.Error(
			"App Context in Error",
			zap.String("Method", "CloseAllConnections"),
			zap.Error(self.appContext.Err()))
		return self.appContext.Err()
	}
	self.logger.Info("", zap.String("Method", "CloseAllConnections"))

	var err error
	for _, connectionInformation := range self.connectionMap {
		err = multierr.Append(
			err,
			self.goFunctionCounter.GoRun("ConnectionManager.CloseAllConnections",
				func(ci *model.ConnectionInformation) func() {
					return func() {
						ci.CancelFunc()
					}
				}(connectionInformation),
			),
		)
	}
	return err
}

func (self *data) CloseConnection(id string) error {
	if self.appContext.Err() != nil {
		self.logger.Error(
			"App Context in Error",
			zap.String("Method", "CloseConnection"),
			zap.Error(self.appContext.Err()))
		return self.appContext.Err()
	}
	self.logger.Info(
		"",
		zap.String("Method", "CloseConnection"),
		zap.String("Id", id))

	if ci, ok := self.connectionMap[id]; ok {
		ci.CancelFunc()
	}
	return nil
}

func (self *data) GetConnections(_ context.Context) ([]*model.ConnectionInformation, error) {
	if self.appContext.Err() != nil {
		self.logger.Error(
			"App Context in Error",
			zap.String("Method", "GetConnections"),
			zap.Error(self.appContext.Err()),
		)
		return nil, self.appContext.Err()
	}
	self.logger.Info("", zap.String("Method", "GetConnections"))

	var result []*model.ConnectionInformation
	for _, ci := range self.connectionMap {
		result = append(result, ci)
	}
	return result, nil
}

func (self *data) NameConnection(id string, name string) error {
	if self.appContext.Err() != nil {
		self.logger.Error(
			"App Context in Error",
			zap.String("Method", "NameConnection"),
			zap.Error(self.appContext.Err()),
		)
		return self.appContext.Err()
	}
	self.logger.Info("",
		zap.String("Method", "NameConnection"),
		zap.String("Id", id),
		zap.String("name", name))

	if ci, ok := self.connectionMap[id]; ok {
		ci.Name = name
		self.dirtyMap[id] = true
		self.publishMessage(model.NewConnectionCreated(id, name, ci.CancelFunc, ci.ConnectionTime, ci.CancelContext))
	}
	return nil
}

func (self *data) RegisterConnection(id string, function context.CancelFunc, CancelContext context.Context) error {
	if self.appContext.Err() != nil {
		self.logger.Error(
			"App Context in Error",
			zap.String("Method", "RegisterConnection"),
			zap.Error(self.appContext.Err()),
		)
		return self.appContext.Err()
	}
	self.logger.Info("",
		zap.String("Method", "RegisterConnection"),
		zap.String("id", id))

	value := model.NewConnectionInformation(id, function, CancelContext)
	self.connectionMap[id] = value
	self.publishMessage(model.NewConnectionCreated(id, "(unassigned)", function, value.ConnectionTime, CancelContext))
	return nil
}

func (self *data) DeregisterConnection(id string) error {
	if self.appContext.Err() != nil {
		self.logger.Error(
			"App Context in Error",
			zap.String("Method", "DeregisterConnection"),
			zap.Error(self.appContext.Err()),
		)
		return self.appContext.Err()
	}
	self.logger.Info("",
		zap.String("Method", "DeregisterConnection"),
		zap.String("id", id),
	)

	delete(self.connectionMap, id)
	self.publishMessage(&model.ConnectionClosed{ConnectionId: id})
	return nil
}

func (self *data) ConnectionInformationReceived(counters *model.PublishRxHandlerCounters) error {
	if connectionInformation, ok := self.connectionMap[counters.ConnectionId]; ok {
		if counters.Direction == model.StreamDirectionInbound {
			connectionInformation.InboundCounters = counters
		} else {
			connectionInformation.OutboundCounters = counters
		}

		if connectionInformation.OutboundCounters != nil && connectionInformation.InboundCounters != nil {
			l := int(math.Max(
				float64(len(connectionInformation.InboundCounters.Counters)),
				float64(len(connectionInformation.OutboundCounters.Counters))))
			grid := make([]*model.LineData, l)
			for i := 0; i < l; i++ {
				grid[i] = &model.LineData{}
			}
			for i, counter := range connectionInformation.OutboundCounters.Counters {
				grid[i].OutValue = counter.Data
			}
			for i, counter := range connectionInformation.InboundCounters.Counters {
				index := l - i - 1
				grid[index].InValue = counter.Data
			}
			connectionInformation.Grid = grid
		}
		self.dirtyMap[counters.ConnectionId] = true
	}
	return nil
}
func (self *data) handleRefreshDataTo(msg *RefreshDataTo) {
	self.ConnectionManagerHelper.Pub(
		&RefreshDataStart{
			SubscriptionName: msg.SubscriptionName,
		},
		msg.SubscriptionName,
	)
	for _, cm := range self.connectionMap {

		gridData := make([]model.LineData, len(cm.Grid), len(cm.Grid))
		for i, lineData := range cm.Grid {
			gridData[i] = *lineData
		}

		self.ConnectionManagerHelper.Pub(
			model.NewConnectionCreated(
				cm.Id,
				cm.Name,
				cm.CancelFunc,
				cm.ConnectionTime,
				cm.CancelContext,
			),
			msg.SubscriptionName)
		self.ConnectionManagerHelper.Pub(
			&model.ConnectionState{
				ConnectionId:   cm.Id,
				CancelContext:  cm.CancelContext,
				CancelFunc:     cm.CancelFunc,
				Name:           cm.Name,
				ConnectionTime: cm.ConnectionTime,
				Grid:           gridData,
			},
			msg.SubscriptionName,
		)

	}
	self.ConnectionManagerHelper.Pub(
		&RefreshDataStop{
			SubscriptionName: msg.SubscriptionName,
		},
		msg.SubscriptionName,
	)
}

func (self *data) handleEmptyQueue(msg *messages.EmptyQueue) error {
	if self.appContext.Err() != nil {
		self.logger.Error(
			"App Context in Error",
			zap.String("Method", "handleEmptyQueue"),
			zap.Error(self.appContext.Err()))
		return self.appContext.Err()
	}
	self.logger.Info("",
		zap.String("Method", "handleEmptyQueue"),
		zap.Int("Count", msg.Count),
		zap.Int("Overall", msg.OverallCount))

	for s, b := range self.dirtyMap {
		if b {
			if connInfo, ok := self.connectionMap[s]; ok {
				gridData := make([]model.LineData, len(connInfo.Grid), len(connInfo.Grid))
				for i, lineData := range connInfo.Grid {
					gridData[i] = *lineData
				}
				self.publishMessage(
					&model.ConnectionState{
						ConnectionId:   connInfo.Id,
						CancelContext:  connInfo.CancelContext,
						CancelFunc:     connInfo.CancelFunc,
						Name:           connInfo.Name,
						ConnectionTime: connInfo.ConnectionTime,
						Grid:           gridData,
					})
			}
		}
	}
	self.dirtyMap = make(map[string]bool)

	return nil
}

func (self *data) publishMessage(msg interface{}) {
	self.ConnectionManagerHelper.Pub(
		msg,
		self.ConnectionManagerHelper.PublishChannelName(),
	)
}

func (self *data) publishDirtyList() {

}

func newData(
	appContext context.Context,
	logger *zap.Logger,
	ConnectionHelper IHelper,
	goFunctionCounter GoFunctionCounter.IService,
) (IData, error) {
	result := &data{
		dirtyMap:                make(map[string]bool),
		appContext:              appContext,
		connectionMap:           make(map[string]*model.ConnectionInformation),
		messageRouter:           messageRouter.NewMessageRouter(),
		logger:                  logger,
		ConnectionManagerHelper: ConnectionHelper,
		goFunctionCounter:       goFunctionCounter,
	}
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	_ = result.messageRouter.Add(result.handleRefreshDataTo)
	return result, nil
}
