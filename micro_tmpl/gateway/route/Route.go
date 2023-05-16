package route

import (
	"common/define"
	"context"
	"micro_tmpl/gateway/handler"
	"sync"
)

type Route struct {
}

var once sync.Once
var instance *Route

func GetRouteInstance() *Route {
	once.Do(func() {
		instance = &Route{}
	})
	return instance
}

func (route *Route) DispatchApiRequest(ctx context.Context, request *define.Request) *define.Response {
	var res = &define.Response{}

	switch request.Module {
	case "player":
		playerHandler := handler.GetHandlerInstance()
		switch request.Method {
		case "login":
			res = playerHandler.Login(ctx, request)
		case "register":
			res = playerHandler.Register(ctx, request)
		}
	}
	return res
}
