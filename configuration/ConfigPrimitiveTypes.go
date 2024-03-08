package configuration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metalim/jsonmap"
	"service/components/controller"
	model2 "service/components/model"
	"strings"
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
	Cache    string      `yaml:"cache"`
	Verb     string      `yaml:"verb"`
}

func (c Controller) adapt(model *model2.Model) (controller.Controller, error) {
	JSON, err := json.Marshal(c.Fallback)
	if err != nil {
		return controller.Controller{}, errors.New(fmt.Sprintf("Json error in Controller : %s", c.Name))
	}
	verb := strings.ToUpper(c.Verb)
	switch verb {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD", "":
	default:
		err := errors.New("Unrecognized HTTP method")
		fmt.Println(err)
		return controller.Controller{}, err
	}
	return controller.Create(c.Name, model, JSON, c.Cors, c.Cache, verb), nil
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
func (m model) adapt(db *sql.DB) (model2.Model, error) {
	f := jsonmap.New()
	for i := 0; i < len(m.JsonTemplate); i++ {
		if m.JsonTemplate[i].Type == "" || m.JsonTemplate[i].Name == "" {
			return model2.Model{}, errors.New("missing Value or Type in Json Template field")
		} else {
			f.Set(m.JsonTemplate[i].Name, m.JsonTemplate[i].Type)
		}
	}
	return model2.Create(m.Name, db, m.QueryTemplate, f), nil
}
