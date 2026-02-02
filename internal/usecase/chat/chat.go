package chat

import (
	"context"

	"github.com/Me1onRind/mr_agent/internal/domain/dialog"
	"github.com/Me1onRind/mr_agent/internal/protocol/http/chat"
)

type ChatUsecase struct {
	DialogDomain *dialog.DialogDomain
}

func NewChatUsecase() *ChatUsecase {
	return &ChatUsecase{
		DialogDomain: dialog.NewDialogDomain(),
	}
}

func (a *ChatUsecase) Chat(ctx context.Context, request *chat.ChatRequest) (*chat.ChatResponse, error) {
	reply, err := a.DialogDomain.Dialog(ctx, request.Msg, &dialog.DialogConfig{
		WithCtx: true,
	})
	if err != nil {
		return nil, err
	}

	return &chat.ChatResponse{
		Msg: reply,
	}, nil
}
