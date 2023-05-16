package lib

import (
	"common/define"
	"github.com/gin-gonic/gin"
)

func CommonError(c *gin.Context, code int, err string) {
	retRes := define.NetResponse{}
	retRes.Code = code
	retRes.Error = err
	c.JSON(200, retRes)
}

func CreateResponseError(code int, err string) *define.Response {
	return &define.Response{
		Code:  code,
		Error: err,
		Data:  map[string]any{},
	}
}
func CreateResponseSuccess(data any) *define.Response {
	return &define.Response{
		Code:  200,
		Error: "",
		Data:  data,
	}
}
func CreateRetResponseError(code int, err string) *define.NetResponse {
	retRes := define.NetResponse{}
	retRes.Code = code
	retRes.Error = err
	return &retRes
}
