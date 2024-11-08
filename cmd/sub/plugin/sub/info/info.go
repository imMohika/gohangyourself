package info

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/handler"
	"github.com/imMohika/gohangyourself/log"
	"strings"
)

type SubCommand struct {
}

func (s SubCommand) Handle(args []string) {
	if len(args) != 1 {
		log.Warn("wrong usage. please provide url")
		return
	}

	url := strings.TrimSpace(args[0])
	pluginHandler, err := handler.GetHandler(url)
	if err != nil {
		log.Fatal(err, "can not get handler", "url", url)
		return
	}

	meta, err := pluginHandler.GetMeta()
	if err != nil {
		log.Fatal(err, "can not get plugin meta", "url", url)
		return
	}

	log.Info(meta.Title)

	if meta.Description != "" {
		log.Info(meta.Description)
	}

	if len(meta.Loaders) > 0 {
		log.Info("Loaders: " + strings.Join(meta.Loaders, ", "))
	}

	log.Info("Last Update: " + humanize.Time(meta.Updated))

	if meta.Downloads > 0 {
		log.Info(fmt.Sprintf("Downloads: %s", humanize.Comma(int64(meta.Downloads))))
	}

	if meta.Source != "" {
		log.Info(fmt.Sprintf("Source: %s", meta.Source))
	}

	if meta.Support != "" {
		log.Info(fmt.Sprintf("Support: %s", meta.Support))
	}

	if meta.Wiki != "" {
		log.Info(fmt.Sprintf("Wiki: %s", meta.Wiki))
	}
}
