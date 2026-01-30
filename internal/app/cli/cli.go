package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/Me1onRind/mr_agent/internal/domain/dialog"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/cache/llm"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/session"
	"github.com/Me1onRind/mr_agent/internal/initialize"
)

type CLIService struct {
	DialogDomain *dialog.DialogDomain
	params       *Params
}

func NewCLIService() *CLIService {
	cli := &CLIService{
		DialogDomain: dialog.NewDialogDomain(),
	}
	return cli
}

func (c *CLIService) Init(ctx context.Context) *CLIService {
	logger.InitLogger(os.Stdout, slog.LevelError, false)
	_ = initialize.InitOpentracing("mr_agent", "0.0.1")(ctx)
	_ = llm.InitLLMCache(ctx)
	_ = session.InitSessionStore(ctx)
	c.params = parseCliParams()
	return c
}

func (c *CLIService) Run(ctx context.Context) error {
	log := logger.CtxLogger(ctx)
	ctx, err := session.NewSession(ctx, &session.Data{})
	if err != nil {
		return err
	}

	fmt.Printf("Hello, here is dialog with llm, dialog mode:%s\n", c.params.Mode.Name())
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF && strings.TrimSpace(input) == "" {
				break
			}
			if err != io.EOF {
				log.Error("read input failed", slog.String("error", err.Error()))
				break
			}
		}
		input = strings.TrimSpace(input)
		if input == "" {
			if err == io.EOF {
				break
			}
			continue
		}
		reply, err := c.DialogDomain.Dialog(ctx, input, c.params.WithContext())
		if err != nil {
			log.Error("c.DialogDomain.Dialog failed", slog.String("error", err.Error()))
			break
		}
		fmt.Printf("\nLLM: %s\n\n", reply)
	}
	return nil
}
