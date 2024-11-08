package cmd

import (
	"flag"
	"github.com/imMohika/gohangyourself/cmd/sub"
	"github.com/imMohika/gohangyourself/cmd/sub/config"
	"github.com/imMohika/gohangyourself/cmd/sub/platform"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin"
	"github.com/imMohika/gohangyourself/cmd/sub/script"
	"github.com/imMohika/gohangyourself/log"
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

var subCommands = map[string]sub.Command{
	"platform": platform.SubCommand{},
	"script":   script.SubCommand{},
	"plugin":   plugin.SubCommand{},
	"config":   config.SubCommand{},
}

func Execute() {
	flag.Usage = func() {
		log.Info(usage)
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
