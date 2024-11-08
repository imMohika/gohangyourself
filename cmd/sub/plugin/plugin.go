package plugin

import (
	"errors"
	"flag"
	"github.com/imMohika/gohangyourself/cmd/sub"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/sub/download"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/sub/info"
	"github.com/imMohika/gohangyourself/log"
	"os"
	"strings"
)

type SubCommand struct {
}

var subCommands = map[string]sub.Command{
	"info":     info.SubCommand{},
	"download": download.SubCommand{},
}

const usage = `Plugin:
	info <url>	print info about the plugin
	download <url>	download the plugin

Note:
	<url> should start with "hangar:", "modrinth:" or "https://"
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
