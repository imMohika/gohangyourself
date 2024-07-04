package cmd

import (
	"flag"
	"gohangyourself/cmd/sub/platform"
	"gohangyourself/cmd/sub/script"
	"log/slog"
	"os"
	"strings"
)

const usage = `Commands:
	platform	download a platform
	plugin		download a plugin
	script		generate a run script
	help		show help
	version		show version

Flags:
	test, t 	just for testing flags
`

var (
	TestFlag bool
)

type SubCommand interface {
	Handle(args []string)
}

var subCommands = map[string]SubCommand{
	"platform": platform.SubCommand{},
	"script":   script.SubCommand{},
}

func Execute() {
	flag.Usage = func() {
		slog.Info(usage)
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	flag.BoolVar(&TestFlag, "test", false, "test flag")
	flag.BoolVar(&TestFlag, "t", false, "test flag")

	subCmd := strings.ToLower(os.Args[1])
	handler, exists := subCommands[subCmd]
	if !exists {
		slog.Error("sub command not found", "subCmd", subCmd)
		flag.Usage()
		os.Exit(1)
	}

	handler.Handle(os.Args[2:])
}
