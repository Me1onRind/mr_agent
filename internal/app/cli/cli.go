package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/Me1onRind/mr_agent/internal/config"
	"github.com/Me1onRind/mr_agent/internal/domain/dialog"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/agent/tools"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/cache/llm"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/db"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/session"
	"github.com/Me1onRind/mr_agent/internal/initialize"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
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
	if err := config.LoadLocalConfig("./conf/config.yaml"); err != nil {
		panic(err)
	}

	logger.InitLogger(os.Stdout, slog.LevelError, false)
	_ = initialize.InitOpentracing("mr_agent", "0.0.1")(ctx)
	_ = llm.InitLLMCache(ctx)
	_ = session.InitSessionStore(ctx)
	_ = db.InitRegistry(ctx, config.LocalCfg.MysqlConfigs)
	if err := tools.InitAgentTools(ctx); err != nil {
		panic(err)
	}
	c.params = parseCliParams()
	return c
}

func (c *CLIService) Run(ctx context.Context) error {
	log := logger.LoggerFromCtx(ctx)
	ctx, err := session.NewSession(ctx, &session.Data{})
	if err != nil {
		return err
	}

	fmt.Printf("Hello, here is dialog with llm, dialog mode:%s\n", c.params.Mode.Name())
	yellow := color.RGB(255, 255, 102)
	blue := color.RGB(153, 255, 255)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          yellow.Sprint("You: "),
		HistoryFile:     "",
		InterruptPrompt: "",
		EOFPrompt:       "",
	})
	if err != nil {
		return err
	}
	defer rl.Close()
	for {
		input, err := rl.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Error("read input failed", slog.String("error", err.Error()))
			break
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		reply, err := c.DialogDomain.Dialog(ctx, input, &dialog.DialogConfig{
			WithCtx: c.params.WithContext(),
		})
		if err != nil {
			log.Error("c.DialogDomain.Dialog failed", slog.String("error", err.Error()))
			break
		}
		blue.Print("\nLLM: ")
		fmt.Printf("%s\n\n", reply)
	}
	return nil
}
