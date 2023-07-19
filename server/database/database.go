package database

import (
	"moco/utils/logger"
	"time"
)

var log = logger.New("database")

type Model = map[string]interface{}

type Database struct {
	models map[string][]Model
}

var db *Database

// New creates a new database
func Init() {
	log.Info("Initializing database")

	db = &Database{
		models: make(map[string][]Model),
	}
}

// GetClauseFn is a function that returns true if the given model matches the clause
type GetClauseFn func(model Model) bool

// ModelTypeExists returns true if the given model type exists
func ModelTypeExists(modelName string) bool {
	return db.models[modelName] != nil
}

// GetAll returns all models of the given type
func GetAll(modelName string) []Model {
	models := db.models[modelName]

	if models == nil {
		log.Warnf("No models found for %s", modelName)
		return []Model{}
	}

	return models
}

// GetOne returns the first model of the given type that matches the given clause
func GetOne(modelName string, where GetClauseFn) *Model {
	for _, model := range db.models[modelName] {
		if where(model) {
			return &model
		}
	}

	return nil
}

// GetMany returns all models of the given type that match the given clause
func GetMany(modelName string, where GetClauseFn) []*Model {
	var result []*Model

	for _, model := range db.models[modelName] {
		if where(model) {
			result = append(result, &model)
		}
	}

	return result
}

// Insert adds a new model to the database
func Insert(modelName string, rawModel Model) Model {
	cpy := ModelDeepCopy(rawModel)

	model := processModel(cpy, modelName)
	db.models[modelName] = append(db.models[modelName], model)

	return model
}

func ModelDeepCopy(model Model) Model {
	newModel := make(Model)

	for key, value := range model {
		if subModel, ok := value.(Model); ok {
			newModel[key] = ModelDeepCopy(subModel)
		} else {
			newModel[key] = value
		}
	}

	return newModel
}

// processModel processes the given model searching for special values
// like _AUTONUM_ and _NOW_ and replaces them with the appropriate values
func processModel(model Model, modelName string) Model {
	for key, value := range model {
		if value == "_AUTONUM_" {
			model[key] = len(db.models[modelName]) + 1
		} else if value == "_NOW_" {
			// 2023-06-29T14:07:23.416666+00:00
			model[key] = time.Now().Format("2006-01-02T15:04:05.999999+00:00")
		}
	}

	return model
}

// Update updates all models of the given type that match the given clause
func Update(modelName string, model Model, where GetClauseFn) ([]Model, int) {
	var updatedModels []Model = []Model{}

	for i, m := range db.models[modelName] {
		if where(m) {
			// TODO: merge instead of replace
			db.models[modelName][i] = model
			updatedModels = append(updatedModels, model)
		}
	}

	return updatedModels, len(updatedModels)
}

// Delete deletes all models of the given type that match the given clause
// returns the number of deleted models
func Delete(modelName string, where GetClauseFn) int {
	var deletedModels int = 0

	for i := len(db.models[modelName]) - 1; i >= 0; i-- {
		if where(db.models[modelName][i]) {
			db.models[modelName] = append(db.models[modelName][:i], db.models[modelName][i+1:]...)
			deletedModels++
		}
	}

	return deletedModels
}

// DeleteOne deletes the first model of the given type that matches the given clause
// returns the number of deleted models (0 or 1)
func DeleteOne(modelName string, where GetClauseFn) int {
	for i, model := range db.models[modelName] {
		if where(model) {
			db.models[modelName] = append(db.models[modelName][:i], db.models[modelName][i+1:]...)
			return 1
		}
	}

	return 0
}

// Reset clears the database
func Reset() {
	db.models = make(map[string][]Model)
}

// ResetModel clears all models of the given type
func ResetModel(modelName string) {
	db.models[modelName] = []Model{}
}

// Dump returns a map of all models in the database
func Dump() map[string][]Model {
	return db.models
}

// Load loads the given models into the database
func Load(models map[string][]Model) {
	db.models = models
}
