package llm

import (
	"context"
	"errors"

	"github.com/openai/openai-go/v3"
)

var llmCache llmCacheIface

var (
	ErrStoreNotInit     = errors.New("session message store not init")
	ErrMessagesNotFound = errors.New("session message not found")
)

func InitLLMCache(ctx context.Context) error {
	llmCache = newSimpleCache()
	return nil
}

func GetMessages(ctx context.Context, seesionId string) ([]openai.ChatCompletionMessageParamUnion, error) {
	if llmCache == nil {
		return nil, ErrStoreNotInit
	}
	return llmCache.GetMessages(ctx, seesionId)
}

func AppendMsgAndAss(ctx context.Context, sessionId string, msg, ass openai.ChatCompletionMessageParamUnion) error {
	if llmCache == nil {
		return ErrStoreNotInit
	}
	return llmCache.AppendMsgAndAss(ctx, sessionId, msg, ass)
}

type llmCacheIface interface {
	GetMessages(ctx context.Context, seesionId string) ([]openai.ChatCompletionMessageParamUnion, error)
	AppendMsgAndAss(ctx context.Context, sessionId string, msg, ass openai.ChatCompletionMessageParamUnion) error
}
