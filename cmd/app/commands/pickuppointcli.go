package commands

import (
	"context"
	"golang.org/x/sync/errgroup"
	"homework/cmd/app/cli"
	"homework/internal/app/core"
	"homework/internal/app/logger"
	"homework/internal/app/rwthread"
	"os/signal"
	"syscall"
)

type PickUpPointCliConsoleCommands struct {
	svc  core.PickUpPointCoreService
	help Command
}

func NewPickUpPointCliConsoleCommands(svc core.PickUpPointCoreService, help Command) *PickUpPointCliConsoleCommands {
	return &PickUpPointCliConsoleCommands{svc: svc, help: help}
}

func (c *PickUpPointCliConsoleCommands) ManagePickUpPointsCommand(args []string) error {
	fs := createFlagSet(c.help)
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(ctx)

	logCtx, stopLog := context.WithCancel(context.Background())
	defer stopLog()

	log := logger.NewLogger()
	go log.Run(logCtx)

	runner := rwthread.NewRunner(log)
	eg.Go(func() error {
		return runner.Run(ctx)
	})

	cmds := cli.NewPickUpPointCommands(c.svc, log, runner)
	cmdMap := map[string]cli.Command{
		"help":   cmds.HelpCommand,
		"exit":   cmds.ExitCommand,
		"create": cmds.CreateCommand,
		"list":   cmds.ListCommand,
		"get":    cmds.GetCommand,
		"update": cmds.UpdateCommand,
		"delete": cmds.DeleteCommand,
	}

	ui := cli.NewConsoleUi(cmdMap)
	eg.Go(func() error {
		return ui.Run(ctx)
	})

	return eg.Wait()
}
