package dialog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/cache/llm"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/client/agent"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/session"
	"github.com/openai/openai-go/v3"
)

type DialogDomain struct {
}

func NewDialogDomain() *DialogDomain {
	d := &DialogDomain{}
	return d
}

func (d *DialogDomain) Dialog(ctx context.Context, msg string, withCtx bool) (string, error) {
	log := logger.CtxLogger(ctx)
	userMsg := openai.UserMessage(msg)
	agentClient := agent.NewAgentClient(ctx)

	var (
		msgs []openai.ChatCompletionMessageParamUnion
		err  error
	)

	sessionId, err := session.GetSessionId(ctx)
	if err != nil {
		return "", err
	}

	if withCtx {
		msgs, err = llm.GetMessages(ctx, sessionId)
		if err != nil && !errors.Is(err, llm.ErrMessagesNotFound) {
			log.Error("LLMCache.GetMessages failed", slog.String("error", err.Error()))
			return "", err
		}
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

	if withCtx {
		if err := llm.AppendMsgAndAss(ctx, sessionId, userMsg, assistant); err != nil {
			log.Warn("LLMCache.AppendMsgAndAss failed", slog.String("error", err.Error()))
		}
	}

	return reply, nil
}
