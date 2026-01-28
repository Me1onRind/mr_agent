package llm

import (
	"context"
	"errors"

	"github.com/openai/openai-go/v3"
)

var (
	ErrMessagesNotFound = errors.New("session message not found")
)

type LLMCache interface {
	GetMessages(ctx context.Context, seesionId string) ([]openai.ChatCompletionMessageParamUnion, error)
	AppendMsgAndAss(ctx context.Context, sessionId string, msg, ass openai.ChatCompletionMessageParamUnion) error
}

type LLMCacheProxy struct {
}

func (l *LLMCacheProxy) GetMessages(ctx context.Context, seesionId string) ([]openai.ChatCompletionMessageParamUnion, error) {
	panic("not implemented") // TODO: Implement
}

func (l *LLMCacheProxy) AppendMsgAndAss(ctx context.Context, sessionId string, msg string, ass string) error {
	panic("not implemented") // TODO: Implement
}
