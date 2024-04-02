package commands

import "flag"

func createFlagSet(help Command) *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.Usage = func() {
		help(nil)
	}
	return fs
}
