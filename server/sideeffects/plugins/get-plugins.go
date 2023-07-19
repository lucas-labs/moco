package plugins

import (
	"moco/server/sideeffects/utils"
	"moco/utils/logger"
)

var log = logger.New("plugins")

type Plugin struct {
	Name   string
	Handle utils.EffectHandler
}

func GetPlugins() []Plugin {
	// TODO: the idea is that in the future, we can get
	// the plugins from a /plugins folder
	// loading them dynamically (maybe .dll files)
	// for now, we just return the plugins we have
	// locally inside this package
	log.Info("Loading plugins")

	return []Plugin{
		{
			Name:   "chat:append",
			Handle: ChatAppendEffect,
		},
	}
}
