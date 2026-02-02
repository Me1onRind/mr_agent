package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Me1onRind/mr_agent/internal/errcode"
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
	}
}

func jsonGateway[A any, B any](c *gin.Context, handler HTTPHandler[A, B]) []byte {
	ctx := c.Request.Context()

	var response *JsonResponse
	var request A
	if err := c.ShouldBind(&request); err != nil {
		response = &JsonResponse{
			Code:    errcode.ParamInvalidCode,
			Message: fmt.Sprintf("Decode Request Body Failed, cause:[%s]", err),
		}
	} else {
		data, err := handler(ctx, &request)
		response = getResponse(data, err)
	}

	jsonData, err := jsoniter.Marshal(response)
	if err != nil {
		logger.LoggerFromCtx(ctx).Error("marshal response failed", slog.String("error", err.Error()))
		jsonData, _ = jsoniter.Marshal(&JsonResponse{
			Code:    errcode.JsonEncodeFailedCode,
			Message: fmt.Sprintf("Encode Response Object Failed, cause:[%s]", err.Error()),
		})
	}
	return jsonData
}

func getResponse(data any, err error) *JsonResponse {
	response := &JsonResponse{}
	if err == nil {
		response.Code = errcode.SuccessCode
		response.Message = "Success"
		response.Data = data
		return response
	}
	response.Message = err.Error()
	if expectErr := errcode.ExtractError(err); expectErr != nil {
		response.Code = expectErr.Code
	} else {
		response.Code = errcode.UnexpectCode
	}

	if errcode.IsWarning(response.Code) {
		response.Data = data
	}
	return response
}
