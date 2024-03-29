package cmdargs

import (
	"errors"
	"os"
)

type Command func(args []string) error

func Run(commands map[string]Command) error {
	if len(os.Args) < 2 {
		return errors.New("no subcommand specified, see `help` subcommand for details")
	}

	cmdName := os.Args[1]
	args := os.Args[2:]

	cmd, ok := commands[cmdName]
	if !ok {
		return errors.New("not such subcommand, see `help` subcommand for details")
	}

	return cmd(args)
}
