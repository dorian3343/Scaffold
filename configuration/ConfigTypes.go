package configuration

import (
	"database/sql"
	"errors"
	"github.com/metalim/jsonmap"
	"github.com/rs/zerolog/log"
	model2 "service/model"
)

// ConfigTypes should include all the non-critical structs and their methods

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
	CORS     bool        `yaml:"CORS"`
}

// Struct representing a single field of a json spec
type TypeSkeleton struct {
	Name string `yaml:"Name"`
	Type string `yaml:"Type"`
}

type model struct {
	QueryTemplate string         `yaml:"query-template"`
	JsonTemplate  []TypeSkeleton `yaml:"json-template"`
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
