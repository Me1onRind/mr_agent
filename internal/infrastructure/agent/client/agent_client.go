package agent

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"go.elastic.co/apm"
)

type AgentClient struct {
	openaiClient *openai.Client
}

func NewAgentClient() *AgentClient {
	a := &AgentClient{}
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("QWEN_TOKEN")),
		option.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
	)
	a.openaiClient = &client
	return a
}

func (a *AgentClient) Chat(ctx context.Context,
	msgs []openai.ChatCompletionMessageParamUnion) (*openai.ChatCompletion, error) {
	return a.ChatWithTools(ctx, msgs, nil)
}

func (a *AgentClient) ChatWithTools(
	ctx context.Context,
	msgs []openai.ChatCompletionMessageParamUnion,
	tools []openai.ChatCompletionToolUnionParam,
) (*openai.ChatCompletion, error) {
	startTime := time.Now()
	log := logger.LoggerFromCtx(ctx)
	defer func() {
		log.Info("new agent done", slog.Int64("latency", time.Since(startTime).Milliseconds()))
	}()
	span, _ := apm.StartSpan(ctx, "new_agent", "openai")
	defer span.End()
	client := a.openaiClient
	params := openai.ChatCompletionNewParams{
		Messages: msgs,
		Model:    "qwen-plus",
	}
	if len(tools) > 0 {
		params.Tools = tools
	}
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), params)

	if err != nil {
		log.Error("client.Chat.Completions.New failed", slog.String("error", err.Error()))
		return nil, err
	}

	//return chatCompletion.Choices[0].Message.Content, nil
	return chatCompletion, nil
}
