package chat

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/cache/llm"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/client/agent"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/Me1onRind/mr_agent/internal/protocol/http/chat"
	"github.com/openai/openai-go/v3"
)

type ChatUsecase struct {
	LLMCache llm.LLMCache
}

func NewChatUsecase() *ChatUsecase {
	return &ChatUsecase{
		LLMCache: llm.NewSimpleCache(),
	}
}

func (a *ChatUsecase) Chat(ctx context.Context, request *chat.ChatRequest) (*chat.ChatResponse, error) {
	log := logger.CtxLogger(ctx)
	userMsg := openai.UserMessage(request.Msg)
	agentClient := agent.NewAgentClient(ctx)

	msgs, err := a.LLMCache.GetMessages(ctx, "session_id")
	if err != nil && !errors.Is(err, llm.ErrMessagesNotFound) {
		log.Error("LLMCache.GetMessages failed", slog.String("error", err.Error()))
		return nil, err
	}

	msgs = append(msgs, userMsg)
	completion, err := agentClient.Chat(ctx, msgs)
	if err != nil {
		return nil, err
	}

	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("chat completion has no choices")
	}

	reply := completion.Choices[0].Message.Content
	assistant := openai.AssistantMessage(reply)

	if err := a.LLMCache.AppendMsgAndAss(ctx, "session_id", userMsg, assistant); err != nil {
		log.Warn("LLMCache.AppendMsgAndAss failed", slog.String("error", err.Error()))
	}

	return &chat.ChatResponse{
		Msg: reply,
	}, nil
}
