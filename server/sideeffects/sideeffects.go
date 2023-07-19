package sideeffects

import (
	"moco/config"
	"moco/server/sideeffects/effects"
	"moco/server/sideeffects/plugins"
	"moco/server/sideeffects/utils"
	"moco/server/structs"
	"moco/utils/logger"
	"net/http"
)

var availableEffects = make(map[string]utils.EffectHandler)
var log = logger.New("sideeffects")

func Init() {
	log.Info("Initializing main side effects")
	availableEffects["create"] = effects.CreateEffect
	availableEffects["retrieve"] = effects.RetrieveEffect

	plugins := plugins.GetPlugins()

	// register plugin side effects
	for _, plugin := range plugins {
		log.Infof("Registering plugin side effect: %s", plugin.Name)
		availableEffects[plugin.Name] = plugin.Handle
	}
}

func Handle(
	w http.ResponseWriter,
	r *http.Request,
	request structs.Request,
	response config.ResponseConfig,
	sideEffects []config.SideEffectConfig,
	defaultHandler utils.DefaultRequestHandler,
) {
	for _, sideEffect := range sideEffects {
		// if sideEffect.Type starts with "plugin:", we remove the "plugin:" part
		// and call the plugin side effect

		var pluginName string

		if len(sideEffect.Type) > 7 && sideEffect.Type[:7] == "plugin:" {
			pluginName = sideEffect.Type[7:]
		} else {
			pluginName = sideEffect.Type
		}

		// if sideEffect.Type exists in effects, call it
		if effect, ok := availableEffects[pluginName]; ok {
			effect(w, r, request, response, sideEffect, defaultHandler)
			return
		} else {
			log.Warnf("Unknown side effect type: %s", pluginName)
		}
	}

	// if no side effect was called, call default handler
	log.Warn("No side effect was called, even though there were side effects configured")
	defaultHandler(w, r, response)
}
