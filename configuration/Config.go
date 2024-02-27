package configuration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/metalim/jsonmap"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"service/controller"
	database2 "service/database"
	model2 "service/model"
)

type database struct {
	InitQuery string `yaml:"init-query"`
	Path      string `yaml:"path"`
}

type Database struct {
	Db        *sql.DB
	InitQuery string
	Path      string
}

func (d database) adapt(db *sql.DB) *Database {
	return &Database{Db: db, InitQuery: d.InitQuery, Path: d.Path}
}

type KeyValue struct {
	Name string `yaml:"Name"`
	Type string `yaml:"Type"`
}
type model struct {
	QueryTemplate string     `yaml:"query-template"`
	JsonTemplate  []KeyValue `yaml:"json-template"`
	Name          string
}

func (m model) adapt(db *sql.DB) model2.Model {
	f := jsonmap.New()
	for i := 0; i < len(m.JsonTemplate); i++ {
		if m.JsonTemplate[i].Type == "" || m.JsonTemplate[i].Name == "" {
			log.Fatal().Err(errors.New("Missing Value in Json Template")).Msg("Something went wrong with JSON template's ")
		} else {
			f.Set(m.JsonTemplate[i].Name, m.JsonTemplate[i].Type)
		}
	}
	return model2.Create(m.Name, db, m.QueryTemplate, f)
}

type Controller struct {
	Fallback interface{} `yaml:"fallback"`
	Name     string      `yaml:"name"`
	Model    string      `yaml:"model"`
	CORS     bool        `yaml:"CORS"`
}

/* Private server config to only be used for constructing the public one*/
type server struct {
	Port      int    `yaml:"port"`
	TargetLog string `yaml:"target-log"`
	Services  []struct {
		Route      string `yaml:"route"`
		Controller string `yaml:"controller"`
	} `yaml:"service(s)"`
}

type Server struct {
	Port      int
	TargetLog string
	Services  map[string]controller.Controller
}

func (s server) adapt(controllers []controller.Controller) Server {
	services := make(map[string]controller.Controller)
	var cont controller.Controller

	for i := 0; i < len(s.Services); i++ {
		for j := 0; j < len(controllers); j++ {
			if controllers[j].Name == s.Services[i].Controller {
				cont = controllers[j]
			}
		}
		services[s.Services[i].Route] = cont
	}

	return Server{Port: s.Port, TargetLog: s.TargetLog, Services: services}
}

/* Private configuration is meant to be adapted to the public one by converting yaml to functions */
type configuration struct {
	Database    database     `yaml:"database"`
	Models      []model      `yaml:"model(s)"`
	Controllers []Controller `yaml:"controller(s)"`
	Server      server       `yaml:"server"`
}

type Configuration struct {
	Database        *Database
	Models          []model2.Model
	Controllers     []controller.Controller
	Server          Server
	DatabaseClosure func()
}

// Adapt adapts the configuration, converting FallbackJSON to actual controllers.
func (c configuration) adapt() *Configuration {
	var controllers []controller.Controller
	// This need a SQL handle for some reason

	if c.Database.Path == "" || c.Database.InitQuery == "" {
		log.Warn().Msg("Missing Database in main.yml : Models are disabled")
		// Set all the models to nil, effectively disabling models
		for i := 0; i < len(c.Controllers); i++ {
			JSON, err := json.Marshal(c.Controllers[i].Fallback)
			if err != nil {
				log.Fatal().Err(err).Msg("JSON error in Controller : " + c.Controllers[i].Name)
			}
			newController := controller.Create(c.Controllers[i].Name, nil, JSON, c.Controllers[i].CORS)
			controllers = append(controllers, newController)
		}

		return &Configuration{
			Database:        nil,
			Controllers:     controllers,
			Models:          nil,
			Server:          c.Server.adapt(controllers),
			DatabaseClosure: nil,
		}

	} else {
		// call closeDB to defer the db close
		db, closeDB := database2.Create(c.Database.InitQuery, c.Database.Path)
		var newmodels []model2.Model
		var controllermodel *model2.Model

		// Adapt all the models to actual data models
		for i := 0; i < len(c.Models); i++ {
			newmodels = append(newmodels, c.Models[i].adapt(db))
		}

		for i := 0; i < len(c.Controllers); i++ {
			JSON, err := json.Marshal(c.Controllers[i].Fallback)
			if err != nil {
				log.Fatal().Err(err).Msg("JSON error in Controller : " + c.Controllers[i].Name)
			}
			// The model the controller should use
			for j := 0; j < len(newmodels); j++ {
				if c.Controllers[i].Model == newmodels[j].Name {
					controllermodel = &newmodels[j]
				}
			}

			newController := controller.Create(c.Controllers[i].Name, controllermodel, JSON, c.Controllers[i].CORS)
			controllers = append(controllers, newController)
		}

		return &Configuration{
			Database:        c.Database.adapt(db),
			Controllers:     controllers,
			Models:          newmodels,
			Server:          c.Server.adapt(controllers),
			DatabaseClosure: closeDB,
		}

	}
}

func create(filename string) (*Configuration, error) {
	// Read YAML file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// Unmarshal YAML data into Configuration struct
	var config configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config.adapt(), nil
}

// Setup the config + logging
func Setup(path string) (*Configuration, func()) {
	var multi zerolog.LevelWriter
	var closeFile func()

	conf, err := create(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Something went wrong with generating config from main.yml")
	}
	targetLog := conf.Server.TargetLog

	if targetLog != "" {
		/* Setup logging :  Get logging file and set MultiLevelWriting*/
		file, err := os.OpenFile(targetLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal().Err(err).Msg("Error opening log file")
		}
		closeFile = func() {
			err := file.Close()
			if err != nil {
				log.Fatal().Err(err).Msg("Error while closing log file")
			}
		}
		multi = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, file)
	} else {
		multi = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	return conf, closeFile
}
