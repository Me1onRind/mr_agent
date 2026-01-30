package llm

import (
	"context"

	"github.com/openai/openai-go/v3"
)

type simpleCache struct {
	data map[string][]openai.ChatCompletionMessageParamUnion
}

func newSimpleCache() *simpleCache {
	return &simpleCache{
		data: make(map[string][]openai.ChatCompletionMessageParamUnion),
	}
}

func (s *simpleCache) GetMessages(ctx context.Context, sessionId string) ([]openai.ChatCompletionMessageParamUnion, error) {
	messages, ok := s.data[sessionId]
	if !ok {
		return nil, ErrMessagesNotFound
	}
	return messages, nil
}

func (s *simpleCache) AppendMsgAndAss(ctx context.Context, sessionId string, msg, ass openai.ChatCompletionMessageParamUnion) error {
	s.data[sessionId] = append(s.data[sessionId], msg, ass)
	return nil
}
