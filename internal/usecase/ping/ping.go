package ping

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/client/agent"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/Me1onRind/mr_agent/internal/protocol/http/ping"
	"github.com/openai/openai-go/v3"
)

type PingUsecase struct {
}

func NewPingUsecase() *PingUsecase {
	return &PingUsecase{}
}

func (p *PingUsecase) Echo(ctx context.Context, request *ping.EchoRequest) (*ping.EchoResponse, error) {
	resp := &ping.EchoResponse{
		Msg: request.Msg,
	}
	return resp, nil
}

func (p *PingUsecase) Panic(ctx context.Context, request *struct{}) (*struct{}, error) {
	panic("This is panic")
}

func (p *PingUsecase) HelloToAgent(ctx context.Context, request *ping.HelloToAgentRequest) (*ping.HelloToAgentResponse, error) {
	log := logger.CtxLogger(ctx)
	agentClient := agent.NewAgentClient(ctx)
	completion, err := agentClient.Chat(ctx, []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(request.Msg),
	})
	if err != nil {
		return nil, err
	}
	log.Info("new agent", slog.String("completion_id", completion.ID))
	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("chat completion has no choices")
	}
	return &ping.HelloToAgentResponse{
		Msg: completion.Choices[0].Message.Content,
	}, nil
}
