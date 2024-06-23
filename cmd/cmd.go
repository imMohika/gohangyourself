package cmd

import (
	"flag"
	"gohangyourself/cmd/sub/platform"
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
	switch subCmd {
	case "version":
		slog.Info("version")
	case "plugin":
		slog.Info("plugin")
	case "script":
		slog.Info("script")
	case "help":
		slog.Info("help")
	case "platform":
		platform.Handle(os.Args[2:])
	default:
		flag.Usage()
	}
}
