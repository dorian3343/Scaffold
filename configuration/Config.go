package configuration

import (
	"database/sql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"service/components/controller"
	model2 "service/components/model"
)

/* Private server config to only be used for constructing the public one*/
type server struct {
	Port      int    `yaml:"port"`
	Static    string `yaml:"static"`
	TargetLog string `yaml:"target-log"`
	Services  []struct {
		Route      string `yaml:"route"`
		Controller string `yaml:"controller"`
	} `yaml:"$service"`
}

type Server struct {
	Port      int
	TargetLog string
	Static    string
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

	return Server{Port: s.Port, Static: s.Static, TargetLog: s.TargetLog, Services: services}
}

/* Private configuration is meant to be adapted to the public one by converting yaml to functions */
type configuration struct {
	Database    database     `yaml:"database"`
	Models      []model      `yaml:"$model"`
	Controllers []Controller `yaml:"$controller"`
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
func (c configuration) adapt() (*Configuration, error) {

	var controllers []controller.Controller
	var databasePointer *Database
	var models []model2.Model
	var databaseClosure func()

	if c.Database.Path == "" || c.Database.InitQuery == "" {
		log.Warn().Msg("Missing Database in main.yml : Models are disabled")
		// Set all the models to nil, effectively disabling models
		for i := 0; i < len(c.Controllers); i++ {
			newController, err := c.Controllers[i].adapt(nil)
			if err != nil {
				return nil, err
			}
			controllers = append(controllers, newController)
		}
		databasePointer = nil
		models = nil
		databaseClosure = nil

	} else {
		// call closeDB to defer the db close
		db, closeDB := createDB(c.Database.InitQuery, c.Database.Path)
		var controllermodel *model2.Model

		// Adapt all the models to actual data models
		for i := 0; i < len(c.Models); i++ {
			adapted, err := c.Models[i].adapt(db)
			if err != nil {
				return nil, err
			}
			models = append(models, adapted)
		}

		for i := 0; i < len(c.Controllers); i++ {
			// The model the controller should use
			for j := 0; j < len(models); j++ {
				if c.Controllers[i].Model == models[j].Name {
					controllermodel = &models[j]
				}
			}
			newController, err := c.Controllers[i].adapt(controllermodel)
			if err != nil {
				return nil, err
			}
			controllers = append(controllers, newController)
		}
		databasePointer = c.Database.adapt(db)
		databaseClosure = closeDB
	}
	return &Configuration{
		Database:        databasePointer,
		Controllers:     controllers,
		Models:          models,
		Server:          c.Server.adapt(controllers),
		DatabaseClosure: databaseClosure,
	}, nil

}

func create(filename string) (*Configuration, error) {
	// Read YAML file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// Unmarshal YAML data into Configuration struct
	var config configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	finalConf, newerr := config.adapt()
	if newerr != nil {
		return nil, newerr
	}
	return finalConf, nil
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

// Setup a database
func createDB(query string, databaseName string) (*sql.DB, func()) {
	/* Check if Database file exists */
	_, err := os.Stat(databaseName)
	if os.IsNotExist(err) {
		_, err2 := os.Create(databaseName)
		if err2 != nil {
			log.Fatal().Err(err).Msg("Something went wrong with creating Database")
		}
	}
	// Create the Connection
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal Error opening sqlite")
	}

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal Error during table setup")
	}

	closeDB := func() {
		if err := db.Close(); err != nil {
			log.Fatal().Err(err).Msg("Failed to close the database")
		}
	}
	return db, closeDB
}
