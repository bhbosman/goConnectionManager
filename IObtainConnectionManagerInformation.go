package goConnectionManager

import (
	"context"
	"github.com/bhbosman/gocommon/model"
)

type IObtainConnectionManagerInformation interface {
	GetConnections(ctx context.Context) ([]*model.ConnectionInformation, error)
}
