package goroutine

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
)

func LogPanicStack(ctx context.Context, err any) {
	stack := debug.Stack()
	logger.CtxLogger(ctx).Error("panic", slog.Any("error", err), slog.String("stack", string(stack)))
	fmt.Println(err)
	fmt.Println(string(stack))
}

func SafeGo[T any](ctx context.Context, f func(c context.Context, args T), args T) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				LogPanicStack(ctx, err)
			}
		}()
		f(ctx, args)
	}()
}
