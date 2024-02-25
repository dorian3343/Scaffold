package configuration

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"service/controller"
	"service/model"
)

type Database struct {
	InitQuery string `yaml:"init-query"`
	Path      string `yaml:"path"`
}

type QueryTemplate struct {
	JsonTemplate map[string]interface{} `yaml:"json-template"`
	Name         string                 `yaml:"name"`
}

type Model struct {
	QueryTemplate QueryTemplate `yaml:"query-template"`
}

type Controller struct {
	Fallback interface{} `yaml:"fallback"`
	Name     string      `yaml:"name"`
	Model    string      `yaml:"model"`
	CORS     bool        `yaml:"CORS"`
}

/* Private server config to only be used for constructing the public one*/
type server struct {
	Port     int `yaml:"port"`
	Services []struct {
		Route      string `yaml:"route"`
		Controller string `yaml:"controller"`
	} `yaml:"service(s)"`
}

type Server struct {
	Port     int
	Services map[string]controller.Controller
}

func (s server) Adapt(controllers []controller.Controller) Server {
	services := make(map[string]controller.Controller)

	for i := 0; i < len(s.Services); i++ {
		var cont controller.Controller
		for j := 0; j < len(controllers); j++ {
			if controllers[j].Name == s.Services[i].Controller {
				cont = controllers[j]
			}
		}
		services[s.Services[i].Route] = cont
	}

	return Server{Port: s.Port, Services: services}
}

/* Private configuration is meant to be adapted to the public one by converting yaml to functions */
type configuration struct {
	Database    Database     `yaml:"database"`
	Models      []Model      `yaml:"model(s)"`
	Controllers []Controller `yaml:"controller(s)"`
	Server      server       `yaml:"server"`
}

type Configuration struct {
	Database    Database
	Models      []model.Model
	Controllers []controller.Controller
	Server      Server
}

// Adapt adapts the configuration, converting FallbackJSON to actual controllers.
func (c configuration) Adapt() *Configuration {
	var controllers []controller.Controller
	for i := 0; i < len(c.Controllers); i++ {
		JSON, err := json.Marshal(c.Controllers[i].Fallback)
		if err != nil {
			log.Fatal().Err(err).Msg("JSON error in Controller : " + c.Controllers[i].Name)
		}
		newController := controller.Create(c.Controllers[i].Name, nil, JSON, c.Controllers[i].CORS)
		controllers = append(controllers, newController)
	}
	return &Configuration{
		Database:    c.Database,
		Controllers: controllers,
		Models:      nil,
		Server:      c.Server.Adapt(controllers),
	}
}

func Create(filename string) (*Configuration, error) {
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
	return config.Adapt(), nil
}
