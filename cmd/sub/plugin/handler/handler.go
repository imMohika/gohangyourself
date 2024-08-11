package handler

import "github.com/imMohika/gohangyourself/cmd/sub/plugin/internal"

//type PluginInfo struct {
//	url string
//	id  string
//}

type PluginHandler interface {
	//Handle(args []string)
	FromURL(url string) (internal.PluginInfo, error)
}
