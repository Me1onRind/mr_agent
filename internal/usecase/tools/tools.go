package tools

import (
	"context"

	agentTools "github.com/Me1onRind/mr_agent/internal/infrastructure/agent/tools"
	"github.com/Me1onRind/mr_agent/internal/protocol/http/tools"
	jsoniter "github.com/json-iterator/go"
)

type ToolsUsecase struct {
}

func NewToolsUsecase() *ToolsUsecase {
	return &ToolsUsecase{}
}

func (t *ToolsUsecase) Call(ctx context.Context, request *tools.CallToolReqest) (*tools.CallToolResponse, error) {
	tool, err := agentTools.DefaultRegistry.Get(request.Tool)
	if err != nil {
		return nil, err
	}

	params, err := jsoniter.Marshal(request.Params)
	if err != nil {
		return nil, err
	}

	result, err := tool.Handler(ctx, params)
	if err != nil {
		return nil, err
	}

	response := tools.CallToolResponse{
		Result: result,
	}

	return &response, nil
}
