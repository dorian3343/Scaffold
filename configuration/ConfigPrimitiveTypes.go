package configuration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/metalim/jsonmap"
	"github.com/rs/zerolog/log"
	"service/components/controller"
	model2 "service/components/model"
)

// ConfigPrimitiveTypes should include all the non-critical structs and their methods

// Database types
type database struct {
	InitQuery string `yaml:"init-query"`
	Path      string `yaml:"path"`
}

type Database struct {
	Db        *sql.DB
	InitQuery string
	Path      string
}

// Attach a db pointer
func (d database) adapt(db *sql.DB) *Database {
	return &Database{Db: db, InitQuery: d.InitQuery, Path: d.Path}
}

// Controller types
type Controller struct {
	Fallback interface{} `yaml:"fallback"`
	Name     string      `yaml:"name"`
	Model    string      `yaml:"model"`
	Cors     string      `yaml:"cors"`
}

func (c Controller) adapt(model *model2.Model) controller.Controller {
	JSON, err := json.Marshal(c.Fallback)
	if err != nil {
		log.Fatal().Err(err).Msg("JSON error in Controller : " + c.Name)
	}

	return controller.Create(c.Name, model, JSON, c.Cors)
}

// Struct representing a single field of a json spec
type JsonSpecSkeleton struct {
	Name string `yaml:"Name"`
	Type string `yaml:"Type"`
}

type model struct {
	QueryTemplate string             `yaml:"query-template"`
	JsonTemplate  []JsonSpecSkeleton `yaml:"json-template"`
	Name          string
}

// Adapt the fake model type into the actual model which has access to the database
func (m model) adapt(db *sql.DB) model2.Model {
	f := jsonmap.New()
	for i := 0; i < len(m.JsonTemplate); i++ {
		if m.JsonTemplate[i].Type == "" || m.JsonTemplate[i].Name == "" {
			log.Fatal().Err(errors.New("missing Value in Json Template")).Msg("Something went wrong with JSON template's ")
		} else {
			f.Set(m.JsonTemplate[i].Name, m.JsonTemplate[i].Type)
		}
	}
	return model2.Create(m.Name, db, m.QueryTemplate, f)
}
