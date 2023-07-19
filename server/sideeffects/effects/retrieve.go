package effects

import (
	"fmt"
	"moco/config"
	db "moco/server/database"
	"moco/server/sideeffects/utils"
	"moco/server/structs"
	"net/http"
)

func RetrieveEffect(
	w http.ResponseWriter,
	r *http.Request,
	request structs.Request,
	response config.ResponseConfig,
	cfg config.SideEffectConfig,
	defaultHandler utils.DefaultRequestHandler,
) {
	if response.Request.IsLoggedIn && r.Header.Get("Authorization") == "" {
		defaultHandler(w, r, response)
	}

	searchBy := cfg.Data["searchBy"].(string)
	searchFrom := cfg.Data["searchFrom"].(string)

	if cfg.Entity == "" {
		http.Error(w, "RetrieveEffect: entity must be defined", http.StatusInternalServerError)
		return
	}

	var identifier any

	// get value to search by using searchFrom and searchBy
	if searchFrom == "params" {
		identifier = request.Params[searchBy]
	} else if searchFrom == "body" {
		identifier = request.Body[searchBy]
	} else if searchFrom == "query" {
		identifier = request.Query[searchBy]
	} else {
		http.Error(w, "ChatAppendEffect: searchFrom must be one of params, body or query", http.StatusInternalServerError)
		return
	}

	// insert the model in the database
	if cfg.Entity != "" {
		model := db.GetOne(cfg.Entity, func(m db.Model) bool {
			return string(fmt.Sprintf("%v", m[searchBy])) == identifier
		})

		if model == nil {
			http.Error(w, fmt.Sprintf("No model found with %s = %s", searchBy, identifier), http.StatusNotFound)
			return
		}

		// create the response
		newResponse := config.ResponseConfig{
			Request:         response.Request,
			Status:          http.StatusCreated,
			ResponseBody:    *model,
			ResponseHeaders: response.ResponseHeaders,
		}

		// call the default handler
		defaultHandler(w, r, newResponse)
	}

}
