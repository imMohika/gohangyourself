package config

import (
	"errors"
	"flag"
	"github.com/imMohika/gohangyourself/cmd/sub"
	"github.com/imMohika/gohangyourself/cmd/sub/config/sub/cp"
	"github.com/imMohika/gohangyourself/log"
	"os"
	"strings"
)

type SubCommand struct {
}

var subCommands = map[string]sub.Command{
	"cp": cp.SubCommand{},
	//"sync": sync.SubCommand{},
}

const usage = `Config:
	cp	copy and process a config file
`

func (s SubCommand) Handle(args []string) {
	flag.Usage = func() {
		log.Info(usage)
	}

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	subCmd := strings.ToLower(args[0])
	handler, exists := subCommands[subCmd]
	if !exists {
		log.Error(errors.New("sub command not found"), "subCmd", subCmd)
		flag.Usage()
		os.Exit(1)
	}

	handler.Handle(args[1:])
}
