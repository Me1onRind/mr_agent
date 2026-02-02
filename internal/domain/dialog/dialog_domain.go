package dialog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	agent "github.com/Me1onRind/mr_agent/internal/infrastructure/agent/client"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/agent/tools"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/cache/llm"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/session"
	jsoniter "github.com/json-iterator/go"
	"github.com/openai/openai-go/v3"
)

type DialogDomain struct {
	AgentClient *agent.AgentClient
}

func NewDialogDomain() *DialogDomain {
	d := &DialogDomain{
		AgentClient: agent.NewAgentClient(),
	}
	return d
}

type DialogConfig struct {
	WithCtx       bool
	SystemMessage string
}

func (d *DialogDomain) Dialog(ctx context.Context, msg string, cfg *DialogConfig) (string, error) {
	log := logger.LoggerFromCtx(ctx)
	userMsg := openai.UserMessage(msg)

	var (
		msgs []openai.ChatCompletionMessageParamUnion
	)

	sessionId, err := session.GetSessionId(ctx)
	if err != nil {
		return "", err
	}

	if cfg.WithCtx {
		msgs, err = llm.GetMessages(ctx, sessionId)
		if err != nil && !errors.Is(err, llm.ErrMessagesNotFound) {
			log.Error("LLMCache.GetMessages failed", slog.String("error", err.Error()))
			return "", err
		}
	}

	if cfg.SystemMessage != "" {
		msgs = append(msgs, openai.SystemMessage(cfg.SystemMessage))
	}
	msgs = append(msgs, userMsg)
	openAITools := tools.DefaultRegistry.OpenAITools()

	var reply string
	maxToolRounds := 5
	for range maxToolRounds {
		completion, err := d.AgentClient.ChatWithTools(ctx, msgs, openAITools)
		if err != nil {
			return "", err
		}
		if len(completion.Choices) == 0 {
			return "", fmt.Errorf("chat completion has no choices")
		}

		message := completion.Choices[0].Message
		toolCalls := message.ToolCalls
		if len(toolCalls) == 0 {
			reply = message.Content
			break
		}

		msgs = append(msgs, message.ToParam())
		msgs, err = callTools(ctx, msgs, toolCalls)
		if err != nil {
			return "", err
		}
	}

	if reply == "" {
		return "", fmt.Errorf("tool rounds exceeded without final response")
	}
	assistant := openai.AssistantMessage(reply)

	if cfg.WithCtx {
		if err := llm.AppendMsgAndAss(ctx, sessionId, userMsg, assistant); err != nil {
			log.Warn("LLMCache.AppendMsgAndAss failed", slog.String("error", err.Error()))
		}
	}

	return reply, nil
}

func callTools(ctx context.Context,
	msgs []openai.ChatCompletionMessageParamUnion,
	toolCalls []openai.ChatCompletionMessageToolCallUnion) ([]openai.ChatCompletionMessageParamUnion, error) {
	log := logger.LoggerFromCtx(ctx)
	for _, toolCall := range toolCalls {
		if toolCall.Type != "function" {
			log.Warn("unsupported tool call type", slog.String("type", toolCall.Type))
			continue
		}
		tool, err := tools.DefaultRegistry.Get(toolCall.Function.Name)
		if err != nil {
			log.Warn("tool not found", slog.String("tool", toolCall.Function.Name))
			msgs = append(msgs, openai.ToolMessage(
				fmt.Sprintf(`{"error":"tool not found","tool":"%s"}`, toolCall.Function.Name),
				toolCall.ID,
			))
			continue
		}
		result, err := tool.Handler(ctx, jsoniter.RawMessage(toolCall.Function.Arguments))
		if err != nil {
			log.Warn("tool handler failed", slog.String("tool", toolCall.Function.Name), slog.String("error", err.Error()))
			msgs = append(msgs, openai.ToolMessage(
				fmt.Sprintf(`{"error":"%s","tool":"%s"}`, err.Error(), toolCall.Function.Name),
				toolCall.ID,
			))
			continue
		}
		resultStr, err := toolResultToString(result)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, openai.ToolMessage(resultStr, toolCall.ID))
	}
	return msgs, nil
}

func toolResultToString(result any) (string, error) {
	if result == nil {
		return "null", nil
	}
	if s, ok := result.(string); ok {
		return s, nil
	}
	encoded, err := jsoniter.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
