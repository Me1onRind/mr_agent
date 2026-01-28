package llm

import (
	"context"
	"fmt"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/openai/openai-go/v3"
)

type SimpleCache struct {
	data map[string][]openai.ChatCompletionMessageParamUnion
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{
		data: make(map[string][]openai.ChatCompletionMessageParamUnion),
	}
}

func (s *SimpleCache) GetMessages(ctx context.Context, sessionId string) ([]openai.ChatCompletionMessageParamUnion, error) {
	messages, ok := s.data[sessionId]
	if !ok {
		return nil, ErrMessagesNotFound
	}
	return messages, nil
}

func (s *SimpleCache) AppendMsgAndAss(ctx context.Context, sessionId string, msg, ass openai.ChatCompletionMessageParamUnion) error {
	if s.data == nil {
		logger.CtxLogger(ctx).Error("data is nil")
		return fmt.Errorf("SimpleCache Data Not Init")
	}
	s.data[sessionId] = append(s.data[sessionId], msg, ass)
	return nil
}
