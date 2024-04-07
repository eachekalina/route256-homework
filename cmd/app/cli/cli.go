package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Command func(ui *ConsoleUi) error

var ErrInvalidInput = errors.New("invalid input")
var ErrExit = errors.New("exiting")

type ConsoleUi struct {
	scanner  *bufio.Scanner
	commands map[string]Command
}

func NewConsoleUi(commands map[string]Command) *ConsoleUi {
	return &ConsoleUi{scanner: bufio.NewScanner(os.Stdin), commands: commands}
}

func (ui *ConsoleUi) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := ui.handleCommand()
			if err != nil {
				return err
			}
		}
	}
}

func (ui *ConsoleUi) handleCommand() error {
	line, err := ui.getLine("")
	if err != nil {
		return err
	}
	cmd, ok := ui.commands[line]
	if !ok {
		fmt.Println("no such command found, use `help` for help")
		return nil
	}
	return cmd(ui)
}

func (ui *ConsoleUi) getLine(prompt string) (string, error) {
	fmt.Printf("%s > ", prompt)
	if !ui.scanner.Scan() {
		return "", ui.scanner.Err()
	}
	line := ui.scanner.Text()
	return strings.TrimSpace(line), nil
}

func (ui *ConsoleUi) getUint(prompt string) (uint64, error) {
	str, err := ui.getLine(prompt)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, ErrInvalidInput
	}
	return n, nil
}
