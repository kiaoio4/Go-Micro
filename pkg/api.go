package pkg

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

func (micro *GoMicro) handlerWrapper(selfServiceName string, h echo.HandlerFunc) echo.HandlerFunc {
	if micro.gossipKVCache != nil {
		return micro.gossipKVCache.GoMicroHandlerWrapper(selfServiceName, h, true)
	}
	return h
}

// SetupWeb Set interface
func (micro *GoMicro) SetupWeb(root echoswagger.ApiRoot, base, selfServiceName string) {
	g := root.Group(API, base)
	g.GET("/path", micro.handlerWrapper(selfServiceName, micro.getTestData)).
		AddParamQuery("", "path", "store path", true).
		AddResponse(http.StatusOK, `
		{
			"code": 0,
			"msg": "OK",
			"data": [
				"/home/workspace"
			]
		}
		`, "", nil).
		AddResponse(http.StatusNotFound, `
		{
			"code": 404,
			"msg": "Not Found"
		}		
		`, nil, nil).
		AddResponse(http.StatusTooManyRequests, `
		{
			"code": 429,
			"msg": "Too Many Requests:"+ id
		}		
		`, nil, nil).
		AddResponse(http.StatusServiceUnavailable, `
		{
			"code":503,
			"msg":"Service Unavailable"
		}	
		`, nil, nil).
		SetOperationId("path").
		SetSummary("Get information of store path")
}
