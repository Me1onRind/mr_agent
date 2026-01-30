package dialog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/cache/llm"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/client/agent"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/openai/openai-go/v3"
)

type DialogDomain struct {
	LLMCache llm.LLMCache
}

func NewDialogDomain() *DialogDomain {
	d := &DialogDomain{
		LLMCache: llm.NewSimpleCache(),
	}
	return d
}

func (d *DialogDomain) Dialog(ctx context.Context, msg string) (string, error) {
	log := logger.CtxLogger(ctx)
	userMsg := openai.UserMessage(msg)
	agentClient := agent.NewAgentClient(ctx)

	msgs, err := d.LLMCache.GetMessages(ctx, "session_id")
	if err != nil && !errors.Is(err, llm.ErrMessagesNotFound) {
		log.Error("LLMCache.GetMessages failed", slog.String("error", err.Error()))
		return "", err
	}

	msgs = append(msgs, userMsg)
	completion, err := agentClient.Chat(ctx, msgs)
	if err != nil {
		return "", err
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("chat completion has no choices")
	}

	reply := completion.Choices[0].Message.Content
	assistant := openai.AssistantMessage(reply)

	if err := d.LLMCache.AppendMsgAndAss(ctx, "session_id", userMsg, assistant); err != nil {
		log.Warn("LLMCache.AppendMsgAndAss failed", slog.String("error", err.Error()))
	}

	return reply, nil
}
