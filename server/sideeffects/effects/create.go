package effects

import (
	"moco/config"
	db "moco/server/database"
	"moco/server/sideeffects/utils"
	"moco/server/structs"
	"net/http"
)

func CreateEffect(
	w http.ResponseWriter,
	r *http.Request,
	_ structs.Request,
	response config.ResponseConfig,
	cfg config.SideEffectConfig,
	defaultHandler utils.DefaultRequestHandler,
) {
	if response.Request.IsLoggedIn && r.Header.Get("Authorization") == "" {
		defaultHandler(w, r, response)
	}

	// insert the model in the database
	if cfg.Entity != "" {
		model := db.Insert(cfg.Entity, cfg.Data)

		// create the response
		newResponse := config.ResponseConfig{
			Request:         response.Request,
			Status:          http.StatusCreated,
			ResponseBody:    model,
			ResponseHeaders: response.ResponseHeaders,
		}

		// call the default handler
		defaultHandler(w, r, newResponse)
	}

}
