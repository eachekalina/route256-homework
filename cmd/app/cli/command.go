package cli

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/app/core"
	"homework/internal/app/logger"
	"homework/internal/app/pickuppoint"
	"homework/internal/app/rwthread"
)

type PickUpPointCommands struct {
	svc    *core.PickUpPointCoreService
	log    *logger.Logger
	runner *rwthread.Runner
}

func NewPickUpPointCommands(svc *core.PickUpPointCoreService, log *logger.Logger, runner *rwthread.Runner) *PickUpPointCommands {
	return &PickUpPointCommands{
		svc:    svc,
		log:    log,
		runner: runner,
	}
}

func (c *PickUpPointCommands) HelpCommand(ui *ConsoleUi) error {
	fmt.Println(`Available commands:

help	show this help
exit	exit program
create	create a new pick-up point
list	list all existing pick-up points
get		show selected pick-up point
update	update pick-up point
delete	delete pick-up point`)
	return nil
}

func (c *PickUpPointCommands) ExitCommand(ui *ConsoleUi) error {
	return ErrExit
}

func (c *PickUpPointCommands) handleInputError(err error) error {
	if errors.Is(err, ErrInvalidInput) {
		c.log.Log("%v", err)
		return nil
	}
	return err
}

func (c *PickUpPointCommands) CreateCommand(ui *ConsoleUi) error {
	var req core.CreatePointRequest
	var err error
	req.Id, err = ui.GetUint("enter point id")
	if err != nil {
		return c.handleInputError(err)
	}
	req.Name, err = ui.GetLine("enter point name")
	if err != nil {
		return c.handleInputError(err)
	}
	req.Address, err = ui.GetLine("enter point address")
	if err != nil {
		return c.handleInputError(err)
	}
	req.Contact, err = ui.GetLine("enter point contact")
	if err != nil {
		return c.handleInputError(err)
	}

	c.runner.RunWrite(func(ctx context.Context) error {
		return c.svc.CreatePoint(ctx, req)
	})
	return nil
}

func (c *PickUpPointCommands) ListCommand(ui *ConsoleUi) error {
	c.runner.RunRead(func(ctx context.Context) error {
		points, err := c.svc.ListPoints(ctx)
		if err != nil {
			return err
		}
		c.log.Log(pickuppoint.ListPoints(points))
		return nil
	})
	return nil
}

func (c *PickUpPointCommands) GetCommand(ui *ConsoleUi) error {
	id, err := ui.GetUint("enter point id")
	if err != nil {
		return c.handleInputError(err)
	}
	c.runner.RunRead(func(ctx context.Context) error {
		point, err := c.svc.GetPoint(ctx, id)
		if err != nil {
			return err
		}
		c.log.Log(pickuppoint.ListPoints([]pickuppoint.PickUpPoint{point}))
		return nil
	})
	return nil
}

func (c *PickUpPointCommands) UpdateCommand(ui *ConsoleUi) error {
	var req core.UpdatePointRequest
	var err error
	req.Id, err = ui.GetUint("enter point id")
	if err != nil {
		return c.handleInputError(err)
	}
	req.Name, err = ui.GetLine("enter new point name")
	if err != nil {
		return c.handleInputError(err)
	}
	req.Address, err = ui.GetLine("enter new point address")
	if err != nil {
		return c.handleInputError(err)
	}
	req.Contact, err = ui.GetLine("enter new point contact")
	if err != nil {
		return c.handleInputError(err)
	}

	c.runner.RunWrite(func(ctx context.Context) error {
		return c.svc.UpdatePoint(ctx, req)
	})
	return nil
}

func (c *PickUpPointCommands) DeleteCommand(ui *ConsoleUi) error {
	id, err := ui.GetUint("enter point id")
	if err != nil {
		return c.handleInputError(err)
	}
	c.runner.RunWrite(func(ctx context.Context) error {
		return c.svc.DeletePoint(ctx, id)
	})
	return nil
}
