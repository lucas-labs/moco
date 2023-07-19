package utils

import (
	"moco/config"
	"moco/server/structs"
	"net/http"
)

type DefaultRequestHandler func(
	w http.ResponseWriter,
	r *http.Request,
	response config.ResponseConfig,
) int

type EffectHandler func(
	w http.ResponseWriter,
	r *http.Request,
	request structs.Request,
	response config.ResponseConfig,
	sideEffects config.SideEffectConfig,
	defaultHandler DefaultRequestHandler,
)
