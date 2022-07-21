package goConnectionManager

import "context"

type ICommandsToConnectionManager interface {
	CloseConnection(id string) error
	CloseAllConnections(ctx context.Context) error
}
