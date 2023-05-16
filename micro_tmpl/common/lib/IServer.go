package lib

import (
	"common/define"
	"context"
)

type IServer interface {
	Run()
	Stop()
	RegisterRoutes(routeName string, handle func(ctx context.Context, request *define.Request) *define.Response)
}
