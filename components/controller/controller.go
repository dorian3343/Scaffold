package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"service/components/model"
	"slices"
)

/*
Type Controller implements:
- Name : This is the identifier the yaml uses to parse it into an actual endpoint
- http.Handler interface : This is the function that tells the controller what to do
- Model : Handles any interaction witht the database
*/
type Controller struct {
	Name     string
	Model    *model.Model
	Fallback []byte
	cors     string
	cache    string
	verb     string
	http.Handler
}

/* Constructor for the controller, outside of package used like this 'Controller.Create(x,y)' */
func Create(name string, datamodel *model.Model, fallback []byte, cors string, cache string, verb string) Controller {
	return Controller{Name: name, Model: datamodel, Fallback: fallback, cors: cors, cache: cache, verb: verb}
}

func (c Controller) handleNoModelRequest(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(c.Fallback)
	if err != nil {
		log.Err(err).Msg("Something went wrong with Fallback")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func (c Controller) handleHeaders(w http.ResponseWriter) {
	if c.cors != "" {
		w.Header().Set("Access-Control-Allow-Origin", c.cors)
	}
	if c.cache != "" {
		w.Header().Set("Cache-Control", c.cache)
	}

}

/* logic is the function to fulfill the http.Handler interface. */
func (c Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if c.verb != "" && c.verb != r.Method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//set Headers
	c.handleHeaders(w)
	if c.Model == nil {
		c.handleNoModelRequest(w)
	} else {
		var results []map[string]interface{}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Err(err).Msg("Something went wrong with reading the request body")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Trace().Msg("Building Query in : " + c.Name)
		// make the db query
		query, err := c.Model.Querybuilder(body)
		if err != nil {
			if err.Error() == "JSON request does not match spec" {
				log.Err(err).Msg("Something went wrong with building query")
				http.Error(w, "JSON request does not match spec", http.StatusBadRequest)
				return
			} else {
				log.Err(err).Msg("Something went wrong with building query")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		log.Trace().Msg("Running Query in : " + c.Name)
		// Queries the database
		result, err := c.Model.Query(query)
		if err != nil {
			log.Err(err).Msg("Something went wrong with querying database")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer result.Close()

		// Convert the result to a slice of maps
		columns, err := result.Columns()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		for result.Next() {
			row := make(map[string]interface{})
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range columns {
				valuePtrs[i] = &values[i]
			}
			if err := result.Scan(valuePtrs...); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			for i, col := range columns {
				val := values[i]
				row[col] = val
			}
			results = append(results, row)
		}

		// Marshal the results into JSON
		resp, err := json.Marshal(results)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Send the JSON response
		log.Trace().Msg("Sending response in  : " + c.Name)
		_, err = w.Write(resp)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// Sets up and attaches the controllers to http at the proper routes
func SetupControllers(services map[string]Controller) {
	/* Setup http handling : Controllers + Models + main http server */
	fmt.Println("---------Tree------------")
	var wrn []string

	for route, handler := range services {
		log.Trace().Msgf("New Route : | %s : %s", route, handler.Name)
		if route == "" {
			log.Fatal().Err(errors.New("Missing route")).Msg("Something went wrong with setting up Controllers")
		}
		http.Handle(route, handler)
		if handler.Name == "" {
			wrn = append(wrn, fmt.Sprintf("Empty controller for Route: '%s'", route))
		}

		// Check for empty fallbacks
		e := []byte("null")
		if slices.Equal(e, handler.Fallback) {
			wrn = append(wrn, fmt.Sprintf("Empty Fallback for Route: '%s'", route))
		}

	}
	fmt.Println("---------Tree------------")

	/* print tree warnings */
	if len(wrn) != 0 {
		for i := 0; i < len(wrn); i++ {
			log.Warn().Msg(wrn[i])
		}
	}

}
