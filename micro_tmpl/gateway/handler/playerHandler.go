package handler

import (
	"common/define"
	"common/lib"
	"context"
	"micro_tmpl/gateway/gatewayDefine"
	"micro_tmpl/gateway/proto"
	"sync"
)

var once sync.Once
var instance *Handler

func GetHandlerInstance() *Handler {
	once.Do(func() {
		instance = &Handler{}
	})
	return instance
}

type Handler struct {
}

func (handler *Handler) Login(ctx context.Context, request *define.Request) *define.Response {
	//参数验证:
	req := gatewayDefine.ReqLoginData{}
	request.DecodeData(&req)
	if len(req.UserName) <= 0 || len(req.Password) <= 0 {
		return lib.CreateResponseError(define.ParamError, define.ParamErrorMsg)
	}
	//grpc
	rpc := GetRpcServer()
	conn, client := rpc.getPlayerConn()
	defer client.Close()
	rpcReq := &proto.ReqGetInfoByName{
		Name: req.UserName,
	}
	rpcRes, err := conn.GetInfoByName(ctx, rpcReq)
	if err != nil {
		return lib.CreateResponseError(define.ServerError, define.ServerErrorMsg)
	}
	if rpcRes.Id == "" {
		return lib.CreateResponseError(define.UserNameOrPasswordError, define.UserNameOrPasswordErrorMsg)
	}
	return lib.CreateResponseSuccess(rpcRes)
}

func (handler *Handler) Register(ctx context.Context, request *define.Request) *define.Response {
	//参数验证:
	req := gatewayDefine.ReqRegisterData{}
	request.DecodeData(&req)
	if len(req.Name) <= 0 || len(req.Password) <= 0 {
		return lib.CreateResponseError(define.ParamError, define.ParamErrorMsg)
	}
	//grpc
	rpc := GetRpcServer()
	conn, client := rpc.getPlayerConn()
	defer client.Close()
	rpcReq := &proto.ReqRegister{
		Name:     req.Name,
		Password: req.Password,
	}
	rpcRes, err := conn.Register(ctx, rpcReq)
	if err != nil {
		return lib.CreateResponseError(define.ServerError, define.ServerErrorMsg)
	}
	if rpcRes.Id == "" {
		return lib.CreateResponseError(define.ServerError, define.ServerErrorMsg)
	}
	return lib.CreateResponseSuccess(rpcRes)
}
