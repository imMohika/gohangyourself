package internal

import (
	"errors"
	"strings"
)

func GetHandlerFromURL(url string) (string, error) {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	switch {
	case strings.HasPrefix(url, "hangar.papermc.io"):
		return "hangar", nil
	case strings.HasPrefix(url, "modrinth.com"):
		return "modrinth", nil
	case strings.HasPrefix(url, "spigotmc.org"):
		return "spigot", nil
	case strings.HasPrefix(url, "dev.bukkit.org"):
		return "bukkit", nil
	case strings.HasPrefix(url, "bukkit.org"):
		return "bukkit", nil
	}

	return "", errors.New("invalid url")
}

//func GetPluginInfo(name string) PluginInfo {
//	var info PluginInfo
//
//	switch {
//	// If it's a url using http then use this:
//	case strings.HasPrefix(name, "https://") || strings.HasPrefix(name, "http://"):
//		getURL(name, &info)
//
//	// If it's file which ends in .json try reading it locally:
//	case strings.HasSuffix(name, ".json"):
//		getLocal(name, &info)
//
//	// If it's a modrinth plugin try getting it from modrinth:
//	case strings.HasPrefix(name, "modrinth:"):
//		getModrinth(name, &info)
//
//	// If it's a spigot plugin try getting it from spigotmc:
//	case strings.HasPrefix(name, "spigot:"),
//		strings.HasPrefix(name, "spigotmc:"):
//		getSpigotmc(name, &info)
//
//	// If it's a bukkit plugin try getting it from bukkit:
//	case strings.HasPrefix(name, "bukkit:"),
//		strings.HasPrefix(name, "bukkitdev:"):
//		getBukkitdev(name, &info)
//
//	// If it's none of the options above try getting it from the repos:
//	default:
//		getRepos(name, &info)
//	}
//
//	return info
//}
