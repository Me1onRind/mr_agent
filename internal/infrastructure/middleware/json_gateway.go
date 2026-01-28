package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HTTPHandler[A any, B any] func(c context.Context, request *A) (data *B, err error)

func JSON[A any, B any](handler HTTPHandler[A, B]) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", jsonGateway(c, handler))
		c.Next()
	}
}

func jsonGateway[A any, B any](c *gin.Context, handler HTTPHandler[A, B]) []byte {
	ctx := c.Request.Context()

	var response *JsonResponse
	var request A
	if err := c.ShouldBind(&request); err != nil {
		response = &JsonResponse{
			Code:    -1,
			Message: fmt.Sprintf("Decode request text fail, cause:[%s]", err),
		}
	} else {
		data, err := handler(ctx, &request)
		response = getResponse(data, err)
	}

	jsonData, err := jsoniter.Marshal(response)
	if err != nil {
		logger.CtxLoggerWithSpandId(ctx).Error("Marshal response fail", slog.String("error", err.Error()))
		jsonData, _ = jsoniter.Marshal(&JsonResponse{
			Code:    -2,
			Message: fmt.Sprintf("JSON Gateway encode response fail, err:[%s]", err.Error()),
		})
	}
	return jsonData
}

func getResponse(data any, err error) *JsonResponse {
	response := &JsonResponse{}
	if err == nil {
		response.Code = 0
		response.Message = "Success"
		response.Data = data
		return response
	}
	response.Message = err.Error()
	//if expectErr := customErr.ExtractError(err); expectErr != nil {
	//response.Code = expectErr.Code
	//} else {
	//response.Code = code.Unexpect
	//}

	//if code.IsWarning(response.Code) {
	//response.Data = data
	//}
	return response
}
