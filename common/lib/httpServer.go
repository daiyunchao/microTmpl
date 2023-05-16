package lib

import (
	"common/define"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpServer struct {
	listenAddress string
	r             *gin.Engine
	encodeKey     string
}

// CreateHttpServer 创建一个Http服务器
func CreateHttpServer(listenAddress string, encodeKey string) *HttpServer {
	httpServer := &HttpServer{
		listenAddress: listenAddress,
		r:             gin.Default(),
		encodeKey:     encodeKey,
	}
	return httpServer
}

// Run 启动服务器
func (http *HttpServer) Run() {
	go http.r.Run(http.listenAddress)
}

// RegisterRoutes 注册路由
func (http *HttpServer) RegisterRoutes(routeName string, handle func(ctx context.Context, req *define.Request) *define.Response) {
	http.r.Use(timeoutMiddleware(time.Second * 60))
	http.r.TrustedPlatform = "X-Client-IP"
	http.r.POST(routeName, func(c *gin.Context) {
		reqId := GetNowOffTimeStamp()
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "reqId", reqId)
		resChan := make(chan *define.NetResponse, 1)
		go func() {
			Logger.Infof("reqId: %d, request Start: %s", reqId, c.Request.URL)

			//处理请求
			isNotDebug := c.GetHeader("debug") != "1"
			request, requestErr := http.handleRequest(c, reqId, false, isNotDebug)
			if requestErr.Error != "" {
				resChan <- requestErr
				return
			}
			//执行主逻辑回调
			handlerRes := handle(ctx, &request)

			//处理返回
			ret := http.handleResponse(reqId, false, isNotDebug, handlerRes)
			resChan <- ret
		}()
		select {
		case retRes := <-resChan:
			c.JSON(200, retRes)
		case <-ctx.Done():
			retRes := define.NetResponse{
				Code:  define.Timeout,
				Error: define.TimeoutMsg,
			}
			c.JSON(200, retRes)
			Logger.Infof("reqId: %d,finish,res: %s", reqId, retRes)
		}
	})
}

func (http *HttpServer) RegisterStatic(routeName string, path string) {
	http.r.Static(routeName, path)
}

// Stop 服务器停止
func (http *HttpServer) Stop() {

}

func (http *HttpServer) handleRequest(c *gin.Context, reqId int64, isAdmin bool, isNotDebug bool) (define.Request, *define.NetResponse) {
	//解析出参数
	request := define.Request{
		Data: make(map[string]any),
	}
	module := c.Param("module")
	method := c.Param("method")
	request.Module = module
	request.Method = method
	request.Ip = c.ClientIP()
	//验证token
	if !isAdmin {
		tokenIsValid := validateHeader(c, module, method)
		if !tokenIsValid {
			return request, CreateRetResponseError(define.TokenError, define.TokenErrorMsg)
		}
	}

	if isNotDebug && !isAdmin {
		//解密:
		newData, err := dataDecode(c.Request.Body, http.encodeKey)
		if err != nil {
			Logger.Warnf("reqId: %d,解析参数错误", reqId)
			return request, CreateRetResponseError(define.ParamError, define.ParamErrorMsg)
		}
		request.Data = newData
	} else {
		//不解密:
		params := define.OriJsonRequest{}
		err := c.ShouldBindJSON(&params)
		if err != nil {
			Logger.Warnf("reqId: %d,解析参数错误", reqId)
			return request, CreateRetResponseError(define.ParamError, define.ParamErrorMsg)
		}
		request.Data = params.Data
	}
	return request, &define.NetResponse{}
}

func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {
				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}
			cancel()
		}()
		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func validateHeader(c *gin.Context, module string, method string) bool {
	return true
}

// 解密
func dataDecode(data io.ReadCloser, key string) (map[string]any, error) {
	bytes, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}
	oriRequest := define.OriRequest{}
	err = json.Unmarshal(bytes, &oriRequest)
	if err != nil {
		return nil, err
	}
	encode := AESEnCode{}
	decode, err := encode.Decode([]byte(oriRequest.Data), key)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]any)
	err = json.Unmarshal(decode, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (http *HttpServer) handleResponse(reqId int64, isAdmin bool, isNotDebug bool, res *define.Response) *define.NetResponse {
	retRes := define.NetResponse{}
	if isNotDebug && !isAdmin {
		//加密
		bytes, err := json.Marshal(res.Data)
		if err != nil {
			Logger.Warnf("reqId: %d,AESMarshal Has Some Error :%v", reqId, err)
			return CreateRetResponseError(define.ServerError, "服务器内部错误")
		}
		if string(bytes) == "{}" {
			//空对象,不需要再加密了
		} else {
			base64Encode := AESEnCode{}
			encode, err := base64Encode.Encode(bytes, http.encodeKey)
			if err != nil {
				Logger.Warnf("reqId: %d,AESEnCode Has Some Error :%v", reqId, err)
				return CreateRetResponseError(define.ServerError, "服务器内部错误")
			}
			retRes.Data = string(encode)
		}
	} else {
		//不加密
		bytes, err := json.Marshal(res.Data)
		if err != nil {
			Logger.Warnf("reqId: %d,AESMarshal Has Some Error :%v", reqId, err)
			return CreateRetResponseError(define.ServerError, "服务器内部错误")
		}
		retRes.Code = res.Code
		retRes.Error = res.Error
		retRes.Data = string(bytes)
	}
	return &retRes
}
