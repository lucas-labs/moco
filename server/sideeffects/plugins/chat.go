package plugins

import (
	"fmt"
	"moco/config"
	db "moco/server/database"
	"moco/server/sideeffects/utils"
	"moco/server/structs"
	"net/http"
	"regexp"
	"time"
)

var humanMessageId = 1
var aiMessageId = 1

func ChatAppendEffect(
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

	appendTo := cfg.Data["appendTo"].(string)
	modelToAppend := db.ModelDeepCopy(cfg.Data["model"].(map[string]interface{}))
	searchBy := cfg.Data["searchBy"].(string)
	searchFrom := cfg.Data["searchFrom"].(string)

	// check that appendTo, searchFrom, searchBy and modelToAppend are not empty
	if appendTo == "" || searchFrom == "" || searchBy == "" || len(modelToAppend) == 0 {
		http.Error(w, "ChatAppendEffect: appendTo, searchFrom, searchBy and modelToAppend must be set", http.StatusInternalServerError)
		return
	}

	// check that the model to append exists
	if !db.ModelTypeExists(cfg.Entity) {
		http.Error(w, "ChatAppendEffect: model type "+cfg.Entity+" does not exist", http.StatusInternalServerError)
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

	// get the model to append to
	model := db.GetOne(cfg.Entity, func(model db.Model) bool {
		return string(fmt.Sprintf("%v", model[searchBy])) == identifier
	})

	if model == nil {
		http.Error(w, "ChatAppendEffect: no model found with "+searchBy+" = "+identifier.(string), http.StatusInternalServerError)
		return
	}

	// check that model has a field called like appendTo and that it is a slice
	if _, ok := (*model)[appendTo]; !ok {
		// create the field if it does not exist
		(*model)[appendTo] = make([]interface{}, 0)
	}

	if _, ok := (*model)[appendTo].([]interface{}); !ok {
		http.Error(w, "ChatAppendEffect: model."+appendTo+" is not a slice", http.StatusInternalServerError)
		return
	}

	compiledModel := compileModel(modelToAppend, request)
	compiledModel["type"] = "human"

	// append modelToAppend to the slice
	(*model)[appendTo] = append(
		(*model)[appendTo].([]interface{}),
		compiledModel,
	)

	// append another message, representing the ai's response
	// as we are a mock, lets randomly choose a response using utils.GetRandomParagraph()
	aiResponse := make(map[string]interface{})
	aiResponse["id"] = aiMessageId
	aiResponse["content"] = utils.GetRandomParagraph()
	aiResponse["created"] = time.Now().Format("2006-01-02T15:04:05.999999+00:00")
	aiResponse["type"] = "ai"
	aiMessageId++

	// append the response to the slice
	(*model)[appendTo] = append(
		(*model)[appendTo].([]interface{}),
		aiResponse,
	)

	newResponse := config.ResponseConfig{
		Request:         response.Request,
		Status:          http.StatusCreated,
		ResponseBody:    aiResponse,
		ResponseHeaders: response.ResponseHeaders,
	}

	defaultHandler(w, r, newResponse)
}

func compileModel(model db.Model, request structs.Request) db.Model {
	// regex that matches _FROM_BODY_:<field_name>
	fromBodyRegex := regexp.MustCompile(`_FROM_BODY_:([a-zA-Z0-9_]+)`)

	for key, value := range model {
		val := fmt.Sprintf("%v", value)

		if val == "_AUTONUM_" {
			model[key] = humanMessageId
			humanMessageId++
		} else if val == "_NOW_" {
			model[key] = time.Now().Format("2006-01-02T15:04:05.999999+00:00")
		} else if fromBodyRegex.MatchString(val) {
			// get the field name
			fieldName := fromBodyRegex.FindStringSubmatch(val)[1]
			val := request.Body[fieldName]
			if val == nil {
				log.Error("ChatAppendEffect: field " + fieldName + " not found in request body")
				return model
			}
			model[key] = val
		}
	}

	return model
}
